Feature: User Login
  As a registered user
  I want to authenticate with valid credentials
  So that I receive access and refresh tokens

  Scenario: POST /v1/auth/login with valid credentials returns tokens and user info
    Given the service is initialized
    When I send a POST request to "/v1/auth/login" with body:
      """
      {
        "email": "gbrayhan@gmail.com",
        "password": "qweqwe"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "security"
    And the JSON response should contain key "data"
    And the JSON response should contain "data.email": "gbrayhan@gmail.com"
    And the JSON response should contain "security.jwtAccessToken": "*"
    And the JSON response should contain "security.jwtRefreshToken": "*"

  Scenario: POST /v1/auth/login with invalid credentials returns error
    Given the service is initialized
    When I send a POST request to "/v1/auth/login" with body:
      """
      {
        "email": "gbrayhan@gmail.com",
        "password": "wrongpass"
      }
      """
    Then the response code should be 403
    And the JSON response should contain key "error"
