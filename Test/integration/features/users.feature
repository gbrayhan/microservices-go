Feature: User Management
  As an API consumer
  I want to manage users
  So that I can perform CRUD operations

  Background:
    # Authentication handled globally

  Scenario: Create a new user successfully
    Given I generate a unique alias as "userName"
    When I send a POST request to "/v1/user" with body:
      """
      {
        "user": "${userName}",
        "email": "${userName}@test.com",
        "firstName": "Test",
        "lastName": "User",
        "password": "pass123",
        "role": "user"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "userID"

  Scenario: Retrieve created user
    When I send a GET request to "/v1/user/${userID}"
    Then the response code should be 200
    And the JSON response should contain "id" with numeric value  ${userID}

  Scenario: Update user first name
    When I send a PUT request to "/v1/user/${userID}" with body:
      """
      {
        "firstName": "Updated"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "firstName": "Updated"

  Scenario: Delete user
    When I send a DELETE request to "/v1/user/${userID}"
    Then the response code should be 200
    And the JSON response should contain "message": "resource deleted successfully"

  Scenario: Search users paginated
    When I send a GET request to "/v1/user/search?page=1&pageSize=5"
    Then the response code should be 200
    And the JSON response should contain key "data"
