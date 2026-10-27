package machinelearning

import (
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

func TestBackpropGen(t *testing.T) { fuzzGen(t, &backpropGen{}) }

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

func TestCrossEntropyNegativeLabelShowsMinus(t *testing.T) {
	g := &crossEntropyGen{}
	found := false
	for i := 0; i < 2000; i++ {
		p := g.Generate(generator.GeneratorContext{Difficulty: rand.Float64()})
		if !strings.Contains(p.Question, "true label 0") {
			continue
		}
		found = true
		if strings.Contains(p.Question, "-ln(1 p)") {
			t.Fatalf("negative-label question still omits the minus: %q", p.Question)
		}
		if !strings.Contains(p.Question, "-ln(1 - p)") {
			t.Fatalf("negative-label question does not show -ln(1 - p): %q", p.Question)
		}
	}
	if !found {
		t.Fatal("did not sample a true-label-0 case")
	}
}
