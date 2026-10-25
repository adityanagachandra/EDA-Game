package numtheory

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/chuma-beep/mathua/internal/generator"
)

func fuzzGen(t *testing.T, gen generator.Generator) {
	t.Helper()
	for i := 0; i < 100; i++ {
		d := rand.Float64()
		p := gen.Generate(generator.GeneratorContext{Difficulty: d})
		if p.Question == "" || p.Answer == "" || p.Explanation == "" {
			t.Errorf("empty field at difficulty=%.2f: q=%q a=%q e=%q", d, p.Question, p.Answer, p.Explanation)
		}
	}
}

func TestDivisibilityGen(t *testing.T) { fuzzGen(t, &divisibilityGen{}) }
func TestGcdEuclideanGen(t *testing.T) { fuzzGen(t, &gcdEuclideanGen{}) }
func TestModularGen(t *testing.T)      { fuzzGen(t, &modularGen{}) }
func TestCongruenceGen(t *testing.T)   { fuzzGen(t, &congruenceGen{}) }
func TestFermatLittleGen(t *testing.T) { fuzzGen(t, &fermatLittleGen{}) }
func TestEulerPhiGen(t *testing.T)     { fuzzGen(t, &eulerPhiGen{}) }
func TestDiophantineGen(t *testing.T)  { fuzzGen(t, &diophantineGen{}) }
func TestCryptoGen(t *testing.T)       { fuzzGen(t, &cryptoGen{}) }

func TestFuzz(t *testing.T) {
	reg := generator.NewRegistry()
	Register(reg)
	for _, id := range reg.Concepts() {
		gen, _ := reg.Get(id)
		t.Run(id, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				d := rand.Float64()
				p := gen.Generate(generator.GeneratorContext{Difficulty: d})
				if p.Question == "" || p.Answer == "" || p.Explanation == "" {
					t.Fatalf("empty field at difficulty=%.2f for %s: q=%q a=%q", d, id, p.Question, p.Answer)
				}
			}
		})
	}
}

func TestDiophantineChoicesHaveOneSolution(t *testing.T) {
	g := &diophantineGen{}
	for i := 0; i < 100; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		if !strings.Contains(p.Question, "\n\nA)") {
			t.Fatalf("diophantine question has no option list: %q", p.Question)
		}
		options := strings.Split(p.Question, "\n\n")[1:]
		if len(options) != 4 {
			t.Fatalf("diophantine question has %d options, want 4: %q", len(options), p.Question)
		}
		var correctX, correctY int
		if _, err := fmt.Sscanf(p.Answer, "(%d, %d)", &correctX, &correctY); err != nil {
			t.Fatalf("cannot parse correct pair %q: %v", p.Answer, err)
		}
		c := 3*correctX + 5*correctY
		solutions := 0
		for _, option := range options {
			text := strings.TrimSpace(option[3:])
			var x, y int
			if _, err := fmt.Sscanf(text, "(%d, %d)", &x, &y); err != nil {
				t.Fatalf("cannot parse option %q: %v", option, err)
			}
			if 3*x+5*y == c {
				solutions++
			}
		}
		if solutions != 1 {
			t.Fatalf("diophantine question has %d correct options, want 1: %q", solutions, p.Question)
		}
		if !strings.Contains(p.Question, p.Answer) {
			t.Fatalf("correct answer %q is not among the options: %q", p.Answer, p.Question)
		}
	}
}
