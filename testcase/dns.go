package testcase

import (
	"context"
	"errors"
	"fmt"
	"net"
)

type DnsTest struct {
	name     string
	host     string
	resolver net.Resolver
}

func init() {
	RegisterTestParser("dns", parseDnsTestCase)
}

func parseDnsTestCase(name string, values map[string]interface{}) (TestCase, error) {
	var testCase DnsTest

	testCase.name = name

	host, ok := values["host"]
	if !ok {
		return nil, errors.New("host key not present in DNS test case.")
	}

	testCase.host, ok = host.(string)
	if !ok {
		return nil, errors.New("host value for DNS test case is not a string")
	}

	testCase.resolver = net.Resolver{}

	return &testCase, nil
}

func (m *DnsTest) GetName() string {
	return m.name
}

func (m *DnsTest) GetType() string {
	return "dns"
}

func (m *DnsTest) String() string {
	return fmt.Sprintf("dns:%s:host=%s", m.name, m.host)
}

func (m *DnsTest) Run(ctx context.Context) TestResult {
	addrs, err := m.resolver.LookupNetIP(ctx, "ip", m.host)
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
