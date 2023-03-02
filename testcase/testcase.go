package testcase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type TestCase interface {
	GetName() string
	GetType() string
	String() string
	Run(context.Context) TestResult
}

type TestResult struct {
	Success    bool
	FailReason string
	FailDesc   string
}

func ReadTestSuiteFile(filename string) ([]TestCase, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Error reading test suite: %s", err)
	}

	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("Error parsing test suite json: %s", err)
	}

	jsonDict, ok := jsonData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Test suite file not a JSON object.")
	}

	testCases, ok := jsonDict["tests"]
	if !ok {
		return nil, fmt.Errorf("Test suite object does not have a 'tests' key.")
	}

	testCasesList, ok := testCases.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Test suite object 'tests' field does not contain a list.")
	}

	var ret []TestCase
	for i, testCase := range testCasesList {
		testCaseDict, ok := testCase.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Test case %d not an object.", i)
		}

		parsedTestCase, err := parseTestCase(testCaseDict)
		if err != nil {
			return nil, fmt.Errorf("Error parsing test case %d: %s", i, err)
		}

		ret = append(ret, parsedTestCase)
	}

	return ret, nil
}

func parseTestCase(jsonDict map[string]interface{}) (TestCase, error) {
	var testName string
	var testType string
	var ok bool
	for k, v := range jsonDict {
		switch k {
		case "name":
			testName, ok = v.(string)
			if !ok {
				panic("Test name not a string")
			}
			if testName == "" {
				panic("Test name is empty")
			}
		case "type":
			testType, ok = v.(string)
			if !ok {
				panic("Test type not a string")
			}
		}
	}

	if testName == "" {
		return nil, fmt.Errorf("Test case has no name")
	}

	if testType == "" {
		return nil, fmt.Errorf("Test case has no type")
	}

	parser, ok := registry[testType]
	if !ok {
		return nil, fmt.Errorf("Could not find a parser for test type: %s", testType)
	}

	testCase, err := parser(testName, jsonDict)
	if err != nil {
		return nil, err
	}

	return testCase, nil
}
