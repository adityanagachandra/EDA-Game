package edakit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/grader"
)

// MinVariety is the fewest distinct situations a concept must supply in one contract sweep.
// The sweep draws 120 problems; a concept below this floor would repeat inside one session.
const MinVariety = 10

var (
	numericAnswer  = regexp.MustCompile(`^[\d\s.,/+\-]+$`)
	doubledCommand = regexp.MustCompile(`\\\\[A-Za-z]{2,}`)
	clicheSubgoal  = regexp.MustCompile(`(?i)^follow the worked|^apply the rule$|^check the result$|^understand\b|^review |^work the example$|^verify$`)
	headingLine    = regexp.MustCompile(`(?m)^## (.+?)\s*$`)
)

type conceptRow struct {
	ID               string   `json:"id"`
	Label            string   `json:"label"`
	Domain           string   `json:"domain"`
	Subdomain        string   `json:"subdomain"`
	GradingType      string   `json:"grading_type"`
	Prerequisites    []string `json:"prerequisites"`
	KeyPrerequisites []string `json:"key_prerequisites"`
	Encompasses      []string `json:"encompasses"`
	MasteryThreshold struct {
		Streak         int     `json:"streak"`
		AvgTimeSeconds float64 `json:"avg_time_seconds"`
	} `json:"mastery_threshold"`
}

type kpRow struct {
	Label    string   `json:"label"`
	Section  string   `json:"section"`
	Subgoals []string `json:"subgoals"`
}

// RepoRoot walks up from the working directory to the directory holding go.mod.
func RepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("go.mod not found above the test directory")
		}
		dir = parent
	}
}

// CheckDomain is the contract every EDA domain package runs from its own test. It loads
// data/concepts/<domain>.json, registers the package's generators onto a fresh registry and
// checks the generators, the lessons and the shards the way the corpus-wide suites will, so a
// domain can be finished and proven in isolation.
func CheckDomain(t *testing.T, register func(*generator.Registry), domain string) {
	t.Helper()
	root := RepoRoot(t)
	raw, err := os.ReadFile(filepath.Join(root, "data", "concepts", domain+".json"))
	if err != nil {
		t.Fatalf("read concept file: %v", err)
	}
	var rows []conceptRow
	if err := json.Unmarshal(raw, &rows); err != nil {
		t.Fatalf("parse concept file: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("concept file %s.json is empty", domain)
	}

	reg := generator.NewRegistry()
	register(reg)
	if d := reg.Duplicates(); len(d) > 0 {
		t.Errorf("concepts registered twice: %v", d)
	}
	declared := map[string]bool{}
	for _, c := range rows {
		declared[c.ID] = true
	}
	for _, id := range reg.Concepts() {
		if !declared[id] {
			t.Errorf("%s: registered but not declared in %s.json", id, domain)
		}
	}

	router := grader.NewRouter()
	for _, c := range rows {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			checkMetadata(t, c, domain)
			checkLesson(t, root, c)
			if _, err := reg.Get(c.ID); err != nil {
				t.Fatalf("no generator registered")
			}
			checkGenerator(t, reg, router, c)
		})
	}
}

func checkMetadata(t *testing.T, c conceptRow, domain string) {
	t.Helper()
	if c.Domain != domain {
		t.Errorf("domain %q, want %q", c.Domain, domain)
	}
	if !strings.HasPrefix(c.Subdomain, domain+".") {
		t.Errorf("subdomain %q must start with %q", c.Subdomain, domain+".")
	}
	if c.Label == "" || strings.ContainsAny(c.Label, "$\\") {
		t.Errorf("label %q is empty or carries markup", c.Label)
	}
	switch c.GradingType {
	case "numeric", "multiple_choice":
	default:
		t.Errorf("grading_type %q: EDA concepts use numeric or multiple_choice", c.GradingType)
	}
	if len(c.Prerequisites) == 0 {
		t.Errorf("no prerequisites: every EDA concept builds on something")
	}
	if len(c.KeyPrerequisites) == 0 {
		t.Errorf("no key_prerequisites: name the prerequisite this concept leans on hardest")
	}
	pre := map[string]bool{}
	for _, p := range c.Prerequisites {
		if p == c.ID {
			t.Errorf("concept lists itself as a prerequisite")
		}
		pre[p] = true
	}
	for _, k := range c.KeyPrerequisites {
		if !pre[k] {
			t.Errorf("key_prerequisite %q is not a prerequisite", k)
		}
	}
	for _, e := range c.Encompasses {
		if e == c.ID {
			t.Errorf("concept encompasses itself")
		}
	}
	if c.MasteryThreshold.Streak <= 0 || c.MasteryThreshold.AvgTimeSeconds <= 0 || c.MasteryThreshold.AvgTimeSeconds > 300 {
		t.Errorf("mastery_threshold %+v is out of range", c.MasteryThreshold)
	}
}

func checkLesson(t *testing.T, root string, c conceptRow) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(root, "data", "lessons", "teaching", c.ID+".md"))
	if err != nil {
		t.Errorf("lesson: %v", err)
		return
	}
	text := string(body)
	if !strings.HasPrefix(text, "# ") {
		t.Errorf("lesson must open with a '# Title' line")
	}
	if first := strings.SplitN(text, "\n", 2)[0]; strings.ContainsAny(first, "$\\") || strings.Contains(first, "](") {
		t.Errorf("lesson title carries markup: %q", first)
	}
	if n := len(regexp.MustCompile(`[A-Za-z]{2,}`).FindAllString(text, -1)); n < 250 {
		t.Errorf("lesson has %d words, want at least 250", n)
	}
	if strings.Count(text, "$")%2 != 0 {
		t.Errorf("lesson has an unbalanced $ delimiter")
	}
	checkGlyphs(t, "lesson", text)
	headings := map[string]bool{}
	for _, m := range headingLine.FindAllStringSubmatch(text, -1) {
		headings[m[1]] = true
	}

	shard, err := os.ReadFile(filepath.Join(root, "data", "lessons", "kp", c.ID+".json"))
	if err != nil {
		t.Errorf("kp shard: %v", err)
		return
	}
	var kps []kpRow
	if err := json.Unmarshal(shard, &kps); err != nil {
		t.Errorf("kp shard: %v", err)
		return
	}
	if len(kps) != 3 {
		t.Errorf("kp shard has %d knowledge points, want 3", len(kps))
	}
	for i, kp := range kps {
		if kp.Label == "" || kp.Section == "" {
			t.Errorf("kp %d: empty label or section", i+1)
		}
		if !headings[kp.Section] {
			t.Errorf("kp %d section %q is not a '## ' heading of the lesson", i+1, kp.Section)
		}
		if len(kp.Subgoals) != 3 {
			t.Errorf("kp %d has %d subgoals, want 3", i+1, len(kp.Subgoals))
		}
		for _, sg := range kp.Subgoals {
			if clicheSubgoal.MatchString(sg) {
				t.Errorf("kp %d subgoal %q is filler; name the move", i+1, sg)
			}
		}
	}
}

func checkGlyphs(t *testing.T, where, s string) {
	t.Helper()
	for _, r := range s {
		if r > 0xFFFF || (r >= 0x2460 && r <= 0x24FF) || (r >= 0x2600 && r <= 0x27BF && !strings.ContainsRune("✓✗✔✕✖✦✧", r)) {
			t.Errorf("%s contains pictograph %q (AGENTS.md rule 1)", where, string(r))
			return
		}
	}
	if !utf8.ValidString(s) {
		t.Errorf("%s is not valid UTF-8", where)
	}
}

// StressDraws is how many problems the hang check draws per concept. A redraw loop that only
// re-rolls one operand can be unsatisfiable for rare combinations of the others: one such loop
// in the geometry domain hung about once in seven thousand draws. Twenty thousand draws find a
// trap of that rarity with near certainty, and take well under a second for a healthy generator.
const StressDraws = 20000

// checkTerminates draws StressDraws problems under a watchdog and fails if they do not finish.
func checkTerminates(t *testing.T, gen generator.Generator) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		diffs := []float64{0.2, 0.3, 0.5, 0.8, 1.0}
		for i := 0; i < StressDraws; i++ {
			gen.Generate(generator.GeneratorContext{Difficulty: diffs[i%len(diffs)], KPOrdinal: i % 4})
		}
	}()
	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatalf("%d draws did not finish in 30s: a redraw loop can spin forever. Bound every redraw loop with generator.MaxDrawAttempts and redraw every operand its condition depends on", StressDraws)
	}
}

func checkGenerator(t *testing.T, reg *generator.Registry, router *grader.Router, c conceptRow) {
	t.Helper()
	gen, _ := reg.Get(c.ID)
	checkTerminates(t, gen)
	custom, hasCustom := gen.(interface {
		Grade(string, string) grader.Result
	})
	seen := map[string]bool{}
	failures := 0
	for kp := 0; kp <= 3; kp++ {
		for _, diff := range []float64{0.3, 0.6, 0.9} {
			for i := 0; i < 10; i++ {
				p, err := reg.GenerateContext(c.ID, generator.GeneratorContext{Difficulty: diff, KPOrdinal: kp})
				if err != nil {
					t.Fatalf("generate: %v", err)
				}
				if failures > 5 {
					return
				}
				if strings.TrimSpace(p.Question) == "" || strings.TrimSpace(p.Answer) == "" {
					t.Errorf("empty question or answer: %+v", p)
					failures++
					continue
				}
				if strings.Count(p.Question, "$")%2 != 0 {
					t.Errorf("unbalanced $ in question: %q", p.Question)
					failures++
				}
				if doubledCommand.MatchString(p.Question) {
					t.Errorf("doubled LaTeX command in question: %q", p.Question)
					failures++
				}
				checkGlyphs(t, "question", p.Question)
				var res grader.Result
				if hasCustom {
					res = custom.Grade(p.Answer, p.Answer)
				} else {
					res = router.Grade(grader.GradingType(c.GradingType), p.Answer, p.Answer)
				}
				if !res.Correct {
					t.Errorf("answer %q does not grade as correct against itself: %s", p.Answer, res.Feedback)
					failures++
				}
				opts := generator.ChoiceOptions(p.Question)
				switch c.GradingType {
				case "numeric":
					if !numericAnswer.MatchString(p.Answer) {
						t.Errorf("numeric answer %q is not keypad-typable", p.Answer)
						failures++
					}
					if len(opts) > 0 {
						t.Errorf("numeric concept served an option block")
						failures++
					}
				case "multiple_choice":
					if len(opts) < 3 {
						t.Errorf("multiple-choice question has %d options, want at least 3: %q", len(opts), p.Question)
						failures++
					}
					found := false
					for _, o := range opts {
						if o == p.Answer {
							found = true
							continue
						}
						wrong := router.Grade(grader.GradingType(c.GradingType), p.Answer, o)
						if hasCustom {
							wrong = custom.Grade(p.Answer, o)
						}
						if wrong.Correct {
							t.Errorf("distractor %q grades as correct for answer %q", o, p.Answer)
							failures++
						}
					}
					if !found {
						t.Errorf("answer %q is not among the options", p.Answer)
						failures++
					}
				}
				seen[canonical(p.Question)] = true
			}
		}
	}
	if len(seen) < MinVariety {
		t.Errorf("only %d distinct situations in 120 draws, want at least %d", len(seen), MinVariety)
	}
}

func canonical(q string) string {
	core := q
	if i := strings.Index(q, "\n\nA) "); i >= 0 {
		core = q[:i]
	}
	opts := generator.ChoiceOptions(q)
	texts := make([]string, 0, len(opts))
	for _, o := range opts {
		texts = append(texts, o)
	}
	sort.Strings(texts)
	return core + "\n" + strings.Join(texts, "\n")
}
