package testcase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type HttpTest struct {
	name               string
	url                string
	expectedStatusCode int
	client             http.Client
}

func init() {
	RegisterTestParser("http", parseHttpTestCase)
}

func parseHttpTestCase(name string, values map[string]interface{}) (TestCase, error) {
	var testCase HttpTest

	testCase.name = name

	url, ok := values["url"]
	if !ok {
		return nil, errors.New("url key not present in HTTP test case.")
	}

	testCase.url, ok = url.(string)
	if !ok {
		return nil, errors.New("url value for HTTP test case is not a string")
	}

	expectedStatusCode, ok := values["expected_status_code"]
	if ok {
		testCase.expectedStatusCode, ok = expectedStatusCode.(int)
		if !ok {
			return nil, errors.New("expected_status_code is not an integer")
		}
	} else {
		testCase.expectedStatusCode = 200
	}

	testCase.client = *http.DefaultClient

	return &testCase, nil
}

func (m *HttpTest) GetName() string {
	return m.name
}

func (m *HttpTest) GetType() string {
	return "http"
}

func (m *HttpTest) String() string {
	return fmt.Sprintf("http:%s:url=%s", m.name, m.url)
}

func (m *HttpTest) Run(ctx context.Context) TestResult {
	req, err := http.NewRequest("GET", m.url, nil)
	if err != nil {
		return TestResult{
			Success:    false,
			FailReason: "REQ_CREATE_FAILED",
			FailDesc:   fmt.Sprintf("Failed creating HTTP request: %s", err),
		}
	}

	req = req.WithContext(ctx)

	res, err := m.client.Do(req)
	if err != nil {
		return TestResult{
			Success:    false,
			FailReason: "HTTP_REQ_FAILED",
			FailDesc:   fmt.Sprintf("HTTP request failed: %s", err),
		}
	}

	if res.StatusCode != m.expectedStatusCode {
		return TestResult{
			Success:    false,
			FailReason: "HTTP_REQ_FAILED",
			FailDesc: fmt.Sprintf("Bad status code: expected %d, got %d",
				m.expectedStatusCode, res.StatusCode),
		}
	}

	return TestResult{
		Success:    true,
		FailReason: "",
		FailDesc:   "",
	}
}

var _ TestCase = (*HttpTest)(nil)
