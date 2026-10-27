package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fixturesPath := flag.String("fixtures", "tools/fsrs-replay/fixtures.json", "synthetic fixture file")
	outPath := flag.String("out", "", "write the report here instead of standard output")
	flag.Parse()

	fixtures, err := loadFixtures(*fixturesPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	report, err := Run(fixtures)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := writeReport(*outPath, report); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
