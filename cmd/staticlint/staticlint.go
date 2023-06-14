package main

import (
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var checkers []*analysis.Analyzer
	for _, el := range staticcheck.Analyzers {
		// all of SA, S and ST checks present
		if strings.Contains(el.Analyzer.Name, "S") {
			checkers = append(checkers, el.Analyzer)
		}
	}
	checkers = append(checkers,
		printf.Analyzer,
		shadow.Analyzer,

		structtag.Analyzer,
		osExitFromMainAnalyzer,
	)

	multichecker.Main(checkers...)
}
