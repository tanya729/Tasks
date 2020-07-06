Feature: Working via Grpc API
  As Grpc client of Grpc API service
  In order to understand that an event was added/changed into database
  I want to get the event via the Grpc API

  Scenario: Adding event into database via Grpc API
    When gRPC I send Create request to API
    Then gRPC added event will be returned with id of the event

  Scenario: Getting the event from database by id of the event via Grpc API
    When gRPC I send Get request with event id to API
    Then gRPC I get response with event

  Scenario: Getting non existing event from database by id via Grpc API
    When gRPC I send Get request with non existing event id to API
    Then gRPC I get response with error code 'Event not found