package main

import (
	"context"
	"firetest/testcase"
	"fmt"
	"os"
	"time"
)

type TestResultPair struct {
	test   testcase.TestCase
	result testcase.TestResult
}

type TestCaseSettings struct {
	timeout time.Duration
}

func runTest(
	ctx context.Context,
	test testcase.TestCase,
	settings TestCaseSettings,
	resultChan chan TestResultPair,
) {
	ctx, cancelFunc := context.WithTimeout(ctx, settings.timeout)
	defer cancelFunc()

	result := test.Run(ctx)
	resultChan <- TestResultPair{test: test, result: result}
}

func main() {
	fmt.Println("Firetest v0.1.0")

	testFile := "tests.json"
	if len(os.Args) > 1 {
		testFile = os.Args[1]
	}

	tests, err := testcase.ReadTestSuiteFile(testFile)
	if err != nil {
		fmt.Printf("Error reading test suite: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Running %d test(s).\n", len(tests))

	resultChan := make(chan TestResultPair)

	for i, test := range tests {
		settings := TestCaseSettings{
			timeout: 10 * time.Second,
		}
		fmt.Printf("Running Test %d: %s\n", i, test.String())
		go runTest(context.Background(), test, settings, resultChan)
	}

	for range tests {
		resultPair := <-resultChan
		fmt.Printf("Finished Test %s: ", resultPair.test.String())
		if resultPair.result.Success {
			fmt.Println("Success")
		} else {
			fmt.Printf("Failed (%s): %s\n",
				resultPair.result.FailReason,
				resultPair.result.FailDesc)
		}
	}

	fmt.Println("Done.")
}
