package testcase

import (
	"context"
	"fmt"
	"net"

	"github.com/mitchellh/mapstructure"
)

type DnsTest struct {
	Name     string
	Host     string
	resolver net.Resolver
}

func init() {
	RegisterTestParser("dns", parseDnsTestCase)
}

func parseDnsTestCase(name string, values map[string]interface{}) (TestCase, error) {
	var testCase DnsTest

	testCase.Name = name

	err := mapstructure.Decode(values, &testCase)
	if err != nil {
		return nil, err
	}

	return &testCase, nil
}

func (m *DnsTest) GetName() string {
	return m.Name
}

func (m *DnsTest) GetType() string {
	return "dns"
}

func (m *DnsTest) String() string {
	return fmt.Sprintf("dns:%s:host=%s", m.Name, m.Host)
}

func (m *DnsTest) Run(ctx context.Context) TestResult {
	addrs, err := m.resolver.LookupNetIP(ctx, "ip", m.Host)
	if err != nil {
		return TestResult{
			Success:    false,
			FailReason: "LOOKUP_FAILED",
			FailDesc:   fmt.Sprintf("Lookup failed: %s", err),
		}
	}

	for _, addr := range addrs {
		if addr.IsPrivate() || addr.IsLoopback() || addr.IsUnspecified() {
			return TestResult{
				Success:    false,
				FailReason: "BOGON",
				FailDesc:   fmt.Sprintf("Got bogon address: %s", addr),
			}
		}
	}

	return TestResult{
		Success:    true,
		FailReason: "",
		FailDesc:   "",
	}
}

var _ TestCase = (*DnsTest)(nil)
