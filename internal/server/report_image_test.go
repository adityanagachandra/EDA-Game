package server

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http/httptest"
	"testing"
)

func multipartReport(t *testing.T, fields map[string]string, image []byte) (*bytes.Buffer, string) {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := w.WriteField(key, value); err != nil {
			t.Fatalf("write field %s: %v", key, err)
		}
	}
	if image != nil {
		fw, err := w.CreateFormFile("image", "screenshot.png")
		if err != nil {
			t.Fatalf("create image part: %v", err)
		}
		if _, err := fw.Write(image); err != nil {
			t.Fatalf("write image: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}
	return &body, w.FormDataContentType()
}

func TestReportScreenshotRoundTripsBehindModerationPermission(t *testing.T) {
	e := newAdminEnv(t)
	body, contentType := multipartReport(t, map[string]string{
		"kind":   "message",
		"source": "contact",
		"reason": "other",
		"detail": "The layout breaks when I open the quiz on a phone.",
	}, tinyPNG)

	req := httptest.NewRequest("POST", "/api/reports", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+e.studentTok)
	rec := httptest.NewRecorder()
	e.mux.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("create report with screenshot: %d %s", rec.Code, rec.Body.String())
	}
	var created struct {
		ID       int64 `json:"id"`
		HasImage bool  `json:"has_image"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if !created.HasImage {
		t.Fatal("create response did not report the screenshot")
	}

	rec = e.do(t, "GET", "/api/admin/reports/"+itoa(created.ID), e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("admin report detail: %d %s", rec.Code, rec.Body.String())
	}
	var detail struct {
		Report struct {
			HasImage bool `json:"has_image"`
		} `json:"report"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &detail); err != nil {
		t.Fatalf("decode detail: %v", err)
	}
	if !detail.Report.HasImage {
		t.Fatal("admin detail did not carry has_image")
	}

	rec = e.do(t, "GET", "/api/admin/reports/"+itoa(created.ID)+"/image", e.adminTok, "")
	if rec.Code != 200 {
		t.Fatalf("admin screenshot: %d %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "image/png" {
		t.Fatalf("content type = %q, want image/png", got)
	}
	if !bytes.Equal(rec.Body.Bytes(), tinyPNG) {
		t.Fatal("served screenshot bytes changed")
	}

	rec = e.do(t, "GET", "/api/admin/reports/"+itoa(created.ID)+"/image", e.studentTok, "")
	if rec.Code != 403 {
		t.Fatalf("learner screenshot read: got %d, want 403", rec.Code)
	}

	rec = e.do(t, "GET", "/api/admin/reports/"+itoa(created.ID)+"/image", "", "")
	if rec.Code != 401 {
		t.Fatalf("anonymous screenshot read: got %d, want 401", rec.Code)
	}
}

func TestReportScreenshotRejectsNonImages(t *testing.T) {
	e := newAdminEnv(t)
	body, contentType := multipartReport(t, map[string]string{
		"kind":   "message",
		"source": "contact",
		"reason": "other",
		"detail": "This file is not a screenshot.",
	}, []byte("this is not an image"))

	req := httptest.NewRequest("POST", "/api/reports", body)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer "+e.studentTok)
	rec := httptest.NewRecorder()
	e.mux.ServeHTTP(rec, req)
	if rec.Code != 400 {
		t.Fatalf("non-image upload: got %d, want 400 (%s)", rec.Code, rec.Body.String())
	}
}

func TestReportWithoutScreenshotHasNoImageRoute(t *testing.T) {
	e := newAdminEnv(t)
	rec := e.do(t, "POST", "/api/reports", e.studentTok,
		`{"kind":"message","source":"contact","reason":"other","detail":"A message without a screenshot."}`)
	if rec.Code != 200 {
		t.Fatalf("file report: %d %s", rec.Code, rec.Body.String())
	}
	var created struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	rec = e.do(t, "GET", "/api/admin/reports/"+itoa(created.ID)+"/image", e.adminTok, "")
	if rec.Code != 404 {
		t.Fatalf("missing screenshot: got %d, want 404", rec.Code)
	}
}
