// sitedump draws a fixed set of practice problems for each concept named on
// stdin (one id per line) and prints them as JSON keyed by concept id. The
// static reader in site/ uses it, since GitHub Pages cannot run the engine.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/all"
)

type item struct {
	KP          int    `json:"kp"`
	Question    string `json:"q"`
	Answer      string `json:"a"`
	Explanation string `json:"e"`
}

func main() {
	per := flag.Int("per", 3, "problems per knowledge point")
	flag.Parse()
	reg := all.Registry()
	out := map[string][]item{}
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		id := strings.TrimSpace(sc.Text())
		if id == "" || !reg.Has(id) {
			continue
		}
		for kp := 1; kp <= 3; kp++ {
			for i := 0; i < *per; i++ {
				ctx := generator.GeneratorContext{Difficulty: 0.4 + 0.25*float64(i), Seed: int64(kp*100 + i), KPOrdinal: kp}
				p, err := reg.GenerateContext(id, ctx)
				if err != nil {
					fmt.Fprintln(os.Stderr, id, err)
					continue
				}
				out[id] = append(out[id], item{kp, p.Question, p.Answer, p.Explanation})
			}
		}
	}
	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
