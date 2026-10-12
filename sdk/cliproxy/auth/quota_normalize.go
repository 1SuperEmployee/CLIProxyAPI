package auth

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Normalize only the default shared allowance. Model, overage, credit and
// additional-limit namespaces are intentionally left to existing enforcement.
func selectionQuotaFromSignals(provider string, q QuotaState) (selectionQuotaSnapshot, bool) {
	h := make(http.Header, len(q.Signals))
	for k, v := range q.Signals {
		h.Set(k, v)
	}
	switch provider {
	case "claude":
		window := func(name string, duration time.Duration) (selectionQuotaWindow, bool) {
			prefix := "Anthropic-Ratelimit-Unified-" + name + "-"
			used, err := strconv.ParseFloat(h.Get(prefix+"Utilization"), 64)
			reset, errReset := strconv.ParseInt(h.Get(prefix+"Reset"), 10, 64)
			if err != nil || errReset != nil || !validFraction(used) {
				return selectionQuotaWindow{}, false
			}
			status := h.Get(prefix + "Status")
			if status != "allowed" && status != "allowed_warning" && status != "rejected" {
				return selectionQuotaWindow{}, false
			}
			if status == "rejected" {
				used = 1
			}
			return selectionQuotaWindow{duration, 1 - used, time.Unix(reset, 0), q.ObservedAt}, true
		}
		weekly, ok := window("7d", 7*24*time.Hour)
		short, okShort := window("5h", 5*time.Hour)
		return selectionQuotaSnapshot{weekly, &short}, ok && okShort
	case "codex":
		var windows []selectionQuotaWindow
		for _, name := range []string{"Primary", "Secondary"} {
			prefix := "X-Codex-" + name + "-"
			minutesText := h.Get(prefix + "Window-Minutes")
			// Codex's default rate-limit contract has optional primary/secondary
			// windows. A partially present window is invalid, not inactive.
			if minutesText == "" && h.Get(prefix+"Used-Percent") == "" && h.Get(prefix+"Reset-At") == "" && h.Get(prefix+"Reset-After-Seconds") == "" {
				continue
			}
			minutes, err := strconv.ParseInt(minutesText, 10, 64)
			if err != nil || minutes < 0 || minutes > 10080 {
				return selectionQuotaSnapshot{}, false
			}
			if minutes == 0 {
				// Zero-duration windows are inactive only with zero/missing use.
				if used := h.Get(prefix + "Used-Percent"); used != "" && used != "0" {
					return selectionQuotaSnapshot{}, false
				}
				continue
			}
			used, errUsed := strconv.ParseFloat(h.Get(prefix+"Used-Percent"), 64)
			if errUsed != nil || !validFraction(used/100) {
				return selectionQuotaSnapshot{}, false
			}
			var reset time.Time
			if absolute := h.Get(prefix + "Reset-At"); absolute != "" {
				seconds, errReset := strconv.ParseInt(absolute, 10, 64)
				if errReset != nil {
					return selectionQuotaSnapshot{}, false
				}
				reset = time.Unix(seconds, 0)
			} else {
				seconds, errReset := strconv.ParseInt(h.Get(prefix+"Reset-After-Seconds"), 10, 64)
				if errReset != nil || seconds <= 0 || seconds > minutes*60 {
					return selectionQuotaSnapshot{}, false
				}
				reset = q.ObservedAt.Add(time.Duration(seconds) * time.Second)
			}
			windows = append(windows, selectionQuotaWindow{time.Duration(minutes) * time.Minute, 1 - used/100, reset, q.ObservedAt})
		}
		snapshot, valid := sharedSelectionWindows(windows)
		return codexSharedAvailability(snapshot, valid, h.Get("X-Codex-Allowed"), h.Get("X-Codex-Limit-Reached"))
	}
	return selectionQuotaSnapshot{}, false
}

func selectionQuotaFromBody(provider string, body []byte, observedAt time.Time) (selectionQuotaSnapshot, bool) {
	var root map[string]json.RawMessage
	if json.Unmarshal(body, &root) != nil || root == nil {
		return selectionQuotaSnapshot{}, false
	}
	switch provider {
	case "claude":
		var windows []selectionQuotaWindow
		for _, name := range []string{"five_hour", "seven_day"} {
			raw, exists := root[name]
			if !exists {
				return selectionQuotaSnapshot{}, false
			}
			if string(raw) == "null" {
				continue
			}
			var w struct {
				Used  *float64 `json:"utilization"`
				Reset *string  `json:"resets_at"`
			}
			if json.Unmarshal(raw, &w) != nil || w.Used == nil || !validFraction(*w.Used/100) {
				return selectionQuotaSnapshot{}, false
			}
			if w.Reset == nil && *w.Used == 0 {
				continue // Provider-confirmed inactive window.
			}
			if w.Reset == nil {
				return selectionQuotaSnapshot{}, false
			}
			reset, err := time.Parse(time.RFC3339Nano, *w.Reset)
			if err != nil {
				return selectionQuotaSnapshot{}, false
			}
			duration := 5 * time.Hour
			if name == "seven_day" {
				duration = 7 * 24 * time.Hour
			}
			windows = append(windows, selectionQuotaWindow{duration, 1 - *w.Used/100, reset, observedAt})
		}
		return sharedSelectionWindows(windows)
	case "codex":
		var rate map[string]json.RawMessage
		if json.Unmarshal(root["rate_limit"], &rate) != nil || rate == nil {
			return selectionQuotaSnapshot{}, false
		}
		var windows []selectionQuotaWindow
		for _, name := range []string{"primary_window", "secondary_window"} {
			raw, exists := rate[name]
			if !exists {
				return selectionQuotaSnapshot{}, false
			}
			if string(raw) == "null" {
				continue
			}
			var w struct {
				Used     *float64 `json:"used_percent"`
				Duration *int64   `json:"limit_window_seconds"`
				Reset    *int64   `json:"reset_at"`
				After    *int64   `json:"reset_after_seconds"`
			}
			if json.Unmarshal(raw, &w) != nil || w.Duration == nil || *w.Duration < 0 || *w.Duration > 604800 || w.Used == nil || !validFraction(*w.Used/100) {
				return selectionQuotaSnapshot{}, false
			}
			if *w.Duration == 0 && *w.Used == 0 {
				continue
			}
			var reset time.Time
			if w.Reset != nil {
				reset = time.Unix(*w.Reset, 0)
			} else if w.After != nil && *w.After > 0 && *w.After <= *w.Duration {
				reset = observedAt.Add(time.Duration(*w.After) * time.Second)
			} else {
				return selectionQuotaSnapshot{}, false
			}
			windows = append(windows, selectionQuotaWindow{time.Duration(*w.Duration) * time.Second, 1 - *w.Used/100, reset, observedAt})
		}
		snapshot, valid := sharedSelectionWindows(windows)
		return codexSharedAvailability(snapshot, valid, string(rate["allowed"]), string(rate["limit_reached"]))
	}
	return selectionQuotaSnapshot{}, false
}

func codexSharedAvailability(snapshot selectionQuotaSnapshot, valid bool, allowed, reached string) (selectionQuotaSnapshot, bool) {
	denied := false
	for _, item := range []struct {
		text        string
		deniedValue bool
	}{{allowed, false}, {reached, true}} {
		if item.text == "" {
			continue
		}
		value, err := strconv.ParseBool(item.text)
		if err != nil {
			return selectionQuotaSnapshot{}, false
		}
		denied = denied || value == item.deniedValue
	}
	if denied {
		snapshot.weekly.remainingFraction = 0
		if snapshot.shortWindow != nil {
			snapshot.shortWindow.remainingFraction = 0
		}
	}
	return snapshot, valid
}

func sharedSelectionWindows(windows []selectionQuotaWindow) (selectionQuotaSnapshot, bool) {
	var out selectionQuotaSnapshot
	for _, w := range windows {
		switch {
		case w.duration == 7*24*time.Hour && out.weekly.duration == 0:
			out.weekly = w
		case w.duration > 0 && w.duration < 7*24*time.Hour && out.shortWindow == nil:
			copy := w
			out.shortWindow = &copy
		default:
			return selectionQuotaSnapshot{}, false
		}
	}
	return out, out.weekly.duration != 0
}

func selectionQuotaProvider(auth *Auth) string {
	if auth == nil || auth.AuthKind() != AuthKindOAuth {
		return ""
	}
	provider := strings.ToLower(strings.TrimSpace(auth.Provider))
	if provider == "claude" || provider == "codex" {
		return provider
	}
	return ""
}
