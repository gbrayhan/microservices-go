Feature: User Login
  As a registered user
  I want to authenticate with valid credentials
  So that I receive access and refresh tokens

  Scenario: POST /login with valid credentials returns tokens and user info
    Given the service is initialized
    When I send a POST request to "/login" with body:
      """
      {
        "email": "gbrayhan@gmail.com",
        "password": "qweqwe"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "accessToken"
    And the JSON response should contain key "refreshToken"
    And the JSON response should contain "userEmail": "gbrayhan@gmail.com"
