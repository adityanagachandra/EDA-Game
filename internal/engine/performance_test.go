package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/all"
	"github.com/chuma-beep/mathua/internal/mastery"
	"github.com/chuma-beep/mathua/internal/storage"
)

// TestTimedQuizPerformance is a targeted, opt-in load test for the durable
// timed-quiz and review submission paths. It is skipped in normal CI because
// it seeds a production-shaped database and takes minutes. Run it with:
//
//	MATHUA_BENCH=1 go test ./internal/engine -run TestTimedQuizPerformance -v
//
// Set MATHUA_BENCH_DSN to a disposable PostgreSQL database to run the same
// workload against PostgreSQL instead of SQLite.
func TestTimedQuizPerformance(t *testing.T) {
	if os.Getenv("MATHUA_BENCH") == "" {
		t.Skip("set MATHUA_BENCH=1 to run the targeted performance test")
	}
	students := envInt("MATHUA_BENCH_STUDENTS", 1200)
	workers := envInt("MATHUA_BENCH_WORKERS", 16)
	if students < 160 {
		t.Fatalf("MATHUA_BENCH_STUDENTS must be at least 160, got %d", students)
	}
	if workers < 1 {
		t.Fatalf("MATHUA_BENCH_WORKERS must be positive, got %d", workers)
	}

	dag, err := concepts.LoadDir("../../data/concepts")
	if err != nil {
		t.Fatalf("load corpus: %v", err)
	}
	reg := all.Registry()
	store, engineName, closeStore := benchStore(t)
	defer closeStore()
	e := New(store, dag, reg, nil, nil)
	tq := store.(storage.TimedQuizStore)

	start := time.Now()
	seed := seedBenchmarkData(t, e, dag, reg, students)
	t.Logf("dataset: engine=%s students=%d concepts=%d seeded_progress=%d baseline_xp=200 seed_time=%s workers=%d",
		engineName, students, len(dag.Order()), seed.progressRows, time.Since(start).Round(time.Millisecond), workers)

	const perGroup = 160
	if students < perGroup*6 {
		t.Fatalf("need at least %d students for six groups, got %d", perGroup*6, students)
	}

	// Start: one fresh student per run.
	startStudents := seed.timed[0:perGroup]
	startResult := runWorkload(workers, perGroup, func(i int) error {
		_, _, err := e.StartTimedQuiz(startStudents[i].ID)
		return err
	})
	report(t, "timed quiz start", startResult)

	// Answer: pre-start runs, then claim, grade and update one item per job.
	answerStudents := seed.timed[perGroup : perGroup*2]
	var answerJobs []benchJob
	for _, st := range answerStudents {
		run, items, err := e.StartTimedQuiz(st.ID)
		if err != nil {
			t.Fatalf("pre-start answer run: %v", err)
		}
		for _, item := range items {
			item := item
			answerJobs = append(answerJobs, benchJob{runID: run.ID, item: item, studentID: st.ID})
		}
	}
	answerResult := runJobList(workers, answerJobs, func(j benchJob) error {
		submissionID := "bench-answer-" + j.item.ID
		ok, err := tq.ClaimTimedQuizItem(j.item.ID, j.runID, submissionID)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("item %s was already claimed", j.item.ID)
		}
		res, err := e.SubmitTimedQuizAnswer(j.studentID, j.item.ConceptID, j.item.ExpectedAnswer, j.item.ExpectedAnswer, 5, j.item.Question, j.item.Difficulty, submissionID)
		if err != nil {
			return err
		}
		_, err = tq.UpdateTimedQuizItem(j.item.ID, j.runID, j.item.ExpectedAnswer, submissionID, res.Correct, 5, res.XP, "answered")
		return err
	})
	report(t, "timed quiz answer", answerResult)

	// Skip: pre-start runs, then claim and mark skipped.
	skipStudents := seed.timed[perGroup*2 : perGroup*3]
	var skipJobs []benchJob
	for _, st := range skipStudents {
		run, items, err := e.StartTimedQuiz(st.ID)
		if err != nil {
			t.Fatalf("pre-start skip run: %v", err)
		}
		for _, item := range items {
			item := item
			skipJobs = append(skipJobs, benchJob{runID: run.ID, item: item, studentID: st.ID})
		}
	}
	skipResult := runJobList(workers, skipJobs, func(j benchJob) error {
		submissionID := "bench-skip-" + j.item.ID
		ok, err := tq.ClaimTimedQuizItem(j.item.ID, j.runID, submissionID)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("item %s was already claimed", j.item.ID)
		}
		_, err = tq.UpdateTimedQuizItem(j.item.ID, j.runID, "", submissionID, false, 0, 0, "skipped")
		return err
	})
	report(t, "timed quiz skip", skipResult)

	// Finalize: pre-start, mark every item skipped, then finalize once.
	finalizeStudents := seed.timed[perGroup*3 : perGroup*4]
	var finalizeRuns []*storage.TimedQuizRun
	for _, st := range finalizeStudents {
		run, items, err := e.StartTimedQuiz(st.ID)
		if err != nil {
			t.Fatalf("pre-start finalize run: %v", err)
		}
		for _, item := range items {
			submissionID := "bench-finalize-" + item.ID
			if _, err := tq.ClaimTimedQuizItem(item.ID, run.ID, submissionID); err != nil {
				t.Fatalf("claim finalize item: %v", err)
			}
			if _, err := tq.UpdateTimedQuizItem(item.ID, run.ID, "", submissionID, false, 0, 0, "skipped"); err != nil {
				t.Fatalf("skip finalize item: %v", err)
			}
		}
		finalizeRuns = append(finalizeRuns, run)
	}
	finalizeResult := runWorkload(workers, len(finalizeRuns), func(i int) error {
		_, err := tq.FinalizeTimedQuizRun(finalizeRuns[i].ID, finalizeRuns[i].StudentID, time.Now().UTC())
		return err
	})
	report(t, "timed quiz finalize", finalizeResult)

	// Review submission: pre-create one due review session per student, then
	// measure only SubmitAnswer, which owns the atomic review write.
	reviewStudents := seed.review[0:perGroup]
	var reviewJobs []reviewJob
	for _, st := range reviewStudents {
		sess, err := store.CreateSession(st.ID)
		if err != nil {
			t.Fatalf("create review session: %v", err)
		}
		q, err := e.NextReviewQuestion(sess.ID, st.ID)
		if err != nil || q == nil {
			t.Fatalf("prime review question: %v, %v", err, q)
		}
		e.mu.Lock()
		expected := e.sessions[sess.ID].expectedAnswer
		e.mu.Unlock()
		reviewJobs = append(reviewJobs, reviewJob{sessionID: sess.ID, studentID: st.ID, attemptID: q.AttemptID, answer: expected})
	}
	reviewResult := runReviewList(workers, reviewJobs, func(j reviewJob) error {
		_, err := e.SubmitAnswer(j.sessionID, j.studentID, j.attemptID, j.answer, 5)
		return err
	})
	report(t, "review submission", reviewResult)

	t.Logf("dataset detail: progress_rows=%d timed_students=%d review_students=%d", seed.progressRows, len(seed.timed), len(seed.review))
}

type benchSeed struct {
	timed        []*storage.Student
	review       []*storage.Student
	progressRows int
}

func benchStore(t *testing.T) (storage.Repository, string, func()) {
	t.Helper()
	if dsn := os.Getenv("MATHUA_BENCH_DSN"); dsn != "" {
		store, err := storage.NewPostgresStore(dsn)
		if err != nil {
			t.Fatalf("open postgres: %v", err)
		}
		return store, "postgres", func() { _ = store.Close() }
	}
	path := filepath.Join(t.TempDir(), "bench.db")
	store, err := storage.NewSQLiteStore(path)
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return store, "sqlite", func() {
		_ = store.Close()
		if info, err := os.Stat(path); err == nil {
			t.Logf("sqlite database size: %.1f MB", float64(info.Size())/(1024*1024))
		}
	}
}

func seedBenchmarkData(t *testing.T, e *Engine, dag *concepts.DAG, reg *generator.Registry, students int) benchSeed {
	t.Helper()
	order := dag.Order()
	if len(order) < 80 {
		t.Fatalf("benchmark needs at least 80 concepts, got %d", len(order))
	}
	// The first 60 concepts form a mastered prefix so the next concepts have
	// satisfied prerequisites. The next 20 are the quiz candidates.
	var practice []*concepts.Concept
	for _, c := range order[60:] {
		if len(practice) >= 20 {
			break
		}
		if _, err := reg.Get(c.ID); err == nil {
			practice = append(practice, c)
		}
	}
	if len(practice) < 10 {
		t.Fatalf("benchmark needs at least 10 practice candidates, got %d", len(practice))
	}
	now := time.Now().UTC()
	last := now.AddDate(0, 0, -3)
	reviewLast := now.AddDate(0, 0, -30)
	reviewDue := now.Add(-time.Hour)
	future := now.AddDate(0, 0, 90)
	stability := 100.0
	reviewStability := 10.0
	difficulty := 5.0

	seed := benchSeed{}
	for i := 0; i < students; i++ {
		st, err := e.CreateStudent(fmt.Sprintf("bench-%04d", i))
		if err != nil {
			t.Fatalf("create student %d: %v", i, err)
		}
		quizStudent := i < students*5/6
		if quizStudent {
			if err := e.repo.AwardXP(st.ID, 200, "bench", "bench-xp-"+st.ID); err != nil {
				t.Fatalf("award benchmark xp: %v", err)
			}
			for _, c := range order[:60] {
				p := &storage.ConceptProgress{
					StudentID: st.ID, ConceptID: c.ID, Status: string(mastery.StatusMastered),
					Streak: 1, BestStreak: 5, AvgResponseTime: 6, Attempts: 6,
					LastAttempted: &last, LastReviewed: &reviewLast, NextReviewDue: &future,
					SM2Repetitions: 4, SM2Interval: 90, SM2EFactor: 2.5,
					FSRSStability: &stability, FSRSDifficulty: &difficulty, FSRSLastReview: &reviewLast,
				}
				if err := e.repo.UpsertProgress(p); err != nil {
					t.Fatalf("seed mastered progress: %v", err)
				}
				seed.progressRows++
			}
			for _, c := range practice {
				p := &storage.ConceptProgress{
					StudentID: st.ID, ConceptID: c.ID, Status: string(mastery.StatusLearning),
					Streak: 1, BestStreak: 1, AvgResponseTime: 8, Attempts: 6,
					LastAttempted: &last, SM2Repetitions: 1, SM2Interval: 1, SM2EFactor: 2.5,
					WeaknessScore: 0.6,
				}
				if err := e.repo.UpsertProgress(p); err != nil {
					t.Fatalf("seed practice progress: %v", err)
				}
				seed.progressRows++
			}
			seed.timed = append(seed.timed, st)
		} else {
			// One mastered, overdue root concept. A root has no prerequisite
			// that could make NextReview skip it.
			var root *concepts.Concept
			for _, c := range order {
				if len(c.Prerequisites) == 0 {
					if _, err := reg.Get(c.ID); err == nil {
						root = c
						break
					}
				}
			}
			if root == nil {
				t.Fatal("no root concept with a generator")
			}
			p := &storage.ConceptProgress{
				StudentID: st.ID, ConceptID: root.ID, Status: string(mastery.StatusMastered),
				Streak: 1, BestStreak: 5, AvgResponseTime: 6, Attempts: 6,
				LastAttempted: &reviewLast, LastReviewed: &reviewLast, NextReviewDue: &reviewDue,
				SM2Repetitions: 4, SM2Interval: 30, SM2EFactor: 2.5,
				FSRSStability: &reviewStability, FSRSDifficulty: &difficulty, FSRSLastReview: &reviewLast,
			}
			if err := e.repo.UpsertProgress(p); err != nil {
				t.Fatalf("seed review progress: %v", err)
			}
			seed.progressRows++
			seed.review = append(seed.review, st)
		}
	}
	return seed
}

type benchJob struct {
	index     int
	runID     string
	studentID string
	item      storage.TimedQuizItem
}

type reviewJob struct {
	sessionID string
	studentID string
	attemptID string
	answer    string
}

type latencyResult struct {
	name      string
	n         int
	errors    int
	locks     int
	wall      time.Duration
	cpu       time.Duration
	rssMB     float64
	durations []time.Duration
}

func runWorkload(workers, jobs int, fn func(i int) error) latencyResult {
	return runJobList(workers, makeIndexJobs(jobs), func(j benchJob) error { return fn(j.index) })
}

func makeIndexJobs(n int) []benchJob {
	out := make([]benchJob, n)
	for i := range out {
		out[i].index = i
	}
	return out
}

func runJobList(workers int, jobs []benchJob, fn func(benchJob) error) latencyResult {
	return collectLatencies(workers, len(jobs), func(i int) error { return fn(jobs[i]) })
}

func runReviewList(workers int, jobs []reviewJob, fn func(reviewJob) error) latencyResult {
	return collectLatencies(workers, len(jobs), func(i int) error { return fn(jobs[i]) })
}

func collectLatencies(workers, n int, fn func(i int) error) latencyResult {
	beforeUser, beforeSys := cpuTime()
	beforeRSS := rssMB()
	start := time.Now()
	var mu sync.Mutex
	res := latencyResult{n: n, durations: make([]time.Duration, 0, n)}
	ch := make(chan int)
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range ch {
				t0 := time.Now()
				err := fn(i)
				d := time.Since(t0)
				mu.Lock()
				res.durations = append(res.durations, d)
				if err != nil {
					res.errors++
					if isLockError(err) {
						res.locks++
					}
				}
				mu.Unlock()
			}
		}()
	}
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	res.wall = time.Since(start)
	afterUser, afterSys := cpuTime()
	res.cpu = (afterUser - beforeUser) + (afterSys - beforeSys)
	res.rssMB = rssMB()
	if beforeRSS > res.rssMB {
		res.rssMB = beforeRSS
	}
	sort.Slice(res.durations, func(i, j int) bool { return res.durations[i] < res.durations[j] })
	return res
}

func report(t *testing.T, name string, r latencyResult) {
	t.Helper()
	if r.n == 0 {
		t.Fatalf("%s: no samples", name)
	}
	p50 := percentile(r.durations, 0.50)
	p95 := percentile(r.durations, 0.95)
	p99 := percentile(r.durations, 0.99)
	throughput := float64(r.n) / r.wall.Seconds()
	cpuPct := 0.0
	if r.wall > 0 {
		cpuPct = float64(r.cpu) / float64(r.wall) * 100
	}
	t.Logf("%s: n=%d throughput=%.1f ops/s p50=%.1fms p95=%.1fms p99=%.1fms errors=%d lock_errors=%d cpu=%.0f%% rss=%.0fMB wall=%s",
		name, r.n, throughput, ms(p50), ms(p95), ms(p99), r.errors, r.locks, cpuPct, r.rssMB, r.wall.Round(time.Millisecond))
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx]
}

func ms(d time.Duration) float64 { return float64(d) / float64(time.Millisecond) }

func isLockError(err error) bool {
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "locked") || strings.Contains(s, "busy") || strings.Contains(s, "deadlock") || strings.Contains(s, "timeout")
}

func cpuTime() (time.Duration, time.Duration) {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return 0, 0
	}
	return time.Duration(ru.Utime.Sec)*time.Second + time.Duration(ru.Utime.Usec)*time.Microsecond,
		time.Duration(ru.Stime.Sec)*time.Second + time.Duration(ru.Stime.Usec)*time.Microsecond
}

func rssMB() float64 {
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		return float64(ms.Sys) / (1024 * 1024)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if kb, err := strconv.ParseFloat(fields[1], 64); err == nil {
					return kb / 1024
				}
			}
		}
	}
	return 0
}

func envInt(name string, fallback int) int {
	raw := os.Getenv(name)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}
