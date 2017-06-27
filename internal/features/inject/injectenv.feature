@injectfeature
Feature: Parameter injection
  Scenario:
    Given some configuration store in the SSM parameter store
    When I create an inject environment
    Then I can read my environment variables based on prefix
