Feature: List events from database via GRPC API for day, week and month
  In order to understand that logic for list events works as expected
  I want to get events list for a day, week and month

  Scenario: Adding events by GRPC API
    When I send request to GRPC for cycle with 3 events for day, week, and month
    And gRPC has no errors in these cases

  Scenario: List events for day via GRPC API
    When gRPC I send request with type day
    Then gRPC I get Events with 1 events
    And gRPC I have not errors

  Scenario: List events for week via GRPC API
    When gRPC I send request with type week
    Then gRPC I get Events with 2 events
    And gRPC I have not errors

  Scenario: List events for month via GRPC API
    When gRPC I send  request with type month
    Then gRPC I get Events with 3 events
    And gRPC I have not errors