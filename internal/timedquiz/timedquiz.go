// Package timedquiz holds the policy for the timed assessment. It is separate
// from the 50 XP mastery check in internal/quiz: a mastery check verifies
// whether a learner is due for a short review lesson, while a timed quiz
// samples already learned concepts under a server-owned deadline.
package timedquiz

import (
	"errors"
	"time"
)

var (
	ErrNotDue       = errors.New("timed quiz is not due")
	ErrNoCandidates = errors.New("no studied concepts are available for a timed quiz")
	ErrNotFound     = errors.New("timed quiz not found")
	ErrExpired      = errors.New("timed quiz deadline has passed")
)

const (
	// GateXP is the amount of positive XP a learner must earn after their
	// rollout baseline before a timed quiz becomes due.
	GateXP = 150

	// QuizSize is the number of concepts sampled in one timed quiz.
	QuizSize = 5

	// Duration is the total server-authoritative run time. Five questions in
	// ten minutes leaves room to read and reason without making speed the
	// subject of the assessment.
	Duration = 10 * time.Minute
)

const (
	StatusActive    = "active"
	StatusCompleted = "completed"

	ItemServed     = "served"
	ItemSubmitting = "submitting"
	ItemAnswered   = "answered"
	ItemSkipped    = "skipped"
)
