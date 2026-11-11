package db

import (
	"testing"
)

func TestCalcCVSS31Score(t *testing.T) {
	cases := []struct {
		vector string
		want   float64
	}{
		// Critical — well-known network-exploitable, no auth, full impact
		{"AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H", 9.8},
		// Scope Changed variant (e.g. EternalBlue-like)
		{"AV:N/AC:L/PR:N/UI:N/S:C/C:H/I:H/A:H", 10.0},
		// High — requires user interaction
		{"AV:N/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H", 8.8},
		// Medium
		{"AV:N/AC:L/PR:L/UI:N/S:U/C:L/I:L/A:N", 5.4},
		// Low — local access only, low impact
		{"AV:L/AC:L/PR:L/UI:N/S:U/C:L/I:N/A:N", 3.3},
		// No impact → 0.0
		{"AV:N/AC:L/PR:N/UI:N/S:U/C:N/I:N/A:N", 0.0},
		// Optional CVSS:3.1/ prefix
		{"CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H", 9.8},
	}
	for _, tc := range cases {
		got, err := CalcCVSS31Score(tc.vector)
		if err != nil {
			t.Errorf("CalcCVSS31Score(%q) error: %v", tc.vector, err)
			continue
		}
		if got != tc.want {
			t.Errorf("CalcCVSS31Score(%q) = %.1f, want %.1f", tc.vector, got, tc.want)
		}
	}

	// Error cases
	_, err := CalcCVSS31Score("AV:N/AC:L")
	if err == nil {
		t.Error("expected error for incomplete vector, got nil")
	}
	_, err = CalcCVSS31Score("AV:X/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H")
	if err == nil {
		t.Error("expected error for unknown AV value, got nil")
	}
}
