package testing

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TestStringRepresentationSuite runs generic tests for String() method of any enum type.
// It encapsulates the common Context/It pattern for string representation tests.
func TestStringRepresentationSuite(testCases func() []EnumTestCase) {
	Context("String", func() {
		It("should return correct string representation", func() {
			RunStringRepresentationTest(testCases())
		})
	})
}

// ExpectToBeTrue asserts that the given value is true with a descriptive message.
// This helper reduces duplication of "Expect(x).To(BeTrue())" patterns in tests.
func ExpectToBeTrue(actual any, msgAndArgs ...any) {
	Expect(actual).To(BeTrue(), msgAndArgs...)
}

// ExpectToBeFalse asserts that the given value is false with a descriptive message.
// This helper reduces duplication of "Expect(x).To(BeFalse())" patterns in tests.
func ExpectToBeFalse(actual any, msgAndArgs ...any) {
	Expect(actual).To(BeFalse(), msgAndArgs...)
}

// ExpectToBeNil asserts that the given value is nil with a descriptive message.
func ExpectToBeNil(actual any, msgAndArgs ...any) {
	Expect(actual).To(BeNil(), msgAndArgs...)
}

// ExpectNotToBeNil asserts that the given value is not nil with a descriptive message.
// This helper reduces duplication of "Expect(x).NotTo(BeNil())" patterns in tests.
func ExpectNotToBeNil(actual any, msgAndArgs ...any) {
	Expect(actual).NotTo(BeNil(), msgAndArgs...)
}

// ExpectToEqual asserts that the actual value equals the expected value with a descriptive message.
// This helper reduces duplication of "Expect(x).To(Equal(y))" patterns in tests.
func ExpectToEqual(actual, expected any, msgAndArgs ...any) {
	Expect(actual).To(Equal(expected), msgAndArgs...)
}

// ExpectToNotEqual asserts that the actual value does not equal the unexpected value.
// This helper reduces duplication of "Expect(x).NotTo(Equal(y))" patterns in tests.
func ExpectToNotEqual(actual, expected any, msgAndArgs ...any) {
	Expect(actual).NotTo(Equal(expected), msgAndArgs...)
}

// ExpectToContain asserts that the actual slice or string contains the expected element or substring.
// This helper reduces duplication of "Expect(x).To(ContainElement(y))" patterns in tests.
func ExpectToContain(actual, expected any, msgAndArgs ...any) {
	Expect(actual).To(ContainElement(expected), msgAndArgs...)
}

// ExpectLength asserts that the actual value has the expected length.
// This helper reduces duplication of "Expect(x).To(HaveLen(n))" patterns in tests.
func ExpectLength(actual any, length int, msgAndArgs ...any) {
	Expect(actual).To(HaveLen(length), msgAndArgs...)
}

// CreateMicroserviceTemplateData creates a TemplateData for microservice projects.
// This helper eliminates duplicate template data creation patterns across tests.
func CreateMicroserviceTemplateData(
	projectName, packageName, packagePath string,
) generated.TemplateData {
	return generated.TemplateData{
		ProjectName: projectName,
		ProjectType: generated.ProjectTypeMicroservice,
		Package: generated.PackageConfig{
			Name: packageName,
			Path: packagePath,
		},
		Database: generated.DatabaseConfig{
			Engine:    generated.DatabaseTypePostgreSQL,
			UseUUIDs:  true,
			UseJSON:   true,
			UseArrays: true,
		},
	}
}

// CreateMinimalTemplateData creates a minimal TemplateData for hobby/simple projects.
func CreateMinimalTemplateData(projectName string) generated.TemplateData {
	return generated.TemplateData{
		ProjectName: projectName,
		ProjectType: generated.ProjectTypeHobby,
		Package:     generated.PackageConfig{Name: "db", Path: "github.com/user/minimal"},
		Database:    generated.DatabaseConfig{Engine: generated.DatabaseTypeSQLite},
	}
}

// CreateFullEmitOptions creates EmitOptions with all fields set for comprehensive testing.
func CreateFullEmitOptions() generated.EmitOptions {
	return generated.EmitOptions{
		EmitJSONTags: true, EmitPreparedQueries: true, EmitInterface: true, EmitEmptySlices: true,
		EmitResultStructPointers: false, EmitParamsStructPointers: false, EmitEnumValidMethod: true,
		EmitAllEnumValues: true, JSONTagsCaseStyle: "camel",
	}
}
