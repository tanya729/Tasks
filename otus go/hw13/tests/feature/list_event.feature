Feature: List events from database via REST API for day, week and month
  In order to understand that logic for list events works as expected
  I want to get events list for a day, week and month

  Scenario: Adding events by REST API
    When I send request to API for cycle with 3 events for day, week, and month
    And has no errors in these cases

  Scenario: List events for day via REST API
    When I send request with type day
    Then I get Events with 1 events
    And I have not errors

  Scenario: List events for week via REST API
    When I send request with type week
    Then I get Events with 2 events
    And I have not errors

  Scenario: List events for month via REST API
    When I send  request with type month
    Then I get Events with 3 events
    And I have not errors