package testcase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

type HttpTest struct {
	Name               string
	Url                string
	ExpectedStatusCode int
	client             http.Client
}

func init() {
	RegisterTestParser("http", parseHttpTestCase)
}

func parseHttpTestCase(name string, values map[string]interface{}) (TestCase, error) {
	var testCase HttpTest

	testCase.Name = name

	// default value
	testCase.ExpectedStatusCode = 200

	err := mapstructure.Decode(values, &testCase)
	if err != nil {
		return nil, err
	}

	return &testCase, nil
}

func (m *HttpTest) GetName() string {
	return m.Name
}

func (m *HttpTest) GetType() string {
	return "http"
}

func (m *HttpTest) String() string {
	return fmt.Sprintf("http:%s:url=%s", m.Name, m.Url)
}

func (m *HttpTest) Run(ctx context.Context) TestResult {
	req, err := http.NewRequest("GET", m.Url, nil)
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

	if res.StatusCode != m.ExpectedStatusCode {
		return TestResult{
			Success:    false,
			FailReason: "HTTP_REQ_FAILED",
			FailDesc: fmt.Sprintf("Bad status code: expected %d, got %d",
				m.ExpectedStatusCode, res.StatusCode),
		}
	}

	return TestResult{
		Success:    true,
		FailReason: "",
		FailDesc:   "",
	}
}

var _ TestCase = (*HttpTest)(nil)
