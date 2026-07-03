Feature: Activity queue
  Activity rows should expose clear release details and available actions.

  @SCN-ACTIVITY-001 @unit
  Scenario: Activity row summarizes a downloading release
    Given a download activity has release title, status, and progress
    When the row display model is built
    Then year, languages, quality, formats, and progress are shown
    And cancellable, importable, and deletable actions match the status

  @SCN-ACTIVITY-002 @unit @integration
  Scenario: Activity actions update visible queue and media state
    Given an activity is visible in the download queue
    When the user cancels or deletes the activity
    Then the queue reflects the updated activity state
    And affected media items are refreshed after cancellation
