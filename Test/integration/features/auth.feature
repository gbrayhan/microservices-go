Feature: User Login and Token Refresh
  As a registered user
  I want to authenticate and refresh my tokens
  So that I can access protected endpoints

  Scenario: POST /login with valid credentials returns tokens and user info
    When I send a POST request to "/login" with body:
      """
      {
        "email": "${START_USER_EMAIL}",
        "password": "${START_USER_PW}"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "accessToken"
    And the JSON response should contain key "refreshToken"
    And the JSON response should contain "email": "${START_USER_EMAIL}"
    And I save the JSON response key "accessToken" as "accessToken"
    And I save the JSON response key "refreshToken" as "refreshToken"

  Scenario: POST /login with invalid credentials returns 401
    When I send a POST request to "/login" with body:
      """
      {
        "email": "${START_USER_EMAIL}",
        "password": "wrongpassword"
      }
      """
    Then the response code should be 401
    And the JSON response should contain error "error": "Invalid credentials"

  Scenario: POST /access-token/refresh with valid refresh token returns new access token
    When I send a POST request to "/access-token/refresh" with body:
      """
      {
        "refreshToken": "${refreshToken}"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "accessToken"
    And the JSON response should contain key "id"
    And the JSON response should contain key "email"
    And I save the JSON response key "accessToken" as "accessToken"

  Scenario: POST /access-token/refresh with invalid refresh token returns 401
    When I send a POST request to "/access-token/refresh" with body:
      """
      {
        "refreshToken": "someInvalidToken"
      }
      """
    Then the response code should be 401
    And the JSON response should contain error "error": "Invalid token"

  Scenario: Access protected endpoint without token
    Given I clear the authentication token
    When I send a GET request to "/api/medicines/1"
    Then the response code should be 401
    And the JSON response should contain error "error": "Authorization header not provided"

  # Re-authenticate so subsequent scenarios have a valid token
  Scenario: Re-authenticate after clearing the token
    When I send a POST request to "/login" with body:
      """
      {
        "email": "${START_USER_EMAIL}",
        "password": "${START_USER_PW}"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "accessToken"
    And the JSON response should contain key "refreshToken"
    And the JSON response should contain "email": "${START_USER_EMAIL}"
    And I save the JSON response key "accessToken" as "accessToken"
    And I save the JSON response key "refreshToken" as "refreshToken"
