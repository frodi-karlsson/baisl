package baisl_test

import (
	"strings"
	"testing"

	"github.com/frodi-karlsson/baisl"
)

type parserTest struct {
	path     string
	expected string
}

type failParserTest struct {
	path          string
	errorContains string
}

var parserTests = []parserTest{
	{"raw/ret2.baisl", "Function main(): void:\n  Block:\n    Return\n\nFunction return2(): int:\n  Block:\n    Return 2\n\n"}, // annoying extra newline i haven't dealt with
	{"raw/retParam.baisl", "Function main(): void:\n  Block:\n    Return\n\nFunction returnParam(a: int): int:\n  Block:\n    Return a\n\n"},
}

var failParserTests = []failParserTest{}

func TestParse(t *testing.T) {
	for _, test := range parserTests {
		sourceFile, err := baisl.GetSourceFile(test.path)

		if err != nil {
			t.Errorf("Error reading file: %s", err)
		}

		parser := baisl.Parser{
			SourceFile: &sourceFile,
		}

		result, err := parser.Parse()
		if err != nil {
			t.Errorf("Error parsing file: %s", err)
		}

		joined := ""
		for _, line := range result {
			joined += line.String(0) + "\n"
		}

		if joined != test.expected {
			t.Errorf("\nExpected <%s>, got <%s>", test.expected, joined)
		}
	}

	for _, test := range failParserTests {
		sourceFile, err := baisl.GetSourceFile(test.path)
		if err != nil {
			t.Errorf("Error reading file: %s", err)
		}
		parser := baisl.Parser{
			SourceFile: &sourceFile,
		}
		_, err = parser.Parse()
		if err == nil {
			t.Errorf("Expected error, got none")
		}

		if !strings.Contains(err.Error(), test.errorContains) {
			t.Errorf("Expected error containing <%s>, got <%s>", test.errorContains, err.Error())
		}
	}
}
