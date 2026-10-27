// Package main implements an offline, synthetic replay harness for the FSRS
// and SM-2 schedulers. It is a comparison tool, not production code. Nothing
// imports it, and it never reads a database or a network.
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/chuma-beep/mathua/internal/fsrs"
	"github.com/chuma-beep/mathua/internal/scheduler"
)

const reportSchemaVersion = 1

// Fixture is one synthetic review history. Ratings drive FSRS. Quality drives
// SM-2. Keeping both explicit lets a fixture describe the same review under
// each model without pretending the two scales are identical.
type Fixture struct {
	ID                string   `json:"id"`
	Description       string   `json:"description"`
	InitialStability  *float64 `json:"initial_stability,omitempty"`
	InitialDifficulty *float64 `json:"initial_difficulty,omitempty"`
	Reviews           []Review `json:"reviews"`
}

type Review struct {
	GapDays int64 `json:"gap_days"`
	Rating  int   `json:"rating"`
	Quality int   `json:"quality"`
}

type StrategyResult struct {
	Name            string  `json:"name"`
	FinalDueDay     int64   `json:"final_due_day"`
	FinalInterval   int     `json:"final_interval"`
	FinalStability  float64 `json:"final_stability,omitempty"`
	FinalDifficulty float64 `json:"final_difficulty,omitempty"`
	MaxStability    float64 `json:"max_stability,omitempty"`
	CapHits         int     `json:"cap_hits,omitempty"`
}

type FixtureReport struct {
	ID               string         `json:"id"`
	Description      string         `json:"description"`
	SM2              StrategyResult `json:"sm2"`
	FSRSCapped       StrategyResult `json:"fsrs_capped"`
	FSRSCorrected    StrategyResult `json:"fsrs_corrected"`
	StabilityDelta   float64        `json:"stability_delta"`
	IntervalDelta    int            `json:"interval_delta"`
	DueDayDelta      int64          `json:"due_day_delta"`
	SM2IntervalDelta int            `json:"sm2_interval_delta"`
	SM2DueDayDelta   int64          `json:"sm2_due_day_delta"`
}

type Summary struct {
	Fixtures                        int     `json:"fixtures"`
	CappedCapHits                   int     `json:"capped_cap_hits"`
	FixturesWithStabilityDifference int     `json:"fixtures_with_stability_difference"`
	FixturesWithIntervalDifference  int     `json:"fixtures_with_interval_difference"`
	FixturesWithDueDayDifference    int     `json:"fixtures_with_due_day_difference"`
	MaxStabilityDelta               float64 `json:"max_stability_delta"`
}

type Report struct {
	SchemaVersion int             `json:"schema_version"`
	Model         string          `json:"model"`
	Fixtures      []FixtureReport `json:"fixtures"`
	Summary       Summary         `json:"summary"`
}

func loadFixtures(path string) ([]Fixture, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read fixtures: %w", err)
	}
	var fixtures []Fixture
	if err := json.Unmarshal(raw, &fixtures); err != nil {
		return nil, fmt.Errorf("decode fixtures: %w", err)
	}
	if len(fixtures) == 0 {
		return nil, fmt.Errorf("fixtures are empty")
	}
	return fixtures, nil
}

func Run(fixtures []Fixture) (Report, error) {
	report := Report{
		SchemaVersion: reportSchemaVersion,
		Model:         "fsrs-6-long-term",
		Fixtures:      make([]FixtureReport, 0, len(fixtures)),
	}
	seen := make(map[string]bool, len(fixtures))
	for _, fixture := range fixtures {
		if err := validateFixture(fixture, seen); err != nil {
			return Report{}, err
		}
		sm2, err := runSM2(fixture)
		if err != nil {
			return Report{}, err
		}
		capped, err := runFSRS(fixture, true)
		if err != nil {
			return Report{}, err
		}
		corrected, err := runFSRS(fixture, false)
		if err != nil {
			return Report{}, err
		}
		row := FixtureReport{
			ID:               fixture.ID,
			Description:      fixture.Description,
			SM2:              sm2,
			FSRSCapped:       capped,
			FSRSCorrected:    corrected,
			StabilityDelta:   corrected.FinalStability - capped.FinalStability,
			IntervalDelta:    corrected.FinalInterval - capped.FinalInterval,
			DueDayDelta:      corrected.FinalDueDay - capped.FinalDueDay,
			SM2IntervalDelta: corrected.FinalInterval - sm2.FinalInterval,
			SM2DueDayDelta:   corrected.FinalDueDay - sm2.FinalDueDay,
		}
		report.Fixtures = append(report.Fixtures, row)
		report.Summary.CappedCapHits += capped.CapHits
		if row.StabilityDelta != 0 {
			report.Summary.FixturesWithStabilityDifference++
		}
		if row.IntervalDelta != 0 {
			report.Summary.FixturesWithIntervalDifference++
		}
		if row.DueDayDelta != 0 {
			report.Summary.FixturesWithDueDayDifference++
		}
		if math.Abs(row.StabilityDelta) > report.Summary.MaxStabilityDelta {
			report.Summary.MaxStabilityDelta = math.Abs(row.StabilityDelta)
		}
	}
	report.Summary.Fixtures = len(fixtures)
	return report, nil
}

func validateFixture(f Fixture, seen map[string]bool) error {
	if f.ID == "" {
		return fmt.Errorf("fixture has no id")
	}
	if seen[f.ID] {
		return fmt.Errorf("duplicate fixture id %q", f.ID)
	}
	seen[f.ID] = true
	if len(f.Reviews) == 0 {
		return fmt.Errorf("fixture %q has no reviews", f.ID)
	}
	for i, r := range f.Reviews {
		if r.GapDays < 0 {
			return fmt.Errorf("fixture %q review %d has a negative gap", f.ID, i)
		}
		if r.Rating < 1 || r.Rating > 4 {
			return fmt.Errorf("fixture %q review %d has an invalid rating", f.ID, i)
		}
		if r.Quality < 0 || r.Quality > 5 {
			return fmt.Errorf("fixture %q review %d has an invalid quality", f.ID, i)
		}
	}
	for name, value := range map[string]*float64{
		"initial_stability":  f.InitialStability,
		"initial_difficulty": f.InitialDifficulty,
	} {
		if value == nil {
			continue
		}
		if math.IsNaN(*value) || math.IsInf(*value, 0) || *value <= 0 {
			return fmt.Errorf("fixture %q has an invalid %s", f.ID, name)
		}
	}
	return nil
}

func runSM2(f Fixture) (StrategyResult, error) {
	state := scheduler.DefaultSM2()
	var day int64
	for _, review := range f.Reviews {
		day += review.GapDays
		state = scheduler.ComputeSM2(state, review.Quality)
	}
	return StrategyResult{
		Name:          "sm2",
		FinalDueDay:   day + int64(state.Interval),
		FinalInterval: state.Interval,
	}, nil
}

func runFSRS(f Fixture, capped bool) (StrategyResult, error) {
	p := fsrs.DefaultParameters()
	first := f.Reviews[0]
	rating := fsrs.Rating(first.Rating)

	stability := p.InitialStability(rating)
	if f.InitialStability != nil {
		stability = *f.InitialStability
	}
	difficulty := p.InitialDifficulty(rating)
	if f.InitialDifficulty != nil {
		difficulty = *f.InitialDifficulty
	}

	result := StrategyResult{Name: "fsrs_corrected"}
	if capped {
		result.Name = "fsrs_capped"
		if stability > fsrs.MaxStabilityDays {
			stability = fsrs.MaxStabilityDays
			result.CapHits++
		}
	}

	var day int64
	var interval int
	var due int64
	for i, review := range f.Reviews {
		day += review.GapDays
		rating = fsrs.Rating(review.Rating)
		if i > 0 {
			retrievability := p.Retrievability(float64(review.GapDays), stability)
			stability = p.NextStability(stability, difficulty, retrievability, rating)
			if capped && stability > fsrs.MaxStabilityDays {
				stability = fsrs.MaxStabilityDays
				result.CapHits++
			}
			difficulty = p.NextDifficulty(difficulty, rating)
		}
		interval = p.Interval(stability)
		due = day + int64(interval)
		if stability > result.MaxStability {
			result.MaxStability = stability
		}
	}
	result.FinalDueDay = due
	result.FinalInterval = interval
	result.FinalStability = stability
	result.FinalDifficulty = difficulty
	return result, nil
}

func writeReport(path string, report Report) error {
	raw, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("encode report: %w", err)
	}
	raw = append(raw, '\n')
	if path == "" {
		_, err = os.Stdout.Write(raw)
		return err
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return fmt.Errorf("write report: %w", err)
	}
	return nil
}
