Feature: notifier work correctly
  In order to understand that logic for notify work
  I want to get event from queue

  Scenario: Adding events by REST API
    When I send request to API with event
    Then I get events from queue
    Then I have event in Events
    And I has no errors
