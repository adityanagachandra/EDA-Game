package main

import (
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/fsrs"
)

// TestReplayPackageHasNoDatabaseOrMutationImports enforces the boundary that
// makes the harness safe to run offline. The comparison core may use the pure
// FSRS and SM-2 packages. It must not reach storage, the engine, the server or
// a network or clock.
func TestReplayPackageHasNoDatabaseOrMutationImports(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package directory: %v", err)
	}
	forbidden := []string{
		"database/",
		"net/",
		"internal/storage",
		"internal/engine",
		"internal/server",
		"internal/auth",
		"internal/ai",
		"os/exec",
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		for _, spec := range file.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				t.Fatalf("unquote import in %s: %v", name, err)
			}
			if path == "net" || path == "time" {
				t.Errorf("%s imports %q, which breaks offline determinism", name, path)
			}
			for _, bad := range forbidden {
				if strings.Contains(path, bad) {
					t.Errorf("%s imports %q, which breaks the replay boundary", name, path)
				}
			}
		}
	}
}

func TestFixturesProduceAStableReport(t *testing.T) {
	fixtures, err := loadFixtures("fixtures.json")
	if err != nil {
		t.Fatalf("load fixtures: %v", err)
	}
	first, err := Run(fixtures)
	if err != nil {
		t.Fatalf("first run: %v", err)
	}
	second, err := Run(fixtures)
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("replay report is not deterministic")
	}
	if first.SchemaVersion != reportSchemaVersion {
		t.Fatalf("schema version = %d, want %d", first.SchemaVersion, reportSchemaVersion)
	}
	if first.Summary.Fixtures != len(fixtures) {
		t.Fatalf("summary fixtures = %d, want %d", first.Summary.Fixtures, len(fixtures))
	}
}

func TestCapFixtureSeparatesStoredStabilityFromInterval(t *testing.T) {
	fixtures, err := loadFixtures("fixtures.json")
	if err != nil {
		t.Fatalf("load fixtures: %v", err)
	}
	report, err := Run(fixtures)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	row := findFixture(t, report, "cap_boundary")

	if row.FSRSCapped.FinalStability != fsrs.MaxStabilityDays {
		t.Errorf("capped stability = %v, want %v", row.FSRSCapped.FinalStability, fsrs.MaxStabilityDays)
	}
	if row.FSRSCorrected.FinalStability <= fsrs.MaxStabilityDays {
		t.Errorf("corrected stability = %v, want above %v", row.FSRSCorrected.FinalStability, fsrs.MaxStabilityDays)
	}
	if row.FSRSCapped.FinalInterval != fsrs.MaxIntervalDays {
		t.Errorf("capped interval = %d, want %d", row.FSRSCapped.FinalInterval, fsrs.MaxIntervalDays)
	}
	if row.FSRSCorrected.FinalInterval != fsrs.MaxIntervalDays {
		t.Errorf("corrected interval = %d, want %d", row.FSRSCorrected.FinalInterval, fsrs.MaxIntervalDays)
	}
	if row.DueDayDelta != 0 {
		t.Errorf("due day delta = %d, want 0 for this fixture", row.DueDayDelta)
	}
	if row.FSRSCapped.CapHits == 0 {
		t.Error("capped fixture did not record a cap hit")
	}
}

func TestAboveCapInputIsNotReduced(t *testing.T) {
	fixtures, err := loadFixtures("fixtures.json")
	if err != nil {
		t.Fatalf("load fixtures: %v", err)
	}
	report, err := Run(fixtures)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	row := findFixture(t, report, "above_cap_input")

	if row.FSRSCapped.FinalStability > fsrs.MaxStabilityDays {
		t.Errorf("capped stability = %v, want <= %v", row.FSRSCapped.FinalStability, fsrs.MaxStabilityDays)
	}
	if row.FSRSCorrected.FinalStability <= fsrs.MaxStabilityDays {
		t.Errorf("corrected stability = %v, want above %v", row.FSRSCorrected.FinalStability, fsrs.MaxStabilityDays)
	}
	if row.FSRSCorrected.FinalInterval != fsrs.MaxIntervalDays {
		t.Errorf("corrected interval = %d, want %d", row.FSRSCorrected.FinalInterval, fsrs.MaxIntervalDays)
	}
}

func TestInvalidFixturesAreRejected(t *testing.T) {
	valid := Fixture{ID: "ok", Reviews: []Review{{GapDays: 0, Rating: 3, Quality: 4}}}
	cases := []struct {
		name    string
		fixture Fixture
	}{
		{"empty id", Fixture{Reviews: valid.Reviews}},
		{"no reviews", Fixture{ID: "empty"}},
		{"negative gap", Fixture{ID: "gap", Reviews: []Review{{GapDays: -1, Rating: 3, Quality: 4}}}},
		{"bad rating", Fixture{ID: "rating", Reviews: []Review{{Rating: 5, Quality: 4}}}},
		{"bad quality", Fixture{ID: "quality", Reviews: []Review{{Rating: 3, Quality: 6}}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Run([]Fixture{tc.fixture}); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
	seen := make(map[string]bool)
	if err := validateFixture(valid, seen); err != nil {
		t.Fatalf("valid fixture rejected: %v", err)
	}
	if err := validateFixture(valid, seen); err == nil {
		t.Fatal("duplicate fixture id accepted")
	}
}

func findFixture(t *testing.T, report Report, id string) FixtureReport {
	t.Helper()
	for _, row := range report.Fixtures {
		if row.ID == id {
			return row
		}
	}
	t.Fatalf("fixture %q not found", id)
	return FixtureReport{}
}
