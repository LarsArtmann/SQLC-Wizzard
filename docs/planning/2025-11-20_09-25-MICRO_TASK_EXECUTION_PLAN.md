# Micro-Task Execution Plan - Critical Code Quality Fixes

**Created:** 2025-11-20_09-25  
**Task Count:** 27 micro-tasks (max 15min each)  
**Total Time:** ~2 hours including verification

---

## 🚨 PHASE 1: CRITICAL SYSTEM FIXES (Tasks 1-8)

### Task 1.1: Analyze project_creator.go nil check issue (3min)
**File:** `internal/creators/project_creator.go:113-123`  
**Actions:**
- Read current implementation
- Identify exact location for nil check
- Plan defensive programming approach
- Determine error message format

### Task 1.2: Add nil check to project_creator.go (5min)
**Actions:**
- Add defensive nil check after line 114 (logging)
- Add wrapped fmt.Errorf with context
- Ensure check happens before config.Marshal call
- Verify error message is clear and actionable

### Task 1.3: Test project_creator.go nil fix (2min)
**Actions:**
- Run existing tests for project_creator
- Verify build still passes
- Check error handling paths work correctly

### Task 1.4: Analyze migration status underflow issue (3min)
**File:** `internal/migration/status.go:103-107`  
**Actions:**
- Read GetPendingMigrations function
- Understand uint underflow scenario
- Plan defensive check implementation

### Task 1.5: Fix uint underflow in migration status (5min)
**Actions:**
- Add defensive check for applied >= total
- Implement return 0 or inconsistency logging
- Ensure function still returns uint type
- Maintain existing API contract

### Task 1.6: Test migration status underflow fix (2min)
**Actions:**
- Run migration status tests
- Verify edge cases handled
- Check no regression in existing functionality

### Task 1.7: Analyze rule transformer semantic inconsistency (5min)
**File:** `internal/validation/rule_transformer.go`  
**Actions:**
- Compare RequireLimit logic between TransformSafetyRules and TransformTypeSafeSafetyRules
- Identify exact expression differences
- Choose consistent "violation when true" convention
- Plan normalization approach

### Task 1.8: Fix rule transformer semantics (10min)
**Actions:**
- Update RequireLimit expressions to match polarity
- Fix MaxRowsWithoutLimit consistency
- Update any helper expression builders
- Ensure both transformers generate equivalent boolean expressions
- Add comments explaining convention choice

### Task 1.9: Test rule transformer semantic fix (5min)
**Actions:**
- Run rule transformer tests
- Verify functional parity between transformers
- Check NoSelectStar, RequireWhere, RequireLimit consistency

---

## ⚡ PHASE 2: IMPORTANT USER EXPERIENCE FIXES (Tasks 9-14)

### Task 2.1: Analyze safety_policy.go rune conversion (3min)
**File:** `internal/domain/safety_policy.go:98-123`  
**Actions:**
- Locate string(rune(i)) usage on line 112
- Verify strconv import status
- Plan decimal index conversion fix

### Task 2.2: Fix rune conversion in safety_policy.go (5min)
**Actions:**
- Add `strconv` import if missing
- Replace `string(rune(i))` with `strconv.Itoa(i)`
- Update error message construction
- Test import resolution

### Task 2.3: Test safety_policy.go rune fix (2min)
**Actions:**
- Run domain safety policy tests
- Verify error messages show decimal indices
- Check no regression in validation logic

### Task 2.4: Analyze nil rules test issue (3min)
**File:** `internal/validation/rule_transformer_test.go:271-280`  
**Actions:**
- Read current "should handle nil rules gracefully" test
- Identify that it passes empty struct instead of nil
- Decide whether to test actual nil or rename for empty struct

### Task 2.5: Fix nil rules test implementation (5min)
**Actions:**
- Choose approach: actual nil test OR rename/update for empty struct
- Update test name and comments accordingly
- Modify test implementation to match intent
- Update Expect assertions for chosen behavior

### Task 2.6: Test rule_transformer_test.go nil fix (2min)
**Actions:**
- Run validation rule transformer tests
- Verify test now accurately reflects its name
- Check no regression in other edge case tests

### Task 2.7: Analyze uintToString optimization opportunity (3min)
**File:** `internal/validation/rule_transformer.go:180-183`  
**Actions:**
- Locate MaxRowsWithoutLimit usage around lines 125-126
- Identify duplicate uintToString calls with same value
- Plan local variable optimization

### Task 2.8: Optimize uintToString in rule transformer (3min)
**Actions:**
- Add local variable: `limitStr := uintToString(limit)`
- Update Rule construction to use limitStr
- Update Message construction to use limitStr
- Verify no change in generated rule logic

### Task 2.9: Test uintToString optimization (2min)
**Actions:**
- Run validation rule transformer tests
- Verify identical output to before optimization
- Check no performance regression

---

## 🛠️ PHASE 3: PROFESSIONAL CODE QUALITY (Tasks 15-27)

### Task 3.1: Analyze redundant WithMessage tests (5min)
**File:** `internal/errors/errors_test.go:146-165`  
**Actions:**
- Read both WithMessage test specs
- Identify redundancy in test coverage
- Plan consolidation into single comprehensive test
- Choose descriptive test name

### Task 3.2: Consolidate WithMessage tests (8min)
**Actions:**
- Remove redundant test spec
- Enhance remaining test to cover both behaviors
- Update test name to reflect dual coverage
- Ensure test verifies: creates details when nil AND sets Message
- Maintain test readability and clarity

### Task 3.3: Test consolidated WithMessage test (2min)
**Actions:**
- Run errors test suite
- Verify consolidated test covers both scenarios
- Check no loss of test coverage

### Task 3.4: Analyze nil error wrapping test (3min)
**File:** `internal/errors/errors_test.go:407-412`  
**Actions:**
- Read current "should handle nil error" test
- Identify missing assertions for Code and Component
- Determine expected default values for sentinel behavior

### Task 3.5: Enhance nil error wrapping test (5min)
**Actions:**
- Keep existing message assertion
- Add Expect assertion for wrapped.Code
- Add Expect assertion for wrapped.Component
- Ensure assertions match package defaults
- Add comments explaining sentinel behavior

### Task 3.6: Test enhanced nil error wrapping test (2min)
**Actions:**
- Run enhanced error wrapping test
- Verify new assertions pass
- Check sentinel behavior is properly documented

### Task 3.7: Analyze Wrapf description clobbering issue (5min)
**File:** `internal/errors/errors_test.go:435-452`  
**Actions:**
- Read Wrapf implementation and tests
- Understand baseErr.Description overwriting issue
- Plan approach: merge descriptions vs clobber

### Task 3.8: Fix Wrapf description handling (10min)
**Actions:**
- Locate Wrapf function implementation
- Modify to merge descriptions instead of overwriting
- Handle empty baseErr.Description gracefully
- Update single WithDescription call with combined string
- Ensure descriptive, readable format for merged descriptions

### Task 3.9: Test Wrapf description fix (5min)
**Actions:**
- Run Wrapf tests
- Verify descriptions are merged correctly
- Check both empty and non-empty base descriptions
- Ensure no regression in existing functionality

### Task 3.10: Analyze CombineErrors test gaps (5min)
**File:** `internal/errors/errors_test.go:477-497`  
**Actions:**
- Read current CombineErrors tests
- Identify missing component/message assertions
- Plan new test case for CombineErrors(nil, appErr)

### Task 3.11: Enhance CombineErrors existing test (5min)
**Actions:**
- Add component assertion for wrapped std error ("unknown")
- Add message/description assertion containing "standard error"
- Enhance existing "should wrap non-application errors" test
- Ensure comprehensive coverage of wrapped std error behavior

### Task 3.12: Add new CombineErrors test case (5min)
**Actions:**
- Add test case for CombineErrors(nil, appErr)
- Verify nils are skipped properly
- Assert count is 1 (not 2)
- Ensure no "Cannot wrap nil error" entry appears
- Use descriptive test name reflecting nil skipping behavior

### Task 3.13: Test CombineErrors enhancements (3min)
**Actions:**
- Run CombineErrors test suite
- Verify all new assertions pass
- Check new nil handling test works correctly
- Ensure no regression in existing CombineErrors functionality

### Task 3.14: Analyze sentinel error message assertions (5min)
**File:** `internal/errors/errors_test.go:525-545`  
**Actions:**
- Identify which sentinel errors exist and their expected Messages
- Plan Expect assertions for each error's Message field
- Determine exact expected message strings

### Task 3.15: Add Message field assertions to sentinel tests (5min)
**Actions:**
- Add Expect assertions for ErrConfigParseFailed.Message
- Add Expect assertions for ErrInvalidValue.Message  
- Add Expect assertions for ErrFileNotFound.Message
- Add Expect assertions for ErrInvalidType.Message
- Update tests accordingly to verify message stability

### Task 3.16: Test sentinel error message assertions (2min)
**Actions:**
- Run Base Errors test suite
- Verify all Message assertions pass
- Confirm message strings match current constants
- Check tests fail appropriately if messages change

---

## 🎯 EXECUTION CHECKLISTS

### Pre-Execution Checklist
- [ ] Current working directory verified
- [ ] Git status is clean (previous changes committed)
- [ ] Test environment is ready
- [ ] All target files identified and accessible
- [ ] Dependencies and imports are available

### During Execution Checklist
- [ ] Each task completed in order
- [ ] Tests run after each micro-task group
- [ ] Build passes continuously
- [ ] No regression introduced
- [ ] Time limits respected

### Post-Execution Checklist
- [ ] All 27 micro-tasks completed
- [ ] Full test suite passes
- [ ] No build errors or warnings
- [ ] Code review passes quality gates
- [ ] Documentation remains accurate
- [ ] Git commit ready with detailed message

---

## 🚨 ROLLBACK TRIGGERS

### Immediate Rollback Conditions
- Any test failure after a fix
- Build breaks or new warnings
- Unexpected runtime behavior
- Performance regression

### Rollback Procedure
1. Stop current task immediately
2. `git status` to identify changed files
3. `git stash` or `git reset --hard HEAD~1`
4. Investigate failure in isolation
5. Fix approach and re-apply carefully
6. Verify with targeted tests only

---

## 📊 SUCCESS METRICS

### Quantitative Metrics
- 0 critical bugs remaining (Tasks 1, 8, 10)
- 100% of test improvements implemented (Tasks 3-7)
- 0 semantic inconsistencies (Tasks 2, 9, 11)
- All 27 micro-tasks completed

### Qualitative Metrics
- Code prevents runtime failures
- Error messages are clear and helpful
- Tests provide comprehensive coverage
- Rule semantics are consistent and predictable
- Code is maintainable and documented

---

**Total Estimated Time:** 120 minutes (2 hours)  
**Buffer Time:** 30 minutes for investigation  
**Maximum Session Time:** 2.5 hours

---

*This micro-task plan ensures systematic, verifiable improvement while maintaining system stability throughout the process.*