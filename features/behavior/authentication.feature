Feature: Authentication
  Admin users need session-based access so protected app areas remain private.

  @SCN-AUTH-001 @api @e2e @critical
  Scenario: Anonymous visitor sees the login experience
    Given no session cookie is present
    When the visitor opens the app
    Then the current session is unauthenticated
    And the login form is shown

  @SCN-AUTH-002 @api @e2e @critical
  Scenario: Admin signs in with valid credentials
    Given the default admin user exists
    When the admin signs in with valid credentials
    Then the current session is authenticated
    And protected navigation is available

  @SCN-AUTH-003 @api
  Scenario: Invalid credentials are rejected
    Given the default admin user exists
    When a visitor signs in with an invalid password
    Then the login request is rejected
    And no authenticated session is created
