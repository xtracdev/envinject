@injectfeature
Feature: Parameter injection
  Scenario:
    Given some configuration store in the SSM parameter store
    When I create an inject environment
    Then I can read my environment variables based on prefix

  Scenario:
    Given a mix of paramstore and environment variables
    When I create an injected environment
    And there are some environment vars not injected
    Then I can access the non-injected variables
