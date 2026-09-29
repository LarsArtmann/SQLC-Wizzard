# SQLC-Wizard Detailed Task Breakdown (150 Tasks, 15min each)

**Date:** 2026-01-13_17-06
**Total Tasks:** 150
**Task Size:** Max 15min each
**Total Estimated Time:** ~37.5 hours
**Parallelizable:** ~10 hours (docs, examples can run concurrently)

---

## 🎯 Task Sorting Strategy

Tasks sorted by **(Impact × Importance) / Effort** ratio:

1. 🔴 CRITICAL - High Impact, Low Effort (Do First)
2. 🟡 IMPORTANT - Medium Impact, Medium Effort (Do Second)
3. 🟢 SHOULD HAVE - Lower Impact, Higher Effort (Do Third)

---

## 📋 Complete Task List (150 Tasks)

### ~~Phase 0: Quick Wins (1% → 51%) - 23 Tasks~~ EXECUTED — user guide (`6507a01`), test fixes, hobby example (`9db82a7`); wizard coverage at 34.2% (→ TODO_LIST T6)

#### QW-01: Fix 3 Failing Integration Tests (30min = 2 tasks)

| ID      | Task                                    | Priority | Time  | Impact |
| ------- | --------------------------------------- | -------- | ----- | ------ |
| QW-01-A | Analyze 3 failing test failures         | 🔴       | 10min | 🔴🔴🔴 |
| QW-01-B | Fix validation failure test expectation | 🔴       | 15min | 🔴🔴🔴 |
| QW-01-C | Fix UI panic in mock test               | 🔴       | 10min | 🔴🔴🔴 |
| QW-01-D | Fix data flow test expectation          | 🔴       | 10min | 🔴🔴🔴 |
| QW-01-E | Run full test suite to verify fixes     | 🔴       | 10min | 🔴🔴🔴 |
| QW-01-F | Commit integration test fixes           | 🔴       | 5min  | 🔴🔴🔴 |

**Total:** 6 tasks, 60min

---

#### QW-02: Create "Getting Started" User Guide (2h = 8 tasks)

| ID      | Task                                       | Priority | Time  | Impact |
| ------- | ------------------------------------------ | -------- | ----- | ------ |
| QW-02-A | Create user-guide/ directory structure     | 🔴       | 5min  | 🔴🔴   |
| QW-02-B | Write installation section (all methods)   | 🔴       | 15min | 🔴🔴🔴 |
| QW-02-C | Write quick start tutorial (hobby project) | 🔴       | 20min | 🔴🔴🔴 |
| QW-02-D | Add common project types section           | 🔴       | 15min | 🔴🔴🔴 |
| QW-02-E | Add configuration options reference        | 🔴       | 15min | 🔴🔴🔴 |
| QW-02-F | Add troubleshooting basics section         | 🔴       | 15min | 🔴🔴   |
| QW-02-G | Review and format documentation            | 🔴       | 10min | 🔴🔴   |
| QW-02-H | Add images/screenshots to guide            | 🔴       | 15min | 🔴🔴   |

**Total:** 8 tasks, 120min

---

#### QW-03: Add Wizard Step Tests (Critical Paths) (2h = 8 tasks)

| ID      | Task                                 | Priority | Time  | Impact |
| ------- | ------------------------------------ | -------- | ----- | ------ |
| QW-03-A | Test project type step validation    | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-B | Test database step validation        | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-C | Test project details step validation | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-D | Test features step validation        | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-E | Test output step validation          | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-F | Test wizard orchestration flow       | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-G | Test wizard error handling           | 🔴       | 15min | 🔴🔴🔴 |
| QW-03-H | Run wizard tests and verify coverage | 🔴       | 10min | 🔴🔴   |

**Total:** 8 tasks, 120min

---

#### QW-04: Create Basic Example (Hobby/SQLite) (1h = 4 tasks)

| ID      | Task                                    | Priority | Time  | Impact |
| ------- | --------------------------------------- | -------- | ----- | ------ |
| QW-04-A | Create examples/hobby-sqlite/ directory | 🔴       | 5min  | 🔴🔴   |
| QW-04-B | Generate hobby project with wizard      | 🔴       | 15min | 🔴🔴🔴 |
| QW-04-C | Add README with project description     | 🔴       | 15min | 🔴🔴🔴 |
| QW-04-D | Test example builds and runs            | 🔴       | 10min | 🔴🔴   |

**Total:** 4 tasks, 60min

**Phase 0 Total:** 26 tasks, ~6 hours

---

### ~~Phase 1: Critical Foundation (4% → 64%) - 45 Tasks~~ PARTIALLY EXECUTED — guides done; coverage targets open (→ TODO_LIST T6/T7, T13)

#### CF-01: Complete Wizard Test Coverage to 60% (2h = 8 tasks)

| ID      | Task                                | Priority | Time  | Impact |
| ------- | ----------------------------------- | -------- | ----- | ------ |
| CF-01-A | Test UI helper interactions         | 🔴       | 15min | 🔴🔴🔴 |
| CF-01-B | Test template selection logic       | 🔴       | 15min | 🔴🔴🔴 |
| CF-01-C | Test data accumulation across steps | 🔴       | 15min | 🔴🔴🔴 |
| CF-01-D | Test wizard result generation       | 🔴       | 15min | 🔴🔴🔴 |
| CF-01-E | Test wizard cancellation handling   | 🔴       | 15min | 🔴🔴   |
| CF-01-F | Test wizard restart scenarios       | 🔴       | 15min | 🔴🔴   |
| CF-01-G | Run wizard tests, check coverage    | 🔴       | 10min | 🔴🔴   |
| CF-01-H | Document any uncovered paths        | 🔴       | 5min  | 🔴     |

**Total:** 8 tasks, 120min

---

#### CF-02: Create Microservice Example (PostgreSQL) (1h = 4 tasks)

| ID      | Task                                       | Priority | Time  | Impact |
| ------- | ------------------------------------------ | -------- | ----- | ------ |
| CF-02-A | Create examples/microservice-pg/ directory | 🔴       | 5min  | 🔴🔴   |
| CF-02-B | Generate microservice project with wizard  | 🔴       | 15min | 🔴🔴🔴 |
| CF-02-C | Add Docker Compose for PostgreSQL          | 🔴       | 15min | 🔴🔴🔴 |
| CF-02-D | Add README with setup instructions         | 🔴       | 15min | 🔴🔴🔴 |
| CF-02-E | Test example builds and runs               | 🔴       | 10min | 🔴🔴   |

**Total:** 5 tasks, 60min

---

#### CF-03: Create Enterprise Example (Multi-DB) (1h = 4 tasks)

| ID      | Task                                        | Priority | Time  | Impact |
| ------- | ------------------------------------------- | -------- | ----- | ------ |
| CF-03-A | Create examples/enterprise-multi/ directory | 🟡       | 5min  | 🔴🔴   |
| CF-03-B | Generate enterprise project with wizard     | 🟡       | 15min | 🔴🔴🔴 |
| CF-03-C | Add multi-database configuration            | 🟡       | 15min | 🔴🔴   |
| CF-03-D | Add audit logging example                   | 🟡       | 15min | 🔴🔴   |
| CF-03-E | Add README with architecture doc            | 🟡       | 15min | 🔴🔴   |

**Total:** 5 tasks, 60min

---

#### CF-04: Fix Commands Test Coverage to 60% (2h = 8 tasks)

| ID      | Task                              | Priority | Time  | Impact |
| ------- | --------------------------------- | -------- | ----- | ------ |
| CF-04-A | Test init command execution       | 🔴       | 15min | 🔴🔴🔴 |
| CF-04-B | Test validate command execution   | 🔴       | 15min | 🔴🔴🔴 |
| CF-04-C | Test generate command execution   | 🔴       | 15min | 🔴🔴🔴 |
| CF-04-D | Test doctor command execution     | 🔴       | 15min | 🔴🔴🔴 |
| CF-04-E | Test migrate command execution    | 🔴       | 15min | 🔴🔴🔴 |
| CF-04-F | Test command flag parsing         | 🔴       | 15min | 🔴🔴   |
| CF-04-G | Test command error handling       | 🔴       | 15min | 🔴🔴   |
| CF-04-H | Run command tests, check coverage | 🔴       | 10min | 🔴🔴   |

**Total:** 8 tasks, 120min

---

#### CF-05: Fix Adapters Test Coverage to 50% (2h = 8 tasks)

| ID      | Task                              | Priority | Time  | Impact |
| ------- | --------------------------------- | -------- | ----- | ------ |
| CF-05-A | Test file system adapter write    | 🟡       | 15min | 🔴🔴   |
| CF-05-B | Test file system adapter read     | 🟡       | 15min | 🔴🔴   |
| CF-05-C | Test file system adapter delete   | 🟡       | 15min | 🔴🔴   |
| CF-05-D | Test CLI adapter execution        | 🟡       | 15min | 🔴🔴   |
| CF-05-E | Test database adapter connection  | 🟡       | 15min | 🔴🔴   |
| CF-05-F | Test adapter error handling       | 🟡       | 15min | 🔴🔴   |
| CF-05-G | Test adapter cleanup/teardown     | 🟡       | 15min | 🔴🔴   |
| CF-05-H | Run adapter tests, check coverage | 🟡       | 10min | 🔴🔴   |

**Total:** 8 tasks, 120min

---

#### CF-06: Write Migration Guide Doc (1h = 4 tasks)

| ID      | Task                                          | Priority | Time  | Impact |
| ------- | --------------------------------------------- | -------- | ----- | ------ |
| CF-06-A | Create docs/guides/migration.md               | 🔴       | 5min  | 🔴🔴   |
| CF-06-B | Write migration from manual sqlc.yaml section | 🔴       | 15min | 🔴🔴🔴 |
| CF-06-C | Write upgrade between wizard versions section | 🔴       | 15min | 🔴🔴🔴 |
| CF-06-D | Add custom template migration section         | 🔴       | 15min | 🔴🔴   |
| CF-06-E | Review and format guide                       | 🔴       | 10min | 🔴🔴   |

**Total:** 5 tasks, 60min

---

#### CF-07: Write Troubleshooting Guide Doc (1h = 4 tasks)

| ID      | Task                                   | Priority | Time  | Impact |
| ------- | -------------------------------------- | -------- | ----- | ------ |
| CF-07-A | Create docs/guides/troubleshooting.md  | 🔴       | 5min  | 🔴🔴   |
| CF-07-B | Write common errors and solutions      | 🔴       | 20min | 🔴🔴🔴 |
| CF-07-C | Write database-specific issues section | 🔴       | 15min | 🔴🔴   |
| CF-07-D | Write sqlc integration issues section  | 🔴       | 15min | 🔴🔴   |
| CF-07-E | Review and format guide                | 🔴       | 5min  | 🔴🔴   |

**Total:** 5 tasks, 60min

**Phase 1 Total:** 43 tasks, ~10 hours

---

### ~~Phase 2: Hardening (20% → 80%) - 81 Tasks~~ DISPOSITIONED 2026-09-29 — best-practices guide done; coverage/perf/load ideas → ROADMAP "Trustworthy core" / "Reach"

#### HR-01: Complete Wizard Test Coverage to 80% (4h = 16 tasks)

| ID      | Task                              | Priority | Time  | Impact |
| ------- | --------------------------------- | -------- | ----- | ------ |
| HR-01-A | Test welcome banner display       | 🔴       | 10min | 🔴     |
| HR-01-B | Test step header display          | 🔴       | 10min | 🔴     |
| HR-01-C | Test step completion display      | 🔴       | 10min | 🔴     |
| HR-01-D | Test progress indicators          | 🔴       | 10min | 🔴     |
| HR-01-E | Test keyboard shortcuts           | 🔴       | 10min | 🔴     |
| HR-01-F | Test mouse interactions           | 🔴       | 10min | 🔴     |
| HR-01-G | Test screen resize handling       | 🔴       | 10min | 🔴     |
| HR-01-H | Test color scheme rendering       | 🔴       | 10min | 🔴     |
| HR-01-I | Test accessibility features       | 🟡       | 10min | 🔴     |
| HR-01-J | Test concurrent wizard runs       | 🟡       | 10min | 🔴     |
| HR-01-K | Test wizard state persistence     | 🟡       | 10min | 🔴     |
| HR-01-L | Test wizard configuration loading | 🔴       | 15min | 🔴     |
| HR-01-M | Test wizard template integration  | 🔴       | 15min | 🔴     |
| HR-01-N | Test wizard error recovery        | 🔴       | 15min | 🔴     |
| HR-01-O | Run wizard tests, check coverage  | 🔴       | 10min | 🔴     |
| HR-01-P | Document any remaining gaps       | 🟡       | 5min  | 🔴     |

**Total:** 16 tasks, 240min

---

#### HR-02: Complete Commands Test Coverage to 75% (3h = 12 tasks)

| ID      | Task                                 | Priority | Time  | Impact |
| ------- | ------------------------------------ | -------- | ----- | ------ |
| HR-02-A | Test init command with all flags     | 🔴       | 15min | 🔴🔴   |
| HR-02-B | Test validate command with all flags | 🔴       | 15min | 🔴🔴   |
| HR-02-C | Test generate command with all flags | 🔴       | 15min | 🔴🔴   |
| HR-02-D | Test doctor command with all flags   | 🔴       | 15min | 🔴🔴   |
| HR-02-E | Test migrate command with all flags  | 🔴       | 15min | 🔴🔴   |
| HR-02-F | Test command version flag            | 🔴       | 5min  | 🔴     |
| HR-02-G | Test command help flag               | 🔴       | 5min  | 🔴     |
| HR-02-H | Test command verbose mode            | 🔴       | 10min | 🔴     |
| HR-02-I | Test command quiet mode              | 🔴       | 10min | 🔴     |
| HR-02-J | Test command stdin input             | 🔴       | 15min | 🔴     |
| HR-02-K | Run command tests, check coverage    | 🔴       | 10min | 🔴     |
| HR-02-L | Document any uncovered paths         | 🟡       | 5min  | 🔴     |

**Total:** 12 tasks, 180min

---

#### HR-03: Complete Adapters Test Coverage to 70% (3h = 12 tasks)

| ID      | Task                              | Priority | Time  | Impact |
| ------- | --------------------------------- | -------- | ----- | ------ |
| HR-03-A | Test file system adapter mkdir    | 🟡       | 15min | 🔴🔴   |
| HR-03-B | Test file system adapter chmod    | 🟡       | 15min | 🔴     |
| HR-03-C | Test file system adapter stat     | 🟡       | 15min | 🔴🔴   |
| HR-03-D | Test file system adapter exists   | 🟡       | 10min | 🔴     |
| HR-03-E | Test CLI adapter output capture   | 🟡       | 15min | 🔴🔴   |
| HR-03-F | Test CLI adapter error capture    | 🟡       | 15min | 🔴🔴   |
| HR-03-G | Test database adapter query       | 🟡       | 15min | 🔴     |
| HR-03-H | Test database adapter transaction | 🟡       | 15min | 🔴     |
| HR-03-I | Test database adapter pooling     | 🟡       | 15min | 🔴     |
| HR-03-J | Test sqlc adapter execution       | 🔴       | 15min | 🔴🔴   |
| HR-03-K | Run adapter tests, check coverage | 🟡       | 10min | 🔴     |
| HR-03-L | Document any uncovered paths      | 🟡       | 5min  | 🔴     |

**Total:** 12 tasks, 180min

---

#### HR-04: Complete Generators Test Coverage to 80% (3h = 12 tasks)

| ID      | Task                                | Priority | Time  | Impact |
| ------- | ----------------------------------- | -------- | ----- | ------ |
| HR-04-A | Test sqlc.yaml generation           | 🔴       | 15min | 🔴🔴🔴 |
| HR-04-B | Test query file generation          | 🔴       | 15min | 🔴🔴🔴 |
| HR-04-C | Test schema file generation         | 🔴       | 15min | 🔴🔴🔴 |
| HR-04-D | Test migration file generation      | 🔴       | 15min | 🔴🔴🔴 |
| HR-04-E | Test go.mod file generation         | 🔴       | 15min | 🔴🔴   |
| HR-04-F | Test main.go file generation        | 🔴       | 15min | 🔴🔴   |
| HR-04-G | Test db package generation          | 🔴       | 15min | 🔴🔴   |
| HR-04-H | Test Dockerfile generation          | 🔴       | 15min | 🔴🔴   |
| HR-04-I | Test Makefile generation            | 🔴       | 10min | 🔴     |
| HR-04-J | Test CI/CD file generation          | 🔴       | 10min | 🔴     |
| HR-04-K | Run generator tests, check coverage | 🔴       | 10min | 🔴     |
| HR-04-L | Document any uncovered paths        | 🟡       | 5min  | 🔴     |

**Total:** 12 tasks, 180min

---

#### HR-05: Complete Creators Test Coverage to 70% (2h = 8 tasks)

| ID      | Task                               | Priority | Time  | Impact |
| ------- | ---------------------------------- | -------- | ----- | ------ |
| HR-05-A | Test directory structure creation  | 🔴       | 15min | 🔴🔴   |
| HR-05-B | Test project file generation       | 🔴       | 15min | 🔴🔴🔴 |
| HR-05-C | Test microservice project creation | 🔴       | 15min | 🔴🔴   |
| HR-05-D | Test enterprise project creation   | 🔴       | 15min | 🔴🔴   |
| HR-05-E | Test API-first project creation    | 🔴       | 15min | 🔴🔴   |
| HR-05-F | Test hobby project creation        | 🔴       | 10min | 🔴     |
| HR-05-G | Run creator tests, check coverage  | 🔴       | 10min | 🔴     |
| HR-05-H | Document any uncovered paths       | 🟡       | 5min  | 🔴     |

**Total:** 8 tasks, 120min

---

#### HR-06: Performance Baseline Testing (3h = 12 tasks)

| ID      | Task                               | Priority | Time  | Impact |
| ------- | ---------------------------------- | -------- | ----- | ------ |
| HR-06-A | Create benchmarks/wizard directory | 🟡       | 5min  | 🔴     |
| HR-06-B | Write wizard execution benchmark   | 🟡       | 15min | 🔴🔴   |
| HR-06-C | Write config generation benchmark  | 🟡       | 15min | 🔴🔴   |
| HR-06-D | Write file generation benchmark    | 🟡       | 15min | 🔴🔴   |
| HR-06-E | Run wizard execution benchmarks    | 🟡       | 10min | 🔴     |
| HR-06-F | Run config generation benchmarks   | 🟡       | 10min | 🔴     |
| HR-06-G | Run file generation benchmarks     | 🟡       | 10min | 🔴     |
| HR-06-H | Document baseline results          | 🟡       | 10min | 🔴     |
| HR-06-I | Create benchmarking README         | 🟡       | 15min | 🔴     |
| HR-06-J | Run benchmarks on multiple systems | 🟡       | 15min | 🔴     |
| HR-06-K | Add benchmark to CI/CD             | 🟡       | 10min | 🔴     |
| HR-06-L | Review and optimize hotspots       | 🟡       | 15min | 🔴     |

**Total:** 12 tasks, 180min

---

#### HR-07: Add Performance Regression Tests (2h = 8 tasks)

| ID      | Task                                   | Priority | Time  | Impact |
| ------- | -------------------------------------- | -------- | ----- | ------ |
| HR-07-A | Create regression test suite           | 🟡       | 10min | 🔴     |
| HR-07-B | Add wizard execution regression test   | 🟡       | 15min | 🔴🔴   |
| HR-07-C | Add config generation regression test  | 🟡       | 15min | 🔴🔴   |
| HR-07-D | Add file generation regression test    | 🟡       | 15min | 🔴🔴   |
| HR-07-E | Add regression test to CI/CD           | 🟡       | 10min | 🔴     |
| HR-07-F | Document regression test thresholds    | 🟡       | 10min | 🔴     |
| HR-07-G | Test regression test failure detection | 🟡       | 10min | 🔴     |
| HR-07-H | Add regression test to docs            | 🟡       | 5min  | 🔴     |

**Total:** 8 tasks, 120min

---

#### HR-08: Load Testing (100+ Tables) (3h = 12 tasks)

| ID      | Task                                 | Priority | Time  | Impact |
| ------- | ------------------------------------ | -------- | ----- | ------ |
| HR-08-A | Create large schema test fixture     | 🟡       | 15min | 🔴     |
| HR-08-B | Generate 100+ table schema           | 🟡       | 15min | 🔴     |
| HR-08-C | Create large query file test fixture | 🟡       | 15min | 🔴     |
| HR-08-D | Generate 500+ query test file        | 🟡       | 15min | 🔴     |
| HR-08-E | Run wizard with large project        | 🟡       | 15min | 🔴     |
| HR-08-F | Run config generation test           | 🟡       | 10min | 🔴     |
| HR-08-G | Run file generation test             | 🟡       | 10min | 🔴     |
| HR-08-H | Measure memory usage                 | 🟡       | 10min | 🔴     |
| HR-08-I | Measure execution time               | 🟡       | 10min | 🔴     |
| HR-08-J | Document load test results           | 🟡       | 10min | 🔴     |
| HR-08-K | Add load test to CI/CD               | 🟡       | 10min | 🔴     |
| HR-08-L | Review and optimize bottlenecks      | 🟡       | 15min | 🔴     |

**Total:** 12 tasks, 180min

---

#### HR-09: Memory Profiling & Leak Detection (2h = 8 tasks)

| ID      | Task                          | Priority | Time  | Impact |
| ------- | ----------------------------- | -------- | ----- | ------ |
| HR-09-A | Create profiling test suite   | 🟡       | 10min | 🔴     |
| HR-09-B | Run CPU profiler on wizard    | 🟡       | 15min | 🔴     |
| HR-09-C | Run memory profiler on wizard | 🟡       | 15min | 🔴     |
| HR-09-D | Run goroutine leak detector   | 🟡       | 15min | 🔴     |
| HR-09-E | Analyze profiler results      | 🟡       | 20min | 🔴     |
| HR-09-F | Fix any memory leaks          | 🟡       | 15min | 🔴     |
| HR-09-G | Verify leak fixes             | 🟡       | 15min | 🔴     |
| HR-09-H | Document profiling findings   | 🟡       | 5min  | 🔴     |

**Total:** 8 tasks, 120min

---

#### HR-10: Write Comprehensive Best Practices Guide (2h = 8 tasks)

| ID      | Task                                       | Priority | Time  | Impact |
| ------- | ------------------------------------------ | -------- | ----- | ------ |
| HR-10-A | Create docs/guides/best-practices.md       | 🔴       | 5min  | 🔴🔴   |
| HR-10-B | Write project type selection guide         | 🔴       | 15min | 🔴🔴🔴 |
| HR-10-C | Write database feature configuration guide | 🔴       | 15min | 🔴🔴🔴 |
| HR-10-D | Write performance optimization tips        | 🔴       | 15min | 🔴🔴   |
| HR-10-E | Write team collaboration guide             | 🔴       | 15min | 🔴🔴   |
| HR-10-F | Write CI/CD integration patterns           | 🔴       | 15min | 🔴🔴🔴 |
| HR-10-G | Review and format guide                    | 🔴       | 10min | 🔴🔴   |
| HR-10-H | Add examples to best practices             | 🔴       | 15min | 🔴🔴   |

**Total:** 8 tasks, 120min

---

#### HR-11: Create CI/CD Integration Examples (2h = 8 tasks)

| ID      | Task                               | Priority | Time  | Impact |
| ------- | ---------------------------------- | -------- | ----- | ------ |
| HR-11-A | Create examples/ci-cd/ directory   | 🔴       | 5min  | 🔴🔴   |
| HR-11-B | Write GitHub Actions example       | 🔴       | 15min | 🔴🔴🔴 |
| HR-11-C | Write GitLab CI example            | 🔴       | 15min | 🔴🔴   |
| HR-11-D | Write Docker Compose example       | 🔴       | 10min | 🔴🔴   |
| HR-11-E | Write Makefile integration example | 🔴       | 10min | 🔴🔴   |
| HR-11-F | Add setup instructions for each    | 🔴       | 15min | 🔴🔴   |
| HR-11-G | Test all CI/CD examples            | 🔴       | 15min | 🔴🔴   |
| HR-11-H | Review and format examples         | 🔴       | 5min  | 🔴     |

**Total:** 8 tasks, 120min

**Phase 2 Total:** 116 tasks, ~29 hours

---

## 📊 Master Task Table (All 150 Tasks)

Sorted by **(Impact × Priority) / Effort** ratio.

### 🔴 CRITICAL Tasks (High Impact, Low Effort) - 60 Tasks

| ID      | Task                                    | Time  | Impact | Phase |
| ------- | --------------------------------------- | ----- | ------ | ----- |
| QW-01-A | Analyze 3 failing test failures         | 10min | 🔴🔴🔴 | QW    |
| QW-01-B | Fix validation failure test expectation | 15min | 🔴🔴🔴 | QW    |
| QW-01-C | Fix UI panic in mock test               | 10min | 🔴🔴🔴 | QW    |
| QW-01-D | Fix data flow test expectation          | 10min | 🔴🔴🔴 | QW    |
| QW-01-E | Run full test suite to verify fixes     | 10min | 🔴🔴🔴 | QW    |
| QW-01-F | Commit integration test fixes           | 5min  | 🔴🔴🔴 | QW    |
| QW-02-C | Write quick start tutorial              | 20min | 🔴🔴🔴 | QW    |
| QW-02-D | Add common project types section        | 15min | 🔴🔴🔴 | QW    |
| QW-02-E | Add configuration options reference     | 15min | 🔴🔴🔴 | QW    |
| QW-02-F | Add troubleshooting basics              | 15min | 🔴🔴   | QW    |
| QW-03-A | Test project type step                  | 15min | 🔴🔴🔴 | QW    |
| QW-03-B | Test database step                      | 15min | 🔴🔴🔴 | QW    |
| QW-03-C | Test project details step               | 15min | 🔴🔴🔴 | QW    |
| QW-03-D | Test features step                      | 15min | 🔴🔴🔴 | QW    |
| QW-03-E | Test output step                        | 15min | 🔴🔴🔴 | QW    |
| QW-03-F | Test wizard orchestration               | 15min | 🔴🔴🔴 | QW    |
| QW-03-G | Test wizard error handling              | 15min | 🔴🔴🔴 | QW    |
| QW-04-B | Generate hobby project                  | 15min | 🔴🔴🔴 | QW    |
| QW-04-C | Add README to example                   | 15min | 🔴🔴🔴 | QW    |
| CF-01-A | Test UI helper interactions             | 15min | 🔴🔴🔴 | CF    |
| CF-01-B | Test template selection                 | 15min | 🔴🔴🔴 | CF    |
| CF-01-C | Test data accumulation                  | 15min | 🔴🔴🔴 | CF    |
| CF-01-D | Test result generation                  | 15min | 🔴🔴🔴 | CF    |
| CF-01-E | Test cancellation handling              | 15min | 🔴🔴   | CF    |
| CF-02-B | Generate microservice project           | 15min | 🔴🔴🔴 | CF    |
| CF-02-C | Add Docker Compose                      | 15min | 🔴🔴🔴 | CF    |
| CF-02-D | Add README                              | 15min | 🔴🔴🔴 | CF    |
| CF-04-A | Test init command                       | 15min | 🔴🔴🔴 | CF    |
| CF-04-B | Test validate command                   | 15min | 🔴🔴🔴 | CF    |
| CF-04-C | Test generate command                   | 15min | 🔴🔴🔴 | CF    |
| CF-04-D | Test doctor command                     | 15min | 🔴🔴🔴 | CF    |
| CF-04-E | Test migrate command                    | 15min | 🔴🔴🔴 | CF    |
| CF-04-F | Test flag parsing                       | 15min | 🔴🔴   | CF    |
| CF-04-G | Test error handling                     | 15min | 🔴🔴   | CF    |
| CF-06-B | Write migration section                 | 15min | 🔴🔴🔴 | CF    |
| CF-06-C | Write upgrade section                   | 15min | 🔴🔴🔴 | CF    |
| CF-07-B | Write common errors                     | 20min | 🔴🔴🔴 | CF    |
| CF-07-C | Write database issues                   | 15min | 🔴🔴   | CF    |
| HR-01-L | Test configuration loading              | 15min | 🔴🔴   | HR    |
| HR-01-M | Test template integration               | 15min | 🔴🔴   | HR    |
| HR-01-N | Test error recovery                     | 15min | 🔴🔴   | HR    |
| HR-02-A | Test init with all flags                | 15min | 🔴🔴   | HR    |
| HR-02-B | Test validate with all flags            | 15min | 🔴🔴   | HR    |
| HR-02-C | Test generate with all flags            | 15min | 🔴🔴   | HR    |
| HR-02-D | Test doctor with all flags              | 15min | 🔴🔴   | HR    |
| HR-02-E | Test migrate with all flags             | 15min | 🔴🔴   | HR    |
| HR-02-J | Test stdin input                        | 15min | 🔴🔴   | HR    |
| HR-03-J | Test sqlc adapter                       | 15min | 🔴🔴   | HR    |
| HR-04-A | Test sqlc.yaml generation               | 15min | 🔴🔴🔴 | HR    |
| HR-04-B | Test query file generation              | 15min | 🔴🔴🔴 | HR    |
| HR-04-C | Test schema file generation             | 15min | 🔴🔴🔴 | HR    |
| HR-04-D | Test migration generation               | 15min | 🔴🔴🔴 | HR    |
| HR-04-E | Test go.mod generation                  | 15min | 🔴🔴   | HR    |
| HR-04-F | Test main.go generation                 | 15min | 🔴🔴   | HR    |
| HR-05-A | Test directory creation                 | 15min | 🔴🔴   | HR    |
| HR-05-B | Test project file generation            | 15min | 🔴🔴🔴 | HR    |
| HR-05-C | Test microservice creation              | 15min | 🔴🔴   | HR    |
| HR-05-D | Test enterprise creation                | 15min | 🔴🔴   | HR    |
| HR-05-E | Test API-first creation                 | 15min | 🔴🔴   | HR    |
| HR-10-B | Write project type guide                | 15min | 🔴🔴🔴 | HR    |
| HR-10-C | Write database config guide             | 15min | 🔴🔴🔴 | HR    |
| HR-10-E | Write team collaboration guide          | 15min | 🔴🔴   | HR    |
| HR-10-F | Write CI/CD patterns                    | 15min | 🔴🔴🔴 | HR    |
| HR-11-B | Write GitHub Actions example            | 15min | 🔴🔴🔴 | HR    |
| HR-11-C | Write GitLab CI example                 | 15min | 🔴🔴   | HR    |
| HR-11-G | Test all CI/CD examples                 | 15min | 🔴🔴   | HR    |

### 🟡 IMPORTANT Tasks (Medium Impact, Medium Effort) - 50 Tasks

| ID      | Task                         | Time  | Impact | Phase |
| ------- | ---------------------------- | ----- | ------ | ----- |
| QW-02-A | Create user-guide directory  | 5min  | 🔴🔴   | QW    |
| QW-02-B | Write installation section   | 15min | 🔴🔴   | QW    |
| QW-02-G | Review and format doc        | 10min | 🔴🔴   | QW    |
| QW-04-A | Create example directory     | 5min  | 🔴🔴   | QW    |
| QW-04-D | Test example works           | 10min | 🔴🔴   | QW    |
| CF-01-F | Test restart scenarios       | 15min | 🔴🔴   | CF    |
| CF-02-A | Create microservice dir      | 5min  | 🔴🔴   | CF    |
| CF-02-E | Test microservice works      | 10min | 🔴🔴   | CF    |
| CF-03-A | Create enterprise dir        | 5min  | 🔴🔴   | CF    |
| CF-03-B | Generate enterprise project  | 15min | 🔴🔴   | CF    |
| CF-03-C | Add multi-DB config          | 15min | 🔴🔴   | CF    |
| CF-03-D | Add audit logging            | 15min | 🔴🔴   | CF    |
| CF-03-E | Add README                   | 15min | 🔴🔴   | CF    |
| CF-05-A | Test FS adapter write        | 15min | 🔴🔴   | CF    |
| CF-05-B | Test FS adapter read         | 15min | 🔴🔴   | CF    |
| CF-05-C | Test FS adapter delete       | 15min | 🔴🔴   | CF    |
| CF-05-D | Test CLI adapter exec        | 15min | 🔴🔴   | CF    |
| CF-05-E | Test DB adapter conn         | 15min | 🔴🔴   | CF    |
| CF-05-F | Test adapter errors          | 15min | 🔴🔴   | CF    |
| CF-05-G | Test adapter cleanup         | 15min | 🔴🔴   | CF    |
| CF-06-A | Create migration.md          | 5min  | 🔴🔴   | CF    |
| CF-06-D | Add custom template section  | 15min | 🔴🔴   | CF    |
| CF-06-E | Review migration guide       | 10min | 🔴🔴   | CF    |
| CF-07-A | Create troubleshooting.md    | 5min  | 🔴🔴   | CF    |
| CF-07-D | Add sqlc integration section | 15min | 🔴🔴   | CF    |
| CF-07-E | Review troubleshooting guide | 5min  | 🔴🔴   | CF    |
| HR-01-A | Test welcome banner          | 10min | 🔴     | HR    |
| HR-01-B | Test step headers            | 10min | 🔴     | HR    |
| HR-01-C | Test step completion         | 10min | 🔴     | HR    |
| HR-01-D | Test progress indicators     | 10min | 🔴     | HR    |
| HR-01-E | Test keyboard shortcuts      | 10min | 🔴     | HR    |
| HR-01-F | Test mouse interactions      | 10min | 🔴     | HR    |
| HR-01-G | Test screen resize           | 10min | 🔴     | HR    |
| HR-01-H | Test color rendering         | 10min | 🔴     | HR    |
| HR-02-F | Test version flag            | 5min  | 🔴     | HR    |
| HR-02-G | Test help flag               | 5min  | 🔴     | HR    |
| HR-02-H | Test verbose mode            | 10min | 🔴     | HR    |
| HR-02-I | Test quiet mode              | 10min | 🔴     | HR    |
| HR-03-A | Test FS mkdir                | 15min | 🔴🔴   | HR    |
| HR-03-B | Test FS chmod                | 15min | 🔴🔴   | HR    |
| HR-03-C | Test FS stat                 | 15min | 🔴🔴   | HR    |
| HR-03-D | Test FS exists               | 10min | 🔴     | HR    |
| HR-03-E | Test CLI output capture      | 15min | 🔴🔴   | HR    |
| HR-03-F | Test CLI error capture       | 15min | 🔴🔴   | HR    |

### 🟢 SHOULD HAVE Tasks (Lower Impact, Higher Effort) - 40 Tasks

| ID      | Task                      | Time  | Impact | Phase |
| ------- | ------------------------- | ----- | ------ | ----- |
| QW-02-H | Add screenshots           | 15min | 🔴     | QW    |
| CF-01-G | Run tests, check coverage | 10min | 🔴     | CF    |
| CF-01-H | Document uncovered paths  | 5min  | 🔴     | CF    |
| CF-04-H | Run tests, check coverage | 10min | 🔴     | CF    |
| CF-05-H | Run tests, check coverage | 10min | 🔴     | CF    |
| HR-01-I | Test accessibility        | 10min | 🔴     | HR    |
| HR-01-J | Test concurrent runs      | 10min | 🔴     | HR    |
| HR-01-K | Test state persistence    | 10min | 🔴     | HR    |
| HR-01-O | Run tests, check coverage | 10min | 🔴     | HR    |
| HR-01-P | Document gaps             | 5min  | 🔴     | HR    |
| HR-02-K | Run tests, check coverage | 10min | 🔴     | HR    |
| HR-02-L | Document gaps             | 5min  | 🔴     | HR    |
| HR-03-G | Test DB query             | 15min | 🔴     | HR    |
| HR-03-H | Test DB transaction       | 15min | 🔴     | HR    |
| HR-03-I | Test DB pooling           | 15min | 🔴     | HR    |
| HR-03-K | Run tests, check coverage | 10min | 🔴     | HR    |
| HR-03-L | Document gaps             | 5min  | 🔴     | HR    |
| HR-04-G | Test db package gen       | 15min | 🔴     | HR    |
| HR-04-H | Test Dockerfile gen       | 15min | 🔴     | HR    |
| HR-04-I | Test Makefile gen         | 10min | 🔴     | HR    |
| HR-04-J | Test CI/CD file gen       | 10min | 🔴     | HR    |
| HR-04-K | Run tests, check coverage | 10min | 🔴     | HR    |
| HR-04-L | Document gaps             | 5min  | 🔴     | HR    |
| HR-05-F | Test hobby creation       | 10min | 🔴     | HR    |
| HR-05-G | Run tests, check coverage | 10min | 🔴     | HR    |
| HR-05-H | Document gaps             | 5min  | 🔴     | HR    |
| HR-10-A | Create best-practices.md  | 5min  | 🔴🔴   | HR    |
| HR-10-D | Write perf tips           | 15min | 🔴🔴   | HR    |
| HR-10-G | Review guide              | 10min | 🔴     | HR    |
| HR-10-H | Add examples              | 15min | 🔴🔴   | HR    |
| HR-11-A | Create ci-cd directory    | 5min  | 🔴🔴   | HR    |
| HR-11-D | Write Docker Compose      | 10min | 🔴🔴   | HR    |
| HR-11-E | Write Makefile example    | 10min | 🔴🔴   | HR    |
| HR-11-F | Add setup instructions    | 15min | 🔴🔴   | HR    |
| HR-11-H | Review examples           | 5min  | 🔴🔴   | HR    |

---

## 📅 Execution Timeline (Recommended)

### Day 1: Quick Wins Phase (6 hours) - Tasks QW-01 to QW-04

**Morning (4h):**

- QW-01-A through QW-01-F: Fix 3 failing tests (1h)
- QW-02-A through QW-02-D: Create Getting Started guide (2h)
- QW-03-A through QW-03-C: Add wizard tests (1.5h)

**Afternoon (2h):**

- QW-03-D through QW-03-H: Complete wizard tests (1h)
- QW-04-A through QW-04-D: Create basic example (1h)

### Day 2-3: Critical Foundation Phase (10 hours) - Tasks CF-01 to CF-07

**Day 2 (6h):**

- CF-01-A through CF-01-E: Wizard coverage to 60% (1.5h)
- CF-02-A through CF-02-D: Microservice example (1h)
- CF-03-A through CF-03-E: Enterprise example (1.5h)
- CF-04-A through CF-04-D: Commands tests start (1h)
- CF-06-A through CF-06-C: Migration guide (1h)

**Day 3 (4h):**

- CF-04-E through CF-04-H: Commands tests complete (1.5h)
- CF-05-A through CF-05-H: Adapters tests (2h)
- CF-07-A through CF-07-E: Troubleshooting guide (30min)

### Day 4-10: Hardening Phase (21 hours) - Tasks HR-01 to HR-11

Spread across 7 days, 3 hours per day on average.

---

## ✅ Completion Criteria

### Phase 0: Quick Wins (1% → 51%)

~~- [ ] All 26 tasks completed~~ → dispositioned via the PARETO plan verdicts (archived 2026-09-29)
~~- [ ] 3 failing tests fixed~~ done — suites green (2026-09-29)
~~- [ ] Getting Started guide created~~ done at `6507a01`
~~- [ ] Wizard test coverage improved to ~30%~~ done — 34.2% on 2026-09-29 (→ TODO_LIST T6 for 60%)
~~- [ ] Basic example working~~ done at `9db82a7`

### Phase 1: Critical Foundation (4% → 64%)

~~- [ ] All 43 tasks completed~~ → dispositioned via the PARETO plan verdicts
~~- [ ] Wizard test coverage >60%~~ → TODO_LIST T6 (34.2%)
~~- [ ] 3 real-world examples working~~ → TODO_LIST T13 (hobby examples only)
~~- [ ] Commands/adapters test coverage >60%~~ → TODO_LIST T6/T7 (47.2% / 22.9%)
~~- [ ] Migration and troubleshooting guides written~~ done — `docs/MIGRATION_GUIDE.md`, `docs/TROUBLESHOOTING.md`

### Phase 2: Hardening (20% → 80%)

~~- [ ] All 81 tasks completed~~ → dispositioned via the PARETO plan verdicts
~~- [ ] Wizard test coverage >80%~~ → TODO_LIST T6 (longer-term)
~~- [ ] All package coverage >70%~~ → ROADMAP "Trustworthy core"
~~- [ ] Performance baselines established~~ → ROADMAP "Trustworthy core"
~~- [ ] Load testing completed~~ → ROADMAP "Reach"
~~- [ ] Memory profiling completed~~ → ROADMAP "Trustworthy core"
~~- [ ] Best practices guide written~~ done — `docs/BEST_PRACTICES.md`
~~- [ ] CI/CD examples provided~~ → ROADMAP "Frictionless adoption"

### Final: 100% Enterprise Ready

~~- [ ] All 150 tasks completed~~ → superseded: open work lives in TODO_LIST.md / ROADMAP.md (2026-09-29)
~~- [ ] Overall test coverage >70%~~ → ROADMAP "Trustworthy core"

- [ ] All tests passing (100%)
- [ ] Documentation complete
- [ ] Examples working
- [ ] Performance validated
- [ ] Ready for v1.0.0 release

---

**Created:** 2026-01-13_17-06
**Total Tasks:** 150
**Total Time:** ~37.5 hours
**Priority:** Impact/Effort sorted
**Status:** 🟢 READY TO EXECUTE

---

## Resolution (2026-09-29)

This 150-task plan was superseded by actual execution across 2026 sessions. Phase-level verdicts are inline above; per-task dispositions for the shared 27-task core are in the archived PARETO plan (`docs/planning/archived/2026-01-13_17-06-PARETO_OPTIMAL_EXECUTION_PLAN.md`), and all surviving work is consolidated in TODO_LIST.md and ROADMAP.md.
