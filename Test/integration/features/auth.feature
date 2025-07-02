Feature: User Login and Token Refresh
  As a registered user
  I want to authenticate and refresh my tokens
  So that I can access protected endpoints

  Scenario: POST /login with valid credentials returns tokens and user info
    When I send a POST request to "/v1/auth/login" with body:
      """
      {
        "email": "${START_USER_EMAIL}",
        "password": "${START_USER_PW}"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "security"
    And the JSON response should contain "data.email" with value "${START_USER_EMAIL}"
    And I save the JSON response key "security.jwtAccessToken" as "accessToken"
    And I save the JSON response key "security.jwtRefreshToken" as "refreshToken"

  Scenario: POST /login with invalid credentials returns 401
    When I send a POST request to "/v1/auth/login" with body:
      """
      {
        "email": "${START_USER_EMAIL}",
        "password": "wrongpassword"
      }
      """
    Then the response code should be 401
    And the JSON response should contain error "error": "email or password does not match"

  Scenario: POST /access-token/refresh with valid refresh token returns new access token
    When I send a POST request to "/v1/auth/access-token" with body:
      """
      {
        "refreshToken": "${refreshToken}"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "security"
    And the JSON response should contain key "data"
    And I save the JSON response key "security.jwtAccessToken" as "accessToken"

  Scenario: POST /access-token/refresh with invalid refresh token returns 401
    When I send a POST request to "/v1/auth/access-token" with body:
      """
      {
        "refreshToken": "someInvalidToken"
      }
      """
    Then the response code should be 401
    And the JSON response should contain error "error": "token contains an invalid number of segments"

  Scenario: Access protected endpoint without token
    Given I clear the authentication token
    When I send a GET request to "/v1/medicine/1"
    Then the response code should be 401
    And the JSON response should contain error "error": "Token not provided"

  # Re-authenticate so subsequent scenarios have a valid token
  Scenario: Re-authenticate after clearing the token
    When I send a POST request to "/v1/auth/login" with body:
      """
      {
        "email": "${START_USER_EMAIL}",
        "password": "${START_USER_PW}"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "security"
    And the JSON response should contain "data.email" with value "${START_USER_EMAIL}"
    And I save the JSON response key "security.jwtAccessToken" as "accessToken"
    And I save the JSON response key "security.jwtRefreshToken" as "refreshToken"
