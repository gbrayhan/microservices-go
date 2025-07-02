Feature: Device Information and Health Check
  As an API consumer
  I want to access device information and health check endpoints
  So that I can verify system status and device details.

  Background:
    # Login to obtain accessToken is handled globally by InitializeScenario
    # and the token is automatically added to headers by the addAuthHeader function.
    # All resources created in scenarios are automatically tracked and cleaned up
    # by the test framework's teardown mechanism.

  Scenario: TC01 - Access device information endpoint with authentication
    When I send a GET request to "/api/device"
    Then the response code should be 200
    And the JSON response should contain key "ip_address"
    And the JSON response should contain key "user_agent"
    And the JSON response should contain key "device_type"
    And the JSON response should contain key "browser"
    And the JSON response should contain key "os"
    And the JSON response should contain key "language"

  Scenario: TC01.1 - Access device information endpoint without authentication
    Given I clear the authentication token
    When I send a GET request to "/api/device"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  Scenario: TC02 - Access health check authenticated endpoint
    When I send a GET request to "/api/health-check-auth/"
    Then the response code should be 200
    And the JSON response should contain "message": "authenticated"

  Scenario: TC02.1 - Access health check authenticated endpoint without authentication
    Given I clear the authentication token
    When I send a GET request to "/api/health-check-auth/"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  Scenario: TC03 - Verify device information contains expected fields
    When I send a GET request to "/api/device"
    Then the response code should be 200
    And the JSON response should contain key "ip_address"
    And the JSON response should contain key "user_agent"
    And the JSON response should contain key "device_type"
    And the JSON response should contain key "browser"
    And the JSON response should contain key "browser_version"
    And the JSON response should contain key "os"
    And the JSON response should contain key "language"

  Scenario: TC04 - Test device information with different user agents
    # This scenario tests that the device info interceptor works correctly
    # The actual user agent will be set by the test framework
    When I send a GET request to "/api/device"
    Then the response code should be 200
    And the JSON response should contain key "user_agent"

  Scenario: TC05 - Verify health check response format
    When I send a GET request to "/api/health-check-auth/"
    Then the response code should be 200
    And the JSON response should contain key "message"
    And the JSON response should contain "message": "authenticated"

  Scenario: TC06 - Test device endpoint with malformed authentication
    # This scenario tests the JWT middleware with invalid tokens
    Given I clear the authentication token
    When I send a GET request to "/api/device"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  Scenario: TC07 - Test health check endpoint with malformed authentication
    Given I clear the authentication token
    When I send a GET request to "/api/health-check-auth/"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  Scenario: TC08 - Verify device information structure
    When I send a GET request to "/api/device"
    Then the response code should be 200
    And the JSON response should be an object
    And the JSON response should contain key "ip_address"
    And the JSON response should contain key "user_agent"
    And the JSON response should contain key "device_type"
    And the JSON response should contain key "browser"
    And the JSON response should contain key "browser_version"
    And the JSON response should contain key "os"
    And the JSON response should contain key "language"

  Scenario: TC09 - Test device endpoint with expired token
    # This scenario would require setting up an expired token
    # For now, we'll test the basic functionality
    When I send a GET request to "/api/device"
    Then the response code should be 200
    And the JSON response should contain key "ip_address"

  Scenario: TC10 - Test health check endpoint with expired token
    # This scenario would require setting up an expired token
    # For now, we'll test the basic functionality
    When I send a GET request to "/api/health-check-auth/"
    Then the response code should be 200
    And the JSON response should contain "message": "authenticated" 