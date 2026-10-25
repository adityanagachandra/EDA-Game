package storage

import (
	"bytes"
	"testing"
)

func TestReportImageRoundTripAndAccountDeletion(t *testing.T) {
	store := newTestStore(t)
	student, err := store.CreateStudent("Image Tester")
	if err != nil {
		t.Fatalf("create student: %v", err)
	}
	id, err := store.CreateReport(QuestionReport{
		ReporterID: student.ID,
		Kind:       "message",
		Reason:     "other",
		Detail:     "A screenshot belongs to this report.",
		Status:     "open",
	})
	if err != nil {
		t.Fatalf("create report: %v", err)
	}

	imageStore, ok := store.(ReportImageStore)
	if !ok {
		t.Fatal("store does not implement ReportImageStore")
	}
	image := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if err := imageStore.SetReportImage(id, "image/png", image); err != nil {
		t.Fatalf("set image: %v", err)
	}

	report, err := store.GetReport(id)
	if err != nil {
		t.Fatalf("get report: %v", err)
	}
	if report == nil || !report.HasImage {
		t.Fatalf("report has_image = %+v, want true", report)
	}
	reports, err := store.ListReports("open", 50, 0)
	if err != nil {
		t.Fatalf("list reports: %v", err)
	}
	listed := false
	for _, r := range reports {
		if r.ID == id {
			listed = r.HasImage
		}
	}
	if !listed {
		t.Fatal("open report list did not carry has_image")
	}
	got, found, err := imageStore.GetReportImage(id)
	if err != nil {
		t.Fatalf("get image: %v", err)
	}
	if !found || got.ContentType != "image/png" || !bytes.Equal(got.Bytes, image) {
		t.Fatalf("image round trip = %+v found=%v", got, found)
	}

	if err := store.DeleteAccount(student.ID); err != nil {
		t.Fatalf("delete account: %v", err)
	}
	if _, found, err := imageStore.GetReportImage(id); err != nil || found {
		t.Fatalf("image survived account deletion: found=%v err=%v", found, err)
	}
	report, err = store.GetReport(id)
	if err != nil {
		t.Fatalf("get deleted report: %v", err)
	}
	if report != nil {
		t.Fatalf("report survived account deletion: %+v", report)
	}
}
