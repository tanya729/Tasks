Feature: notifier work correctly
  In order to understand that logic for notify work
  I want to get event from queue

  Scenario: Adding events by GRPC API
    When I send request to GRPC API with event
    Then gRPC I get events from queue
    Then gRPC I have event in Events
    And gRPC I has no errors
