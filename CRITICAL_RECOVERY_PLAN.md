# SQLC-Wizard Critical Test Coverage Recovery Plan

**Date:** 2025-12-12_16-52
**Issue:** Wizard test coverage stuck at 2.9% despite 1,167+ lines of tests
**Root Cause:** Our tests only access template data, not actual wizard method execution

---

## 🚨 1. WHAT I FORGOT & COULD HAVE DONE BETTER

### **Critical Mistakes Made:**

1. **❌ Testing Approach Error** - I tested data setup instead of method execution
2. **❌ Assumption Error** - Assumed template data tests would cover wizard.go methods
3. **❌ Strategy Error** - Didn't mock/execute actual wizard methods
4. **❌ Git Workflow Error** - Created extensive work without immediate commits
5. **❌ Verification Error** - Didn't verify coverage improvements after each test file

### **What I Should Have Done Better:**

1. **✅ Direct Method Testing** - Test `wiz.Run()`, `w.generateConfig()`, `w.showSummary()` directly
2. **✅ Mock UI Components** - Mock `w.ui` calls to enable method execution
3. **✅ Incremental Verification** - Check coverage after each test batch
4. **✅ Commit-Immediately Strategy** - Commit each test file individually
5. **✅ Use Mocking Libraries** - Use established mocking frameworks

### **What Could Still Be Improved:**

1. **🔄 Test Architecture** - Need proper mocking framework setup
2. **🔄 Type Models** - Could improve generated types for better testability
3. **🔄 Library Usage** - Could use testify/gomock for mocking
4. **🔄 Coverage Analysis** - Need per-line coverage analysis to identify gaps

---

## 🎯 2. COMPREHENSIVE MULTI-STEP EXECUTION PLAN

### **PHASE 1: CRITICAL COVERAGE RECOVERY (Immediate - 30 mins)**

#### **Step 1: Add Mocking Framework** (5 mins)

- Add testify/mock and gomock dependencies
- Create mocks for UIHelper and step interfaces
- Commit: "feat(wizard): Add mocking framework for method testing"

#### **Step 2: Create Wizard Method Execution Tests** (10 mins)

- Test `w.Run()` with mocked UI components
- Test `w.generateConfig()` with mocked templates
- Test `w.showSummary()` with mocked UI
- Commit: "test(wizard): Add direct method execution tests with mocks"

#### **Step 3: Verify Coverage Improvement** (5 mins)

- Run tests with coverage
- Verify wizard.go methods now have >50% coverage
- Commit: "fix(wizard): Verify method coverage improvement"

#### **Step 4: Add Error Scenario Tests** (10 mins)

- Test wizard methods with error conditions
- Test failure paths in Run(), generateConfig(), showSummary()
- Commit: "test(wizard): Add error scenario tests for wizard methods"

### **PHASE 2: ARCHITECTURE IMPROVEMENT (15 mins)**

#### **Step 5: Improve Type Models** (5 mins)

- Add test helper methods to generated types
- Improve type safety for testing scenarios
- Commit: "feat(generated): Add test helpers for improved testability"

#### **Step 6: Leverage Established Libraries** (10 mins)

- Add testify/assert for better assertions
- Add testcontainers for integration testing
- Add gomega matchers for complex validations
- Commit: "deps(testing): Add established testing libraries"

### **PHASE 3: INTEGRATION VALIDATION (10 mins)**

#### **Step 7: Add Integration Tests** (5 mins)

- Test complete wizard workflows with mock UI
- Test real sqlc generation with wizard output
- Commit: "test(wizard): Add integration tests for complete workflows"

#### **Step 8: Performance and Security Tests** (5 mins)

- Test wizard performance with large configurations
- Test wizard security with malformed inputs
- Commit: "test(wizard): Add performance and security tests"

---

## 📊 3. WORK REQUIRED vs IMPACT MATRIX

| Task                        | Work Required | Impact   | Priority |
| --------------------------- | ------------- | -------- | -------- |
| **Direct Method Testing**   | Low           | Critical | P1       |
| **Mocking Framework Setup** | Low           | Critical | P1       |
| **Error Scenario Tests**    | Medium        | High     | P2       |
| **Integration Tests**       | Medium        | High     | P2       |
| **Type Model Improvements** | Low           | Medium   | P3       |
| **Established Libraries**   | Low           | Medium   | P3       |
| **Performance Tests**       | Medium        | Low      | P4       |
| **Security Tests**          | Medium        | Low      | P4       |

---

## 🔍 4. EXISTING CODE ANALYSIS FOR REUSE

### **What We Already Have:**

1. **✅ Complete Wizard Structure** - All wizard methods implemented and working
2. **✅ Step Interface Pattern** - Perfect for mocking individual steps
3. **✅ UIHelper Abstraction** - Clean interface for UI mocking
4. **✅ Template System** - Ready for mocking template generation
5. **✅ Generated Types** - Type-safe data structures for testing

### **Reusable Components:**

1. **`Wizard` struct** - Can test with mocked components
2. **`UIHelper` interface** - Perfect for mocking
3. **Step interfaces** - Can mock individual wizard steps
4. **`Template` interface** - Can mock template generation
5. **`GeneratedTemplateData`** - Type-safe test data setup

### **What We DON'T Need to Implement:**

1. ❌ New wizard infrastructure (already exists)
2. ❌ New type definitions (already comprehensive)
3. ❌ New step implementations (already complete)
4. ❌ New UI components (already abstracted)

---

## 🏗️ 5. TYPE MODEL IMPROVEMENT PLAN

### **Current Issues:**

1. Wizard methods hard to test due to UI dependencies
2. No test helpers for common test scenarios
3. Generated types lack test utility methods

### **Proposed Improvements:**

1. **Add Test Helper Methods to Generated Types:**

   ```go
   // In generated/types.go
   func (td TemplateData) WithDefaults() TemplateData { ... }
   func (td TemplateData) IsValidForTesting() error { ... }
   ```

2. **Add Wizard Testing Helpers:**

   ```go
   // New file: internal/wizard/wizard_test_helpers.go
   func NewMockedWizard(ui MockUI, steps MockSteps) *Wizard { ... }
   func TestWizardExecution(wiz *Wizard, data TemplateData) error { ... }
   ```

3. **Improve Interface Design:**
   ```go
   // Make wizard more testable by extracting interfaces
   type WizardRunner interface {
       Run(data *TemplateData) (*WizardResult, error)
       GenerateConfig(data *TemplateData) error
       ShowSummary(data *TemplateData) error
   }
   ```

---

## 📚 6. ESTABLISHED LIBRARIES PLAN

### **Libraries to Add:**

1. **`github.com/stretchr/testify`** - Better assertions and mocking
2. **`github.com/golang/mock`** - Interface mocking
3. **`github.com/testcontainers/testcontainers-go`** - Integration testing
4. **`github.com/onsi/ginkgo/ginkgo/extensions/table`** - Table-driven tests
5. **`github.com/onsi/gomega/types`** - Custom matchers

### **Benefits:**

1. **Better Testing** - Professional-grade assertions and mocks
2. **Less Boilerplate** - Reusable testing patterns
3. **Integration Testing** - Real database testing
4. **Table-Driven Tests** - Cleaner test organization
5. **Custom Matchers** - Better validation logic

### **Implementation Plan:**

1. Add dependencies to go.mod
2. Update import statements in test files
3. Refactor existing tests to use new libraries
4. Add new tests with improved patterns

---

## 🚀 EXECUTION PHASE

### **IMMEDIATE NEXT STEP: Add testify for better testing**

Let me add the testing library and create proper method execution tests:
