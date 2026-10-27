// Package fire holds the runtime mode for FIRe, the implicit-repetition
// ledger. FIRe is an observation at rest: it may record that a learner would
// have earned implicit credit, but it must not move a due date, change mastery
// or alter what the scheduler serves.
package fire

import "strings"

// Mode is the resolved startup mode for FIRe. There is deliberately no active
// mode: the only two states are "keep the ledger" and "do not process FIRe at
// all". Anything unrecognized resolves to Shadow so a typo cannot silently
// disable the ledger or enable scheduling effects.
type Mode string

const (
	// Shadow records ledger rows only. It never writes ConceptProgress, FSRS
	// state, due dates or review timestamps, and the scheduler ignores its
	// banked ordinals.
	Shadow Mode = "shadow"
	// Off does not allocate credit and does not apply FIRe outbox effects.
	Off Mode = "off"
)

// Parse resolves the configured value. It is total and fail-closed: empty,
// whitespace, malformed and unknown values all become Shadow.
func Parse(raw string) Mode {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(Off):
		return Off
	default:
		return Shadow
	}
}

// CollectsLedger reports whether FIRe should record observations. Only Shadow
// does. The method is the single question the engine and outbox ask, so a new
// mode cannot accidentally inherit ledger writes.
func (m Mode) CollectsLedger() bool { return m == Shadow }

// String keeps log lines readable and stable.
func (m Mode) String() string { return string(m) }
