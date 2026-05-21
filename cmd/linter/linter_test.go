package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPanicChecker(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "panictest")
}

func TestForbiddenExitChecker(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), Analyzer, "exitinmain")
}
