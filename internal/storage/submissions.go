package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// The study path writes the same four rows on every accepted answer:
// the submission record, the attempt, the concept progress and the topic
// speed. They are one transaction here because a crash between them used to
// leave a learner with progress but no attempt, or an attempt but no XP.
const upsertProgressSQL = `
	INSERT INTO concept_progress
		(student_id, concept_id, status, streak, best_streak,
		 avg_response_time, attempts, last_attempted, last_reviewed,
		 next_review_due, sm2_repetitions, sm2_interval, sm2_efactor,
		 mastered_at, weakness_score, placement_seeded, is_provisional,
		 fsrs_stability, fsrs_difficulty, fsrs_last_review)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(student_id, concept_id) DO UPDATE SET
		status            = excluded.status,
		streak            = excluded.streak,
		best_streak       = excluded.best_streak,
		avg_response_time = excluded.avg_response_time,
		attempts          = excluded.attempts,
		last_attempted    = excluded.last_attempted,
		last_reviewed     = excluded.last_reviewed,
		next_review_due   = excluded.next_review_due,
		sm2_repetitions   = excluded.sm2_repetitions,
		sm2_interval      = excluded.sm2_interval,
		sm2_efactor       = excluded.sm2_efactor,
		mastered_at       = excluded.mastered_at,
		weakness_score    = excluded.weakness_score,
		placement_seeded  = excluded.placement_seeded,
		is_provisional    = excluded.is_provisional,
		fsrs_stability    = excluded.fsrs_stability,
		fsrs_difficulty   = excluded.fsrs_difficulty,
		fsrs_last_review  = excluded.fsrs_last_review
`

const upsertTopicSpeedSQL = `
	INSERT INTO student_topic_speed (student_id, concept_id, efactor, interval, repetitions, learning_speed, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(student_id, concept_id) DO UPDATE SET
		efactor = excluded.efactor,
		interval = excluded.interval,
		repetitions = excluded.repetitions,
		learning_speed = excluded.learning_speed,
		updated_at = excluded.updated_at
`

type attemptInserter func(tx *sql.Tx, entry AttemptEntry) (int64, error)

// commitStudySubmissionTx is shared by both SQL stores. rebind turns the
// SQLite `?` spelling into the store's placeholders. Grading happens before
// this function is called, so the transaction contains only durable writes.
func commitStudySubmissionTx(db *sql.DB, rebind func(string) string, insertAttempt attemptInserter, c StudySubmissionCommit) (StudySubmissionOutcome, error) {
	if c.StudentID == "" || c.SubmissionID == "" || c.Operation == "" {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: missing scope")
	}
	if c.RequestHash == "" {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: missing request hash")
	}
	tx, err := db.Begin()
	if err != nil {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: begin: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC().UnixNano()
	res, err := tx.Exec(rebind(`
		INSERT INTO answer_submissions
			(student_id, operation, submission_id, request_hash, response_json, attempt_id, created_at)
		VALUES (?, ?, ?, ?, ?, 0, ?)
		ON CONFLICT(student_id, operation, submission_id) DO NOTHING
	`), c.StudentID, c.Operation, c.SubmissionID, c.RequestHash, c.ResponseJSON, now)
	if err != nil {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: insert: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		existing, err := getStudySubmissionTx(tx, rebind, c.StudentID, c.Operation, c.SubmissionID)
		if err != nil {
			return StudySubmissionOutcome{}, err
		}
		if existing.RequestHash != c.RequestHash {
			return StudySubmissionOutcome{}, ErrStudySubmissionPayloadMismatch
		}
		if existing.ResponseJSON == "" {
			return StudySubmissionOutcome{}, ErrStudySubmissionInProgress
		}
		return StudySubmissionOutcome{
			Replayed:     true,
			ResponseJSON: existing.ResponseJSON,
			AttemptID:    existing.AttemptID,
		}, nil
	}

	attemptID, err := insertAttempt(tx, c.Attempt)
	if err != nil {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: attempt: %w", err)
	}
	if _, err := tx.Exec(rebind(`
		UPDATE answer_submissions SET attempt_id = ?
		WHERE student_id = ? AND operation = ? AND submission_id = ?
	`), attemptID, c.StudentID, c.Operation, c.SubmissionID); err != nil {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: link attempt: %w", err)
	}
	if c.Progress != nil {
		if _, err := tx.Exec(rebind(upsertProgressSQL),
			c.Progress.StudentID, c.Progress.ConceptID, c.Progress.Status, c.Progress.Streak, c.Progress.BestStreak,
			c.Progress.AvgResponseTime, c.Progress.Attempts,
			nullTime(c.Progress.LastAttempted), nullTime(c.Progress.LastReviewed),
			nullTime(c.Progress.NextReviewDue), c.Progress.SM2Repetitions, c.Progress.SM2Interval,
			c.Progress.SM2EFactor, nullTime(c.Progress.MasteredAt), c.Progress.WeaknessScore,
			boolToInt(c.Progress.PlacementSeeded), boolToInt(c.Progress.IsProvisional),
			nullFloat(c.Progress.FSRSStability), nullFloat(c.Progress.FSRSDifficulty), nullTime(c.Progress.FSRSLastReview),
		); err != nil {
			return StudySubmissionOutcome{}, fmt.Errorf("study submission: progress: %w", err)
		}
	}
	if c.TopicSpeed != nil {
		ts := *c.TopicSpeed
		if ts.LearningSpeed < 0.5 {
			ts.LearningSpeed = 0.5
		}
		if ts.LearningSpeed > 2.0 {
			ts.LearningSpeed = 2.0
		}
		if ts.EFactor == 0 {
			ts.EFactor = 2.5
		}
		if _, err := tx.Exec(rebind(upsertTopicSpeedSQL),
			ts.StudentID, ts.ConceptID, ts.EFactor, ts.Interval, ts.Repetitions, ts.LearningSpeed,
			time.Now().UTC().Format("2006-01-02 15:04:05"),
		); err != nil {
			return StudySubmissionOutcome{}, fmt.Errorf("study submission: topic speed: %w", err)
		}
	}
	if c.XPAmount != 0 {
		eventKey := fmt.Sprintf("attempt:%d", attemptID)
		xpRes, err := tx.Exec(rebind(`
			INSERT INTO xp_events (student_id, amount, reason, event_key, created_at)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(student_id, event_key) DO NOTHING
		`), c.StudentID, c.XPAmount, c.XPReason, eventKey, now)
		if err != nil {
			return StudySubmissionOutcome{}, fmt.Errorf("study submission: xp event: %w", err)
		}
		if n, _ := xpRes.RowsAffected(); n > 0 {
			today := time.Now().UTC().Format("2006-01-02")
			if _, err := tx.Exec(rebind(`
				UPDATE students
				SET xp_total = CASE WHEN xp_total + ? < 0 THEN 0 ELSE xp_total + ? END,
				    xp_today = CASE
				        WHEN xp_date = ? THEN CASE WHEN xp_today + ? < 0 THEN 0 ELSE xp_today + ? END
				        ELSE CASE WHEN ? < 0 THEN 0 ELSE ? END
				    END,
				    xp_date  = ?
				WHERE id = ?
			`), c.XPAmount, c.XPAmount, today, c.XPAmount, c.XPAmount, c.XPAmount, c.XPAmount, today, c.StudentID); err != nil {
				return StudySubmissionOutcome{}, fmt.Errorf("study submission: xp counters: %w", err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return StudySubmissionOutcome{}, fmt.Errorf("study submission: commit: %w", err)
	}
	return StudySubmissionOutcome{ResponseJSON: c.ResponseJSON, AttemptID: attemptID}, nil
}

func getStudySubmissionTx(tx *sql.Tx, rebind func(string) string, studentID, operation, submissionID string) (*StudySubmission, error) {
	var out StudySubmission
	var created int64
	if err := tx.QueryRow(rebind(`
		SELECT student_id, operation, submission_id, request_hash, response_json, attempt_id, created_at
		FROM answer_submissions
		WHERE student_id = ? AND operation = ? AND submission_id = ?
	`), studentID, operation, submissionID).Scan(
		&out.StudentID, &out.Operation, &out.SubmissionID, &out.RequestHash,
		&out.ResponseJSON, &out.AttemptID, &created,
	); err != nil {
		return nil, fmt.Errorf("study submission: read: %w", err)
	}
	out.CreatedAt = time.Unix(0, created).UTC()
	return &out, nil
}

func getStudySubmission(db *sql.DB, rebind func(string) string, studentID, operation, submissionID string) (*StudySubmission, error) {
	var out StudySubmission
	var created int64
	if err := db.QueryRow(rebind(`
		SELECT student_id, operation, submission_id, request_hash, response_json, attempt_id, created_at
		FROM answer_submissions
		WHERE student_id = ? AND operation = ? AND submission_id = ?
	`), studentID, operation, submissionID).Scan(
		&out.StudentID, &out.Operation, &out.SubmissionID, &out.RequestHash,
		&out.ResponseJSON, &out.AttemptID, &created,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("study submission: read: %w", err)
	}
	out.CreatedAt = time.Unix(0, created).UTC()
	return &out, nil
}
