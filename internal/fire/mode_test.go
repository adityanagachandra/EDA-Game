package fire

import "testing"

func TestParseFallsBackToShadow(t *testing.T) {
	cases := []struct {
		raw  string
		want Mode
	}{
		{"", Shadow},
		{"   ", Shadow},
		{"bogus", Shadow},
		{"active", Shadow},
		{"shadow", Shadow},
		{"SHADOW", Shadow},
		{" off ", Off},
		{"OFF", Off},
	}
	for _, tc := range cases {
		if got := Parse(tc.raw); got != tc.want {
			t.Errorf("Parse(%q) = %q, want %q", tc.raw, got, tc.want)
		}
	}
}

func TestOnlyShadowCollectsLedger(t *testing.T) {
	if !Shadow.CollectsLedger() {
		t.Fatal("shadow must collect ledger observations")
	}
	if Off.CollectsLedger() {
		t.Fatal("off must not collect ledger observations")
	}
	if Parse("active").CollectsLedger() != true {
		t.Fatal("unknown modes must fail closed to shadow")
	}
}
