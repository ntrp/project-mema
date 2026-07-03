Feature: External integrations
  Indexer and metadata integrations should be tested against realistic local
  mocks so provider changes are exercised at the boundary.

  @SCN-INTEGRATIONS-001 @integration @api
  Scenario: Torznab indexer capabilities are read from a local mock
    Given a local Torznab mock exposes capabilities and releases
    When the indexer connection is tested
    Then the test succeeds
    And the reported category count matches the mock capabilities

  @SCN-INTEGRATIONS-002 @integration
  Scenario: Metadata provider search can use a local mock
    Given a local metadata provider mock exposes movie and series results
    When metadata search is executed through the provider boundary
    Then normalized media results are returned

  @SCN-INTEGRATIONS-003 @unit
  Scenario: Download clients queue, report, and cancel downloads
    Given SABnzbd and Transmission clients return realistic API payloads
    When downloads are added, inspected, and cancelled through the integration service
    Then download ids, progress, files, and cancel messages reflect the provider responses
