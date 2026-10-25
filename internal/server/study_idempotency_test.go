package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/chuma-beep/mathua/internal/concepts"
)

func studyPracticeQuestion(t *testing.T, mux *http.ServeMux, token, conceptID string) (question, answer string) {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/lessons/"+conceptID+"/practice?count=1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("practice: %d %s", rec.Code, rec.Body.String())
	}
	var set struct {
		Questions []struct {
			Question string `json:"question"`
			Answer   string `json:"answer"`
		} `json:"questions"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &set); err != nil {
		t.Fatalf("decode practice: %v", err)
	}
	if len(set.Questions) == 0 {
		t.Fatal("expected a served question")
	}
	return set.Questions[0].Question, set.Questions[0].Answer
}

func studyAnswerBody(conceptID, answer, question, submissionID string) string {
	b, _ := json.Marshal(map[string]interface{}{
		"concept_id":    conceptID,
		"answer":        answer,
		"elapsed":       5.0,
		"question":      question,
		"submission_id": submissionID,
	})
	return string(b)
}

// TestStudyAnswer_RequiresSubmissionID closes the legacy write path: a study
// answer that awards XP must carry a stable client key, so an old client gets
// a clear 400 instead of an unsafe retry.
func TestStudyAnswer_RequiresSubmissionID(t *testing.T) {
	stores := quizSchedulingStores(t)
	cs := []concepts.Concept{{
		ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
	}}
	for name, store := range stores {
		t.Run(name, func(t *testing.T) {
			_, mux := quizSchedulingServerWith(t, store, cs)
			token, studentID := quizGuest(t, mux)
			question, answer := studyPracticeQuestion(t, mux, token, "a")

			missing := quizSchedulingPost(t, mux, "/api/study/answer", token,
				studyAnswerBody("a", answer, question, ""))
			if missing.Code != 400 {
				t.Fatalf("missing submission id: got %d %s", missing.Code, missing.Body.String())
			}
			attempts, _ := store.GetRecentAttemptsForConcept(studentID, "a", 20)
			if len(attempts) != 0 {
				t.Fatalf("missing submission id wrote %d attempts", len(attempts))
			}
			student, _ := store.GetStudent(studentID)
			if student.XPTotal != 0 {
				t.Fatalf("missing submission id awarded XP %d", student.XPTotal)
			}
		})
	}
}

// TestStudyAnswer_SameSubmissionReplays pins the contract the client relies on:
// a retry carries the same submission id, gets the same body and writes no
// second attempt or XP event.
func TestStudyAnswer_SameSubmissionReplays(t *testing.T) {
	stores := quizSchedulingStores(t)
	cs := []concepts.Concept{{
		ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
	}}
	for name, store := range stores {
		t.Run(name, func(t *testing.T) {
			_, mux := quizSchedulingServerWith(t, store, cs)
			token, studentID := quizGuest(t, mux)
			question, answer := studyPracticeQuestion(t, mux, token, "a")
			body := studyAnswerBody("a", answer, question, "submission-1")

			first := quizSchedulingPost(t, mux, "/api/study/answer", token, body)
			if first.Code != 200 {
				t.Fatalf("first answer: %d %s", first.Code, first.Body.String())
			}
			second := quizSchedulingPost(t, mux, "/api/study/answer", token, body)
			if second.Code != 200 {
				t.Fatalf("replay: %d %s", second.Code, second.Body.String())
			}
			if first.Body.String() != second.Body.String() {
				t.Fatalf("replay body differs:\nfirst  %s\nsecond %s", first.Body.String(), second.Body.String())
			}

			attempts, err := store.GetRecentAttemptsForConcept(studentID, "a", 20)
			if err != nil {
				t.Fatalf("attempts: %v", err)
			}
			if len(attempts) != 1 {
				t.Fatalf("expected one attempt after replay, got %d", len(attempts))
			}
			student, err := store.GetStudent(studentID)
			if err != nil {
				t.Fatalf("student: %v", err)
			}
			if student.XPTotal <= 0 {
				t.Fatalf("expected XP for a correct answer, got %d", student.XPTotal)
			}
			xpAfterReplay := student.XPTotal

			mismatch := quizSchedulingPost(t, mux, "/api/study/answer", token,
				studyAnswerBody("a", "wrong", question, "submission-1"))
			if mismatch.Code != 409 {
				t.Fatalf("expected 409 for a reused id with a different answer, got %d: %s", mismatch.Code, mismatch.Body.String())
			}
			attempts, _ = store.GetRecentAttemptsForConcept(studentID, "a", 20)
			if len(attempts) != 1 {
				t.Fatalf("payload mismatch wrote an attempt: %d", len(attempts))
			}
			student, _ = store.GetStudent(studentID)
			if student.XPTotal != xpAfterReplay {
				t.Fatalf("payload mismatch moved XP from %d to %d", xpAfterReplay, student.XPTotal)
			}

			// A different id is a legitimate second answer, not a replay.
			other := quizSchedulingPost(t, mux, "/api/study/answer", token,
				studyAnswerBody("a", answer, question, "submission-2"))
			if other.Code != 200 {
				t.Fatalf("second logical answer: %d %s", other.Code, other.Body.String())
			}
			attempts, _ = store.GetRecentAttemptsForConcept(studentID, "a", 20)
			if len(attempts) != 2 {
				t.Fatalf("expected two attempts after a new id, got %d", len(attempts))
			}
		})
	}
}

// TestStudyAnswer_ReplayPreservesHaltedFlag covers the response field that
// changes behaviour on the client: the second miss halts the run, and a retry
// must return the same halted body rather than a fresh, unhalted verdict.
func TestStudyAnswer_ReplayPreservesHaltedFlag(t *testing.T) {
	stores := quizSchedulingStores(t)
	cs := []concepts.Concept{{
		ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
	}}
	for name, store := range stores {
		t.Run(name, func(t *testing.T) {
			_, mux := quizSchedulingServerWith(t, store, cs)
			token, studentID := quizGuest(t, mux)
			question, _ := studyPracticeQuestion(t, mux, token, "a")

			first := quizSchedulingPost(t, mux, "/api/study/answer", token,
				studyAnswerBody("a", "wrong-one", question, "sub-halt-1"))
			if first.Code != 200 {
				t.Fatalf("first miss: %d %s", first.Code, first.Body.String())
			}
			secondBody := studyAnswerBody("a", "wrong-two", question, "sub-halt-2")
			second := quizSchedulingPost(t, mux, "/api/study/answer", token, secondBody)
			if second.Code != 200 {
				t.Fatalf("second miss: %d %s", second.Code, second.Body.String())
			}
			var body map[string]interface{}
			if err := json.Unmarshal(second.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode second: %v", err)
			}
			if halted, _ := body["halted"].(bool); !halted {
				t.Fatalf("second miss did not halt: %s", second.Body.String())
			}
			replay := quizSchedulingPost(t, mux, "/api/study/answer", token, secondBody)
			if replay.Code != 200 {
				t.Fatalf("replay: %d %s", replay.Code, replay.Body.String())
			}
			if replay.Body.String() != second.Body.String() {
				t.Fatalf("halted replay differs:\nsecond %s\nreplay %s", second.Body.String(), replay.Body.String())
			}
			attempts, _ := store.GetRecentAttemptsForConcept(studentID, "a", 20)
			if len(attempts) != 2 {
				t.Fatalf("expected two attempts, got %d", len(attempts))
			}
		})
	}
}

// TestStudyAnswer_ConcurrentSameSubmissionWritesOnce runs the retry race
// directly: eight requests with one id must produce one attempt and one XP
// award, not eight.
func TestStudyAnswer_ConcurrentSameSubmissionWritesOnce(t *testing.T) {
	stores := quizSchedulingStores(t)
	cs := []concepts.Concept{{
		ID: "a", Label: "A", Domain: "d", GradingType: "numeric", Prerequisites: []string{},
		MasteryThreshold: concepts.MasteryThreshold{Streak: 3, AvgTimeSeconds: 60},
	}}
	for name, store := range stores {
		t.Run(name, func(t *testing.T) {
			_, mux := quizSchedulingServerWith(t, store, cs)
			token, studentID := quizGuest(t, mux)
			question, answer := studyPracticeQuestion(t, mux, token, "a")
			body := studyAnswerBody("a", answer, question, "submission-concurrent")

			const workers = 8
			codes := make(chan int, workers)
			bodies := make(chan string, workers)
			var wg sync.WaitGroup
			wg.Add(workers)
			for i := 0; i < workers; i++ {
				go func() {
					defer wg.Done()
					rec := quizSchedulingPost(t, mux, "/api/study/answer", token, body)
					codes <- rec.Code
					bodies <- rec.Body.String()
				}()
			}
			wg.Wait()
			close(codes)
			close(bodies)
			for code := range codes {
				if code != 200 {
					t.Fatalf("concurrent replay returned %d", code)
				}
			}
			var first string
			for b := range bodies {
				if first == "" {
					first = b
					continue
				}
				if b != first {
					t.Fatalf("concurrent replay bodies differ:\n%s\n%s", first, b)
				}
			}

			attempts, err := store.GetRecentAttemptsForConcept(studentID, "a", 20)
			if err != nil {
				t.Fatalf("attempts: %v", err)
			}
			if len(attempts) != 1 {
				t.Fatalf("expected one attempt from concurrent retries, got %d", len(attempts))
			}
			student, _ := store.GetStudent(studentID)
			if student.XPTotal <= 0 {
				t.Fatalf("expected one XP award, got %d", student.XPTotal)
			}
		})
	}
}
