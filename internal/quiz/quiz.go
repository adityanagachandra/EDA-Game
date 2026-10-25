package quiz

import (
	"sort"
	"sync"
	"time"

	"github.com/chuma-beep/mathua/internal/concepts"
	"github.com/chuma-beep/mathua/internal/generator"
)

// QuizSize is the number of concepts in one routine mastery check. Selection
// and Start share it so a caller cannot choose six and silently have the sixth
// dropped.
const QuizSize = 5

// Session is one run of the routine mastery check, behind the XP gate and on
// its own grading path (attemptID). The server owns admission and reuse: an
// anonymous run is created per request, while an identified learner gets the
// run already in progress instead of a second one.
type Session struct {
	sync.Mutex
	ID            string
	StudentID     string
	Order         []*concepts.Concept
	Index         int
	Attempts      []Attempt
	LastProblem   *generator.Problem
	LastConceptID string
	// QuestionStarted is when the pending question was served. The server owns
	// this clock: a refresh resumes the same start time, and an answer is
	// measured from it rather than from the client's report alone.
	QuestionStarted time.Time
	// Difficulties holds weakness→difficulty per concept (80% target, 0.3-1.0).
	Difficulties map[string]float64
}

type Attempt struct {
	ConceptID string
	Correct   bool
}

type Engine struct {
	dag      *concepts.DAG
	registry *generator.Registry
}

func NewEngine(dag *concepts.DAG, reg *generator.Registry) *Engine {
	return &Engine{dag: dag, registry: reg}
}

// Start accepts the concepts chosen for one quiz, capped at QuizSize.
// An empty selection falls back to the first concepts in DAG order, which is
// the anonymous demo path; identified learners are refused before they get
// here when there are no eligible concepts.
func (e *Engine) Start(concepts []*concepts.Concept) *Session {
	if len(concepts) == 0 {
		concepts = e.dag.Order()
	}
	if len(concepts) > QuizSize {
		concepts = concepts[:QuizSize]
	}
	return &Session{
		Order:        concepts,
		Difficulties: make(map[string]float64, len(concepts)),
	}
}

func (e *Engine) NextQuestion(s *Session) (*generator.Problem, string, error) {
	s.Lock()
	defer s.Unlock()
	if s.Index >= len(s.Order) {
		s.LastProblem = nil
		s.QuestionStarted = time.Time{}
		return nil, "", nil
	}
	c := s.Order[s.Index]
	// 80% difficulty target via weakness map (engine.DifficultyFor); fallback 0.8.
	difficulty := 0.8
	if d, ok := s.Difficulties[c.ID]; ok && d > 0 {
		difficulty = d
	}
	p, err := e.registry.Generate(c.ID, difficulty)
	if err != nil {
		return nil, "", err
	}
	s.LastProblem = &p
	s.LastConceptID = c.ID
	s.QuestionStarted = time.Now().UTC()
	return s.LastProblem, c.ID, nil
}

// QuestionElapsed returns seconds since the pending question was served. It
// is the server's own measurement, so a refresh or a client clock change
// cannot shorten it. A zero start time means the question predates this field;
// the caller should fall back to the client's report rather than inventing a
// duration.
func (s *Session) QuestionElapsed(now time.Time) float64 {
	s.Lock()
	defer s.Unlock()
	if s.QuestionStarted.IsZero() || s.LastProblem == nil {
		return 0
	}
	elapsed := now.Sub(s.QuestionStarted).Seconds()
	if elapsed < 0 {
		return 0
	}
	return elapsed
}

func (e *Engine) RecordAnswer(s *Session, conceptID string, correct bool) {
	s.Lock()
	defer s.Unlock()
	s.Attempts = append(s.Attempts, Attempt{ConceptID: conceptID, Correct: correct})
	// Advance index after answer
	if s.Index < len(s.Order) && s.Order[s.Index].ID == conceptID {
		s.Index++
	}
}

func (e *Engine) IsComplete(s *Session) bool {
	s.Lock()
	defer s.Unlock()
	return s.Index >= len(s.Order)
}

// HasAttempts reports whether any answer was recorded this session,
// used to avoid resetting the XP gate baseline on an all-skipped quiz.
func (e *Engine) HasAttempts(s *Session) bool {
	s.Lock()
	defer s.Unlock()
	return len(s.Attempts) > 0
}

// SkipQuestion advances past the current question without recording an
// attempt, the skip escape hatch for questions that can't be answered.
// No grade, no XP, no remedial: the concept simply goes untested this quiz.
func (e *Engine) SkipQuestion(s *Session) {
	s.Lock()
	defer s.Unlock()
	if s.Index < len(s.Order) {
		s.Index++
	}
}

// PickQuizConcepts is the anonymous and no-state fallback. Its weakness map
// holds every concept the learner has encountered; the keys are the candidate
// pool. It falls back to the start of the DAG only when the map is empty, which
// is the demo path with no learner identity.
func PickQuizConcepts(dag *concepts.DAG, weakness map[string]float64) []*concepts.Concept {
	return pickQuizConcepts(dag, weakness, nil, true)
}

// PickQuizConceptsFromCandidates selects from an identified learner's eligible
// concepts. The weakness map is the candidate pool: its keys are the concepts
// the caller has already checked for learning evidence and prerequisite
// eligibility, and its values order them. A nil result means there was no
// eligible candidate, never that the caller should fall back to the DAG.
//
// recent holds concepts from the latest quiz-sized window of quiz attempts.
// They are used only after every non-recent candidate, so a learner does not
// see the same topics on consecutive checks while alternatives exist. A short
// pool still returns every candidate; recency is a preference, not a filter
// that can empty an otherwise valid quiz.
func PickQuizConceptsFromCandidates(dag *concepts.DAG, weakness map[string]float64, recent map[string]bool) []*concepts.Concept {
	return pickQuizConcepts(dag, weakness, recent, false)
}

type quizCandidate struct {
	c *concepts.Concept
	w float64
}

func pickQuizConcepts(dag *concepts.DAG, weakness map[string]float64, recent map[string]bool, fallback bool) []*concepts.Concept {
	var list []quizCandidate
	for _, c := range dag.Order() {
		if w, ok := weakness[c.ID]; ok {
			list = append(list, quizCandidate{c, w})
		}
	}
	if len(list) == 0 {
		if !fallback {
			return nil
		}
		order := dag.Order()
		if len(order) > QuizSize {
			order = order[:QuizSize]
		}
		return order
	}
	// Weakest first, then DAG order by id. The id tie-break matters because
	// equal weakness is common (every healthy concept is 0.0), and an unstable
	// sort would make the same learner's quiz depend on map or runtime order.
	sort.Slice(list, func(i, j int) bool {
		if list[i].w != list[j].w {
			return list[i].w > list[j].w
		}
		return list[i].c.ID < list[j].c.ID
	})

	var fresh, seenRecently []quizCandidate
	for _, s := range list {
		if recent[s.c.ID] {
			seenRecently = append(seenRecently, s)
			continue
		}
		fresh = append(fresh, s)
	}

	out := appendDiverse(nil, fresh, QuizSize)
	if len(out) < QuizSize {
		out = appendDiverse(out, seenRecently, QuizSize)
	}
	return out
}

// appendDiverse adds pool to out, one concept per subdomain first, then fills
// the remaining places from the same pool. It never adds a duplicate concept
// and never exceeds limit.
func appendDiverse(out []*concepts.Concept, pool []quizCandidate, limit int) []*concepts.Concept {
	seenID := make(map[string]bool, len(out))
	seenSubdomain := make(map[string]bool, len(out))
	for _, c := range out {
		seenID[c.ID] = true
		seenSubdomain[c.Subdomain] = true
	}
	for _, s := range pool {
		if len(out) >= limit {
			return out
		}
		if seenID[s.c.ID] || seenSubdomain[s.c.Subdomain] {
			continue
		}
		out = append(out, s.c)
		seenID[s.c.ID] = true
		seenSubdomain[s.c.Subdomain] = true
	}
	for _, s := range pool {
		if len(out) >= limit {
			return out
		}
		if seenID[s.c.ID] {
			continue
		}
		out = append(out, s.c)
		seenID[s.c.ID] = true
	}
	return out
}
