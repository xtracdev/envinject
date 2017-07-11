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

  Scenario:
    Given a var defined both in the param store and the environment
    When I lookup the var
    Then the param store value is returned

  Scenario:
    Given a mixed environment
    When I enumerate the vars in the environment
    And the same value is in both the env and the parame store
    Then the param store vars values are returned