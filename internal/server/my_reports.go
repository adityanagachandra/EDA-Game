package server

import (
	"net/http"
	"time"

	"github.com/chuma-beep/mathua/internal/storage"
)

// GET /api/reports/mine returns the authenticated learner's own question
// reports, newest first.
//
// A report is unread when it has a verdict newer than the last time the learner
// opened this list. The report itself carries no reader field: the same row is
// unread for its reporter and read for a moderator, so the marker belongs to
// the account, not to the report.
func (s *Server) handleMyReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sid, _ := r.Context().Value(authStudentKey{}).(string)
	if sid == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	reports, err := s.repo.ListReportsByReporter(sid, 100)
	if err != nil {
		writeError(w, "failed to load reports", http.StatusInternalServerError)
		return
	}
	if reports == nil {
		reports = []storage.QuestionReport{}
	}
	seen, err := s.repo.GetReportSeenAt(sid)
	if err != nil {
		writeError(w, "failed to load report history", http.StatusInternalServerError)
		return
	}
	seenAt, _ := time.Parse(time.RFC3339, seen)
	unread := 0
	for i := range reports {
		if reports[i].Status != "resolved" && reports[i].Status != "dismissed" {
			continue
		}
		resolvedAt, err := time.Parse(time.RFC3339, reports[i].ResolvedAt)
		if err != nil {
			continue
		}
		if seenAt.IsZero() || resolvedAt.After(seenAt) {
			reports[i].Unread = true
			unread++
		}
	}
	writeJSON(w, map[string]interface{}{"reports": reports, "unread": unread})
}

// POST /api/reports/mine/read marks every current verdict as seen.
func (s *Server) handleMyReportsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sid, _ := r.Context().Value(authStudentKey{}).(string)
	if sid == "" {
		writeError(w, "authentication required", http.StatusUnauthorized)
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if err := s.repo.SetReportSeenAt(sid, now); err != nil {
		writeError(w, "failed to mark reports read", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]interface{}{"ok": true, "seen_at": now})
}
