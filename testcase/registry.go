package testcase

type TestCaseParser func(name string, values map[string]interface{}) (TestCase, error)

var registry = make(map[string]TestCaseParser)

func RegisterTestParser(name string, parser TestCaseParser) {
	registry[name] = parser
}
