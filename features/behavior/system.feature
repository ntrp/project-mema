Feature: System runtime
  Runtime utilities should report and route observable system state reliably.

  @SCN-SYSTEM-001 @unit
  Scenario: Configuration loads environment overrides
    Given environment variables are set for runtime configuration
    When configuration is loaded
    Then explicit values override defaults
    And invalid optional values fall back safely

  @SCN-SYSTEM-002 @unit
  Scenario: Event subscribers receive published events
    Given an event subscriber is active
    When an event is published
    Then the subscriber receives the event data
    And cancelled subscribers stop receiving new events

  @SCN-SYSTEM-003 @unit
  Scenario: Rate limit headers produce retry delays
    Given a provider returns retry or reset headers
    When the headers are inspected
    Then the user-facing delay is derived from the first valid header

  @SCN-SYSTEM-004 @unit
  Scenario: Static web handler falls back to the app shell
    Given a built web app has an app shell and static asset
    When root, asset, or unknown routes are requested
    Then root and unknown routes serve the app shell
    And existing assets are served directly

  @SCN-SYSTEM-005 @unit
  Scenario: Tool detection reports missing tools
    Given a required external tool is not on the path
    When tool detection runs
    Then the result reports the tool as unavailable with an error

  @SCN-SYSTEM-006 @api
  Scenario: Admin inspects system status and runtime log settings
    Given the admin is signed in
    When system status, log level, log file settings, and event settings are read or updated
    Then the API responses expose the current runtime state and persisted settings

  @SCN-SYSTEM-009 @api
  Scenario: Runtime system endpoints expose health, tools, and protected log files
    Given the admin is signed in
    When health, tool status, log level updates, and log file downloads are requested
    Then public health is available
    And tool status is listed for the signed-in user
    And invalid log levels and unavailable log files return safe errors

  @SCN-SYSTEM-007 @unit
  Scenario: Runtime log manager publishes and persists entries
    Given runtime logging is configured with subscribers and file output
    When structured log entries are emitted
    Then subscribers receive buffered entries
    And log files are listed and protected from unsafe paths

  @SCN-SYSTEM-008 @unit
  Scenario: Admin observes runtime jobs, events, and logs
    Given the admin opens system observability settings
    When jobs, events, and log views have no loaded runtime data yet
    Then each view still exposes its live status, filters, and empty-state guidance
    And destructive actions remain explicit confirmation flows
