package utils_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

// Helper function to create test cases for string transformations
type stringTestCase struct {
	input    string
	expected string
}

// Helper function to run multiple test cases
func runStringTests(testFunc func(string) string, testCases []stringTestCase) {
	for _, tc := range testCases {
		result := testFunc(tc.input)
		Expect(result).To(Equal(tc.expected), "for input: %s", tc.input)
	}
}

// Helper function to test edge cases
func testEdgeCases(testFunc func(string) string, emptyExpected, singleExpected string) {
	Expect(testFunc("")).To(Equal(emptyExpected))
	Expect(testFunc("A")).To(Equal(singleExpected))
}

// Helper function for truncate test cases
type truncateTestCase struct {
	input    string
	length   int
	expected string
}

// Helper function to run truncate test cases
func runTruncateTests(testFunc func(string, int) string, testCases []truncateTestCase) {
	for _, tc := range testCases {
		result := testFunc(tc.input, tc.length)
		Expect(result).To(Equal(tc.expected),
			"for input: %s, length: %d", tc.input, tc.length)
	}
}

// caseConversionTestCase represents a test case for string case conversion functions
type caseConversionTestCase struct {
	testCases     []stringTestCase
	edgeCaseFunc  func(string) string
	edgeExpected  string
	edgeLeading   string
	edgeTrailing  string
}

// runCaseConversionTests runs comprehensive tests for case conversion functions
func runCaseConversionTests(testFunc func(string) string, description string, testCase caseConversionTestCase) {
	It("should convert CamelCase to "+description, func() {
		runStringTests(testFunc, testCase.testCases)
	})

	It("should handle edge cases", func() {
		testEdgeCases(testFunc, testCase.edgeExpected, "a")

		if testCase.edgeLeading != "" {
			Expect(testFunc("-Leading")).To(Equal(testCase.edgeLeading))
		}
		if testCase.edgeTrailing != "" {
			Expect(testFunc("Trailing"+testCase.edgeTrailing[8:])).To(Equal(testCase.edgeTrailing))
		}
	})
}

var _ = Describe("StringToCamelCase", func() {
	It("should convert snake_case to CamelCase", func() {
		testCases := []stringTestCase{
			{"snake_case", "SnakeCase"},
			{"simple", "Simple"},
			{"multiple_words_here", "MultipleWordsHere"},
			{"alreadyCamelCase", "Alreadycamelcase"},
			{"", ""},
			{"single", "Single"},
		}

		runStringTests(utils.StringToCamelCase, testCases)
	})

	It("should handle edge cases", func() {
		testEdgeCases(utils.StringToCamelCase, "", "A")

		Expect(utils.StringToCamelCase("_leading_underscore")).To(Equal("LeadingUnderscore"))
		Expect(utils.StringToCamelCase("trailing_underscore_")).To(Equal("TrailingUnderscore"))
		Expect(utils.StringToCamelCase("__multiple___underscores__")).To(Equal("MultipleUnderscores"))
	})
})

var _ = Describe("StringToSnakeCase", func() {
	testCase := caseConversionTestCase{
		testCases: []stringTestCase{
			{"CamelCase", "camel_case"},
			{"Simple", "simple"},
			{"MultipleWordsHere", "multiple_words_here"},
			{"already_snake_case", "already_snake_case"},
			{"", ""},
			{"Single", "single"},
			{"XMLHttpRequest", "xmlhttp_request"},
			{"UserID", "user_id"},
		},
		edgeCaseFunc: utils.StringToSnakeCase,
		edgeExpected: "",
		edgeLeading: "-_leading",
		edgeTrailing: "trailing_",
	}

	runCaseConversionTests(utils.StringToSnakeCase, "snake_case", testCase)
})

var _ = Describe("StringToKebabCase", func() {
	testCase := caseConversionTestCase{
		testCases: []stringTestCase{
			{"CamelCase", "camel-case"},
			{"Simple", "simple"},
			{"MultipleWordsHere", "multiple-words-here"},
			{"already-kebab-case", "already-kebab-case"},
			{"", ""},
			{"Single", "single"},
			{"XMLHttpRequest", "xmlhttp-request"},
			{"UserID", "user-id"},
		},
		edgeCaseFunc: utils.StringToKebabCase,
		edgeExpected: "",
		edgeLeading: "--leading",
		edgeTrailing: "trailing-",
	}

	runCaseConversionTests(utils.StringToKebabCase, "kebab-case", testCase)
})

var _ = Describe("Pluralize and Singularize", func() {
	It("should pluralize and singularize common nouns correctly", func() {
		pluralCases := []stringTestCase{
			{"user", "users"},
			{"item", "items"},
			{"order", "orders"},
			{"query", "queries"},
			{"schema", "schemas"},
			{"index", "indices"},
			{"status", "statuses"},
			{"match", "matches"},
			{"box", "boxes"},
			{"person", "people"},
			{"child", "children"},
		}

		singularCases := []stringTestCase{
			{"users", "user"},
			{"items", "item"},
			{"orders", "order"},
			{"queries", "query"},
			{"schemas", "schema"},
			{"indices", "index"},
			{"statuses", "status"},
			{"matches", "match"},
			{"boxes", "box"},
			{"people", "person"},
			{"children", "child"},
		}

		runStringTests(utils.Pluralize, pluralCases)
		runStringTests(utils.Singularize, singularCases)
	})

	It("should handle edge cases for pluralization", func() {
		Expect(utils.Pluralize("")).To(Equal(""))
		Expect(utils.Pluralize("data")).To(Equal("data"))   // Already plural
		Expect(utils.Pluralize("sheep")).To(Equal("sheep")) // Irregular plural
	})

	It("should handle edge cases for singularization", func() {
		Expect(utils.Singularize("")).To(Equal(""))
		Expect(utils.Singularize("data")).To(Equal("data"))   // Already singular
		Expect(utils.Singularize("sheep")).To(Equal("sheep")) // Irregular plural
	})
})

var _ = Describe("IsValidIdentifier", func() {
	It("should validate Go identifiers correctly", func() {
		validCases := []string{"validIdentifier", "ValidIdentifier", "valid_identifier", "valid123", "_private", "a", "_123"}
		invalidCases := []string{"123invalid", "invalid-char", "invalid.char", "", "invalid space"}

		for _, valid := range validCases {
			Expect(utils.IsValidIdentifier(valid)).To(BeTrue(), "for valid input: %s", valid)
		}

		for _, invalid := range invalidCases {
			Expect(utils.IsValidIdentifier(invalid)).To(BeFalse(), "for invalid input: %s", invalid)
		}
	})
})

var _ = Describe("EscapeSQLIdentifier", func() {
	It("should escape SQL identifiers properly", func() {
		testCases := []stringTestCase{
			{"simple", `"simple"`},
			{"with space", `"with space"`},
			{"with-dash", `"with-dash"`},
			{"123starting", `"123starting"`},
			{"already_escaped", `"already_escaped"`},
			{"", `""`},
			{"table.name", `"table.name"`},
		}

		runStringTests(utils.EscapeSQLIdentifier, testCases)
	})
})

var _ = Describe("File Extension Functions", func() {
	It("should get file extension correctly", func() {
		testCases := []stringTestCase{
			{"file.txt", ".txt"},
			{"file.sql", ".sql"},
			{"file.go", ".go"},
			{"file", ""},
			{".hidden", ".hidden"},
			{"file.tar.gz", ".gz"},
			{"path/to/file.txt", ".txt"},
			{"", ""},
		}

		runStringTests(utils.GetFileExtension, testCases)
	})

	It("should check file extension correctly", func() {
		Expect(utils.HasExtension("file.txt", ".txt")).To(BeTrue())
		Expect(utils.HasExtension("file.txt", "txt")).To(BeTrue())
		Expect(utils.HasExtension("file.txt", ".sql")).To(BeFalse())
		Expect(utils.HasExtension("file", ".txt")).To(BeFalse())
		Expect(utils.HasExtension("", ".txt")).To(BeFalse())
	})
})

var _ = Describe("String Manipulation", func() {
	It("should truncate strings correctly", func() {
		testCases := []truncateTestCase{
			{"short", 10, "short"},
			{"exactlyten", 10, "exactlyten"},
			{"toolong", 5, "to..."},
			{"", 10, ""},
			{"a", 1, "a"},
		}

		runTruncateTests(utils.TruncateString, testCases)
	})

	It("should pad strings correctly", func() {
		Expect(utils.PadString("short", 10, " ")).To(Equal("short     "))
		Expect(utils.PadString("exactly", 7, " ")).To(Equal("exactly"))
		Expect(utils.PadString("toolong", 5, " ")).To(Equal("toolong")) // No truncation
		Expect(utils.PadString("", 5, "x")).To(Equal("xxxxx"))
		Expect(utils.PadString("a", 3, "0")).To(Equal("a00"))
	})
})

var _ = Describe("Error Handling", func() {
	It("should handle nil or empty inputs gracefully", func() {
		testEdgeCases(utils.StringToCamelCase, "", "A")
		testEdgeCases(utils.StringToSnakeCase, "", "a")
		testEdgeCases(utils.StringToKebabCase, "", "a")

		Expect(utils.IsValidIdentifier("")).To(BeFalse())
		Expect(utils.EscapeSQLIdentifier("")).To(Equal(`""`))
	})
})

var _ = Describe("Performance", func() {
	It("should handle large strings efficiently", func() {
		largeSnake := "this_is_a_very_long_snake_case_string_with_many_words"
		largeCamel := "thisIsAVeryLongCamelCaseStringWithManyWords"

		// These should complete quickly and not panic
		Expect(utils.StringToCamelCase(largeSnake)).NotTo(BeEmpty())
		Expect(utils.StringToSnakeCase(largeCamel)).NotTo(BeEmpty())
		Expect(utils.StringToKebabCase(largeCamel)).NotTo(BeEmpty())
	})
})
