Feature: Media detail experience
  Media detail views should summarize files, releases, and navigation targets in
  predictable ways.

  @SCN-MEDIA-001 @unit
  Scenario: Media files are grouped for display
    Given a media item has file paths, metadata, and a quality profile
    When the file display model is built
    Then movie and series rows expose relative paths, episode numbers, codecs, languages, size, and missing rows

  @SCN-MEDIA-002 @unit
  Scenario: Release results can be filtered and sorted
    Given release search results contain torrent, nzb, warning, and error matches
    When filters and sort options are applied
    Then the visible results match source, quality, size, and score constraints
    And severe match problems stay ahead of normal results

  @SCN-MEDIA-003 @unit
  Scenario: Media hero and navigation helpers expose stable labels
    Given media details and app section identifiers
    When hero labels and navigation hrefs are requested
    Then images, runtime text, status labels, monitor hints, and routes are stable

  @SCN-MEDIA-004 @unit
  Scenario: App routes and discovery sections map observable state
    Given route paths, media metadata, and discover blacklist entries
    When the app shell route and discovery display models are derived
    Then selected views, submenu sections, related sections, and filtered results match the route and blacklist state

  @SCN-MEDIA-005 @api
  Scenario: Admin manages the discovery blacklist
    Given the admin is signed in
    When media is added to, listed from, filtered by, and removed from the discovery blacklist
    Then blacklisted media is omitted from discovery results
    And each API response reflects the persisted blacklist state

  @SCN-MEDIA-006 @api
  Scenario: Signed-in users create and inspect media requests
    Given a user is signed in
    When a media request is created, listed, and fetched
    Then each API response reflects the requested media and pending status

  @SCN-MEDIA-007 @api
  Scenario: Admin manages media item monitoring options
    Given the admin is signed in
    When a media item is created, listed, and updated
    Then each API response reflects the media item and monitoring settings

  @SCN-MEDIA-008 @api
  Scenario: Signed-in users search and inspect provider metadata
    Given a signed-in user and a local metadata provider mock are available
    When media search, autocomplete, advanced search, discovery, and details endpoints are used
    Then provider results and metadata details are returned from the local mock

  @SCN-MEDIA-009 @api
  Scenario: Admin manages download activity lifecycle
    Given the admin is signed in
    And a queued download activity exists
    When download activity is listed, cancelled, and deleted
    Then each API response reflects the activity lifecycle side effects

  @SCN-MEDIA-010 @unit
  Scenario: Library discovery classifies local media files
    Given a library folder contains movies, series episodes, hidden files, and non-video files
    When the local scanner discovers media files
    Then discovered files expose stable titles, years, media kinds, paths, and safe-match flags

  @SCN-MEDIA-011 @integration
  Scenario: Storage persists release search snapshots
    Given the settings database is available
    When release search candidates and provider errors are replaced for a media item
    Then listed results are ordered for user review
    And stale releases and errors are removed from the current snapshot

  @SCN-MEDIA-012 @unit
  Scenario: Media request surfaces expose list and approval state
    Given pending and approved media requests exist
    When the requests surface renders list, detail, and empty states
    Then request cards, detail facts, tags, approval controls, and missing-request messages are visible

  @SCN-MEDIA-013 @integration
  Scenario: Storage lists wanted media by observable availability state
    Given the settings database is available
    When monitored, unmonitored, and actively downloading media items exist
    Then only monitored media without local files or active downloads is listed as missing
