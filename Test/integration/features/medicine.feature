Feature: Medicine Management
  As an API consumer
  I want to manage medicines
  So that I can perform CRUD operations

  Background:
    # Authentication handled globally

  Scenario: Create a new medicine successfully
    Given I generate a unique alias as "medName"
    When I send a POST request to "/v1/medicine" with body:
      """
      {
        "name": "${medName}",
        "description": "Test medicine",
        "eanCode": "${medName}-EAN",
        "laboratory": "TestLab"
      }
      """
    Then the response code should be 200
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "medicineID"

  Scenario: Retrieve created medicine
    When I send a GET request to "/v1/medicine/${medicineID}"
    Then the response code should be 200
    And the JSON response should contain "id" with numeric value  ${medicineID}

  Scenario: Update medicine description
    When I send a PUT request to "/v1/medicine/${medicineID}" with body:
      """
      {
        "description": "Updated description"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "description": "Updated description"

  Scenario: Delete medicine
    When I send a DELETE request to "/v1/medicine/${medicineID}"
    Then the response code should be 200
    And the JSON response should contain "message": "resource deleted successfully"

  Scenario: Search medicines paginated
    When I send a GET request to "/v1/medicine/search?page=1&pageSize=10"
    Then the response code should be 200
    And the JSON response should contain key "data"
