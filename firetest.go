package main

import (
	"context"
	"firetest/testcase"
	"fmt"
	"os"
)

type TestResultPair struct {
	test   testcase.TestCase
	result testcase.TestResult
}

func runTest(test testcase.TestCase, resultChan chan TestResultPair) {
	result := test.Run(context.Background())
	resultChan <- TestResultPair{test: test, result: result}
}

func main() {
	fmt.Println("Firetest v0.1.0")

	tests, err := testcase.ReadTestSuiteFile("tests.json")
	if err != nil {
		fmt.Printf("Error reading test suite: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running %d test(s).\n", len(tests))

	resultChan := make(chan TestResultPair)

	for i, test := range tests {
		fmt.Printf("Running Test %d: %s\n", i, test.String())
		go runTest(test, resultChan)
	}

	for range tests {
		resultPair := <-resultChan
		fmt.Printf("Finished Test %s: ", resultPair.test.String())
		if resultPair.result.Success {
			fmt.Println("Success.")
		} else {
			fmt.Printf("Failed (%s): %s\n",
				resultPair.result.FailReason,
				resultPair.result.FailDesc)
		}
	}

	fmt.Println("Done.")
}
