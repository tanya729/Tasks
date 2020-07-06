Feature: Working via REST API
  As HTTP client of REST API service
  In order to understand that an event was added/changed into database
  I want to get the event via the REST API

  Scenario: Adding event into database via REST API
    When I send Create request to API
    Then added event will be returned with id of the event
    And GetError has no errors in both cases

  Scenario: Getting the event from database by id of the event via REST API
    When I send Get request with event id to API
    Then I get response with event
    And GetError has no errors

  Scenario: Getting non existing event from database by id via REST API
    When I send Get request with non existing event id to API
    Then I get response with error code 'Event not found'