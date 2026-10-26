package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

func TestMyReportsAreScopedToTheCallerAndCarryUnreadVerdicts(t *testing.T) {
	s, mux, store := reportsTestMux(t)
	token, student := signupReportLearner(t, s, "Ada Example", "ada")
	_, other := signupReportLearner(t, s, "Bob Example", "bob")

	now := time.Now().UTC()
	resolvedID, err := store.CreateReport(storage.QuestionReport{
		ReporterID: student.ID,
		ConceptID:  "a",
		Kind:       "question",
		Question:   "what is 2+2?",
		CreatedAt:  now,
	})
	if err != nil {
		t.Fatalf("create report to resolve: %v", err)
	}
	// CreateReport stores only the reporter's fields. A verdict is written by
	// the moderation path, which is the path that sets resolved_at and is what
	// the unread comparison must read.
	if _, err := store.UpdateReportStatusAudited("moderator", resolvedID, "resolved", "the expected answer was wrong"); err != nil {
		t.Fatalf("resolve report: %v", err)
	}
	if _, err := store.CreateReport(storage.QuestionReport{
		ReporterID: student.ID,
		ConceptID:  "b",
		Kind:       "question",
		Question:   "what is 3+3?",
		Status:     "open",
		CreatedAt:  now,
	}); err != nil {
		t.Fatalf("create open report: %v", err)
	}
	if _, err := store.CreateReport(storage.QuestionReport{
		ReporterID: other.ID,
		ConceptID:  "a",
		Kind:       "question",
		Question:   "someone else's report",
		Status:     "open",
		CreatedAt:  now,
	}); err != nil {
		t.Fatalf("create other report: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/reports/mine", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body.String())
	}
	var got struct {
		Reports []storage.QuestionReport `json:"reports"`
		Unread  int                      `json:"unread"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got.Reports) != 2 {
		t.Fatalf("got %d reports, want the caller's 2", len(got.Reports))
	}
	if got.Unread != 1 {
		t.Errorf("unread = %d, want 1 resolved report", got.Unread)
	}
	for _, report := range got.Reports {
		if report.ReporterID != student.ID {
			t.Errorf("report %d belongs to %s, not the caller", report.ID, report.ReporterID)
		}
		if report.ReporterUsername != "ada" {
			t.Errorf("reporter username = %q, want ada", report.ReporterUsername)
		}
	}

	read := httptest.NewRequest(http.MethodPost, "/api/reports/mine/read", nil)
	read.Header.Set("Authorization", "Bearer "+token)
	readRec := httptest.NewRecorder()
	mux.ServeHTTP(readRec, read)
	if readRec.Code != http.StatusOK {
		t.Fatalf("mark read status %d: %s", readRec.Code, readRec.Body.String())
	}

	after := httptest.NewRequest(http.MethodGet, "/api/reports/mine", nil)
	after.Header.Set("Authorization", "Bearer "+token)
	afterRec := httptest.NewRecorder()
	mux.ServeHTTP(afterRec, after)
	if afterRec.Code != http.StatusOK {
		t.Fatalf("second read status %d: %s", afterRec.Code, afterRec.Body.String())
	}
	got = struct {
		Reports []storage.QuestionReport `json:"reports"`
		Unread  int                      `json:"unread"`
	}{}
	if err := json.Unmarshal(afterRec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode second read: %v", err)
	}
	if got.Unread != 0 {
		t.Errorf("unread after read = %d, want 0", got.Unread)
	}
}

func TestMyReportsRequiresAuthentication(t *testing.T) {
	_, mux, _ := reportsTestMux(t)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/reports/mine", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d, want 401", rec.Code)
	}
}
