Feature: Settings controls
  Settings controls should preserve user choices while keeping dense forms easy
  to scan.

  @SCN-SETTINGS-001 @unit
  Scenario: Selected multi-select options stay first
    Given a settings multi-select has selected and unselected options
    When the options are ordered for display
    Then selected options are shown first
    And the relative order inside selected and unselected groups is preserved

  @SCN-SETTINGS-002 @api
  Scenario: Admin manages tags
    Given the admin is signed in
    When a tag is created, renamed, listed, and deleted
    Then each response reflects the requested tag state
    And the deleted tag is no longer listed

  @SCN-SETTINGS-003 @api
  Scenario: Admin manages language aliases
    Given the admin is signed in
    When a language is created, updated, listed, and deleted
    Then each response reflects the requested language state
    And the deleted language is no longer listed

  @SCN-SETTINGS-004 @api
  Scenario: Admin manages an indexer configuration
    Given the admin is signed in
    When an indexer is created, updated, listed, tested, and deleted
    Then each response reflects the requested indexer state
    And the indexer test succeeds against the local mock provider

  @SCN-SETTINGS-005 @api
  Scenario: Admin validates a download client configuration
    Given the admin is signed in
    When a download client configuration is submitted for testing
    Then the validation response reports the local mock result

  @SCN-SETTINGS-006 @unit
  Scenario: Settings forms normalize optional fields
    Given a settings form has whitespace and blank optional values
    When the form is normalized for the API
    Then required values are trimmed
    And blank optional values are omitted

  @SCN-SETTINGS-007 @unit
  Scenario: File naming templates show useful token guidance
    Given a file naming template query or template
    When suggestions or examples are requested
    Then matching tokens are ranked for the user
    And known tokens render with realistic example values

  @SCN-SETTINGS-008 @unit
  Scenario: Quality size controls preserve valid ranges
    Given quality size slider values are edited
    When values are converted and validated
    Then invalid ranges produce a user-facing error
    And valid slider output preserves minimum, preferred, and maximum ordering

  @SCN-SETTINGS-009 @unit
  Scenario: UI API helpers send normalized requests
    Given the UI calls settings and media API helpers
    When the generated API client returns success, empty data, or errors
    Then helpers return safe defaults, throw user-facing errors, and send normalized request bodies

  @SCN-SETTINGS-010 @unit
  Scenario: Language and quality catalogs expose stable user choices
    Given language and quality settings catalogs are loaded
    When options are prepared for selectors and grouped controls
    Then common languages stay first with display codes
    And unknown selected languages remain visible
    And qualities are grouped by their resolution family

  @SCN-SETTINGS-011 @unit
  Scenario: Storage normalizes settings inputs
    Given settings inputs contain whitespace, duplicates, and invalid option values
    When they are normalized before persistence
    Then names and aliases are trimmed and deduplicated
    And invalid media profile and monitoring options fall back or fail predictably

  @SCN-SETTINGS-012 @unit
  Scenario: UI settings helpers preserve editable state
    Given settings records are loaded into edit forms
    When the user saves custom formats, media profiles, languages, users, or provider links
    Then copied form state does not mutate the source records
    And normalized request payloads trim optional values and preserve valid choices

  @SCN-SETTINGS-013 @integration
  Scenario: Storage persists managed users
    Given the settings database is available
    When managed users are created, updated, queried, and deleted
    Then password hashes verify correctly
    And duplicate usernames and missing users return domain errors

  @SCN-SETTINGS-014 @integration
  Scenario: Storage persists external integration settings
    Given the settings database is available
    When indexers, download clients, and metadata providers are changed
    Then list, health, session token, and delete side effects are persisted

  @SCN-SETTINGS-015 @integration
  Scenario: Storage records system events and search cache state
    Given the settings database is available
    When system events, indexer cache entries, metadata cache entries, and search history are recorded
    Then list, stats, lookup, and delete operations reflect the stored side effects

  @SCN-SETTINGS-016 @api
  Scenario: Admin manages library folders and path mappings
    Given the admin is signed in
    When library folder options, folders, scans, and path mappings are changed
    Then each API response reflects the filesystem and persisted mapping side effects

  @SCN-SETTINGS-017 @api
  Scenario: Admin manages custom formats and tests release parsing
    Given the admin is signed in
    When a custom format is created, updated, listed, parsed against a release name, and deleted
    Then each API response reflects the requested custom format state
    And parsing reports the matching custom format conditions

  @SCN-SETTINGS-018 @api
  Scenario: Admin manages persisted download clients
    Given the admin is signed in
    When a download client is created, updated, listed, tested, and deleted
    Then each API response reflects the requested download client state
    And the test endpoint reports the local validation result

  @SCN-SETTINGS-019 @api
  Scenario: Admin updates file naming templates
    Given the admin is signed in
    When file naming templates are read and updated
    Then the API response persists trimmed template values

  @SCN-SETTINGS-020 @api
  Scenario: Admin inspects and clears indexer search cache
    Given the admin is signed in
    And indexer search cache and history entries exist
    When indexer search settings, cache entries, and history are read or cleared
    Then the API responses expose counts and deleted side effects

  @SCN-SETTINGS-021 @api
  Scenario: Admin inspects and clears metadata cache
    Given the admin is signed in
    And metadata cache and history entries exist
    When metadata cache entries and history are read or cleared
    Then the API responses expose counts and deleted side effects

  @SCN-SETTINGS-022 @api
  Scenario: Admin manages metadata providers
    Given the admin is signed in
    And a local metadata provider mock is available
    When a metadata provider is created, updated, listed, tested, and deleted
    Then each API response reflects the requested metadata provider state
    And the test endpoint reports the local metadata validation result

  @SCN-SETTINGS-023 @api
  Scenario: Admin manages media profiles and quality sizes
    Given the admin is signed in
    And seeded quality choices are available
    When media profiles and quality size settings are created, updated, listed, and deleted
    Then each API response reflects the requested profile and quality size state

  @SCN-SETTINGS-024 @api
  Scenario: Admin manages application users
    Given the admin is signed in
    When a managed user is created, updated, listed, and deleted
    Then each API response reflects the requested managed user state

  @SCN-SETTINGS-025 @api
  Scenario: Admin manages DLNA renderer profiles
    Given the admin is signed in
    And seeded DLNA renderer profiles are available
    When renderer profiles and device overrides are created, updated, reset, imported, exported, and deleted
    Then each API response reflects the requested DLNA profile state
