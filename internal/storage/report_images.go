package storage

import (
	"database/sql"
	"fmt"
	"time"
)

// Report screenshots are one row per report, stored as a BLOB/BYTEA rather
// than a path on disk. The moderation record is the authority for whether the
// image exists, and an account deletion already removes its reports; the
// foreign key cascade keeps the two from drifting apart.
func (s *SQLiteStore) SetReportImage(reportID int64, contentType string, data []byte) error {
	_, err := s.db.Exec(`
		INSERT INTO report_images (report_id, content_type, bytes, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(report_id) DO UPDATE SET
			content_type = excluded.content_type,
			bytes = excluded.bytes,
			created_at = excluded.created_at
	`, reportID, contentType, data, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("set report image: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetReportImage(reportID int64) (ReportImage, bool, error) {
	var img ReportImage
	err := s.db.QueryRow(
		`SELECT content_type, bytes FROM report_images WHERE report_id = ?`, reportID,
	).Scan(&img.ContentType, &img.Bytes)
	if err == sql.ErrNoRows {
		return ReportImage{}, false, nil
	}
	if err != nil {
		return ReportImage{}, false, fmt.Errorf("get report image: %w", err)
	}
	return img, true, nil
}

func (s *PostgresStore) SetReportImage(reportID int64, contentType string, data []byte) error {
	_, err := s.db.Exec(`
		INSERT INTO report_images (report_id, content_type, bytes, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(report_id) DO UPDATE SET
			content_type = excluded.content_type,
			bytes = excluded.bytes,
			created_at = excluded.created_at
	`, reportID, contentType, data, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("set report image: %w", err)
	}
	return nil
}

func (s *PostgresStore) GetReportImage(reportID int64) (ReportImage, bool, error) {
	var img ReportImage
	err := s.db.QueryRow(
		`SELECT content_type, bytes FROM report_images WHERE report_id = $1`, reportID,
	).Scan(&img.ContentType, &img.Bytes)
	if err == sql.ErrNoRows {
		return ReportImage{}, false, nil
	}
	if err != nil {
		return ReportImage{}, false, fmt.Errorf("get report image: %w", err)
	}
	return img, true, nil
}
