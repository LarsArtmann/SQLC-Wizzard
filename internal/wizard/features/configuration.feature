Feature: SQLC Configuration Generation
  As a developer
  I want to generate sqlc configurations
  So that I can quickly setup type-safe SQL

  Scenario: Generate microservice configuration
    Given I choose microservice template
    And I select PostgreSQL database
    And I set project name to "my-api"
    And I set output directory to "./generated"
    When I run sqlc-wizard init
    Then sqlc.yaml should be created
    And it should contain PostgreSQL configuration
    And example queries should be generated
    And directories should be created

  Scenario: Generate enterprise configuration
    Given I choose enterprise template
    And I select PostgreSQL database
    And I enable strict functions
    And I enable UUID support
    When I run sqlc-wizard init
    Then sqlc.yaml should contain enterprise settings
    And UUID types should be enabled
    And strict validation should be enabled

  Scenario: Configuration validation errors
    Given I provide empty project name
    When I attempt to generate configuration
    Then validation should fail
    And clear error message should be shown
    And no files should be created

  Scenario: Template selection workflow
    Given I start sqlc-wizard init
    When I view available templates
    Then I should see all project types:
      | hobby         |
      | microservice  |
      | enterprise     |
      | api-first     |
      | analytics     |
      | testing       |
      | multi-tenant  |
      | library       |
    And I should see all database types:
      | postgresql    |
      | mysql         |
      | sqlite        |

  Scenario: File generation validation
    Given I have valid configuration
    When I run file generation
    Then sqlc.yaml should be valid YAML
    And it should be parseable by sqlc
    And generated Go code should compile
    And test files should be created

  Scenario: Database-specific configurations
    Given I select MySQL database
    When I generate configuration
    Then MySQL driver settings should be included
    And MySQL-specific types should be generated

  Scenario: Output directory handling
    Given output directory does not exist
    When I generate configuration
    Then output directory should be created
    And it should have correct permissions

  Scenario: Package path validation
    Given I provide invalid package path
    When I attempt to generate configuration
    Then validation should fail
    And package path format error should be shown

  Scenario: Feature toggle handling
    Given I enable JSON tags
    And I enable prepared queries
    When I generate configuration
    Then JSON tags should be enabled in generated code
    And prepared query methods should be generated

  Scenario: Schema generation workflow
    Given I have SQL schema files
    When I generate Go types
    Then Go types should match schema
    And table relationships should be preserved
    And nullable types should be handled correctly