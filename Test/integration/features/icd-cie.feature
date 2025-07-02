Feature: ICD-CIE Management
  As an API consumer
  I want to manage ICD-CIE records
  So that I can perform CRUD operations and search for them independently.

  Background:
    # Login to obtain accessToken is handled globally by InitializeScenario
    # and the token is automatically added to headers by the addAuthHeader function.
    # All resources created in scenarios are automatically tracked and cleaned up
    # by the test framework's teardown mechanism.

  Scenario: TC01 - Create a new ICD-CIE record successfully
    Given I generate a unique alias as "newCieCode"
    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${newCieCode}",
        "description": "Test ICD-CIE record for integration testing",
        "chapterNo": "I",
        "chapterTitle": "Certain infectious and parasitic diseases"
      }
      """
    Then the response code should be 201
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "icdCieID"
    And the JSON response should contain "cieVersion": "CIE-10"
    And the JSON response should contain "code": "${newCieCode}"
    And the JSON response should contain "description": "Test ICD-CIE record for integration testing"
    And the JSON response should contain "chapterNo": "I"
    And the JSON response should contain "chapterTitle": "Certain infectious and parasitic diseases"

  Scenario: TC01.1 - Attempt to create an ICD-CIE record with missing required fields
    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "description": "ICD-CIE record without required fields",
        "chapterNo": "II"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC01.2 - Attempt to create an ICD-CIE record with duplicate code
    Given I generate a unique alias as "duplicateCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${duplicateCieCode}",
        "description": "First ICD-CIE record for duplicate test",
        "chapterNo": "III",
        "chapterTitle": "Diseases of the blood and blood-forming organs"
      }
      """
    Then the response code should be 201

    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${duplicateCieCode}",
        "description": "Second ICD-CIE record with same code",
        "chapterNo": "IV",
        "chapterTitle": "Endocrine, nutritional and metabolic diseases"
      }
      """
    Then the response code should be 500 or 409
    And the JSON response should contain key "error"

  Scenario: TC02 - Retrieve all ICD-CIE records
    When I send a GET request to "/api/icd-cie"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC03 - Retrieve a specific ICD-CIE record
    Given I generate a unique alias as "retrieveCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${retrieveCieCode}",
        "description": "ICD-CIE record for retrieval test",
        "chapterNo": "V",
        "chapterTitle": "Mental and behavioural disorders"
      }
      """
    And I save the JSON response key "id" as "retrieveIcdCieID"
    When I send a GET request to "/api/icd-cie/${retrieveIcdCieID}"
    Then the response code should be 200
    And the JSON response should contain "code": "${retrieveCieCode}"
    And the JSON response should contain "description": "ICD-CIE record for retrieval test"
    And the JSON response should contain key "id"

  Scenario: TC03.1 - Attempt to retrieve a non-existent ICD-CIE record
    When I send a GET request to "/api/icd-cie/999999"
    Then the response code should be 404
    And the JSON response should contain error "error": "ICDCie record not found"

  Scenario: TC03.2 - Attempt to retrieve an ICD-CIE record with invalid ID format
    When I send a GET request to "/api/icd-cie/invalidID"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid ID"

  Scenario: TC04 - Update an existing ICD-CIE record
    Given I generate a unique alias as "updateCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${updateCieCode}",
        "description": "ICD-CIE record for update test",
        "chapterNo": "VI",
        "chapterTitle": "Diseases of the nervous system"
      }
      """
    And I save the JSON response key "id" as "updateIcdCieID"
    And I generate a unique alias as "updatedCieCode"
    When I send a PUT request to "/api/icd-cie/${updateIcdCieID}" with body:
      """
      {
        "cieVersion": "CIE-11",
        "code": "${updatedCieCode}",
        "description": "Updated ICD-CIE record description",
        "chapterNo": "VII",
        "chapterTitle": "Diseases of the eye and adnexa"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "cieVersion": "CIE-11"
    And the JSON response should contain "code": "${updatedCieCode}"
    And the JSON response should contain "description": "Updated ICD-CIE record description"
    And the JSON response should contain "chapterNo": "VII"
    And the JSON response should contain "chapterTitle": "Diseases of the eye and adnexa"

  Scenario: TC04.1 - Update ICD-CIE record with partial fields (only description)
    Given I generate a unique alias as "partialUpdateCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${partialUpdateCieCode}",
        "description": "Original description for partial update test",
        "chapterNo": "VI",
        "chapterTitle": "Diseases of the nervous system"
      }
      """
    And I save the JSON response key "id" as "partialUpdateIcdCieID"
    When I send a PUT request to "/api/icd-cie/${partialUpdateIcdCieID}" with body:
      """
      {
        "description": "Updated description only"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "description": "Updated description only"
    And the JSON response should contain "cieVersion": "CIE-10"
    And the JSON response should contain "code": "${partialUpdateCieCode}"
    And the JSON response should contain "chapterNo": "VI"
    And the JSON response should contain "chapterTitle": "Diseases of the nervous system"

  Scenario: TC04.2 - Update ICD-CIE record with multiple partial fields
    Given I generate a unique alias as "multiPartialUpdateCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${multiPartialUpdateCieCode}",
        "description": "Original description for multi partial update",
        "chapterNo": "VI",
        "chapterTitle": "Original chapter title"
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateIcdCieID"
    When I send a PUT request to "/api/icd-cie/${multiPartialUpdateIcdCieID}" with body:
      """
      {
        "cieVersion": "CIE-11",
        "description": "Updated description",
        "chapterNo": "VII"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "cieVersion": "CIE-11"
    And the JSON response should contain "description": "Updated description"
    And the JSON response should contain "chapterNo": "VII"
    And the JSON response should contain "code": "${multiPartialUpdateCieCode}"
    And the JSON response should contain "chapterTitle": "Original chapter title"

  Scenario: TC04.3 - Update ICD-CIE record with empty request
    Given I generate a unique alias as "emptyUpdateCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${emptyUpdateCieCode}",
        "description": "ICD-CIE record for empty update test",
        "chapterNo": "VI",
        "chapterTitle": "Diseases of the nervous system"
      }
      """
    And I save the JSON response key "id" as "emptyUpdateIcdCieID"
    When I send a PUT request to "/api/icd-cie/${emptyUpdateIcdCieID}" with body:
      """
      {
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "No fields to update"

  Scenario: TC04.4 - Attempt to update a non-existent ICD-CIE record
    When I send a PUT request to "/api/icd-cie/999999" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "NONEXISTENT",
        "description": "Non-existent ICD-CIE record",
        "chapterNo": "VIII",
        "chapterTitle": "Diseases of the ear and mastoid process"
      }
      """
    Then the response code should be 404
    And the JSON response should contain error "error": "ICDCie record not found"

  Scenario: TC05 - Delete an ICD-CIE record
    Given I generate a unique alias as "deleteCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${deleteCieCode}",
        "description": "ICD-CIE record to be deleted",
        "chapterNo": "IX",
        "chapterTitle": "Diseases of the circulatory system"
      }
      """
    And I save the JSON response key "id" as "deleteIcdCieID"
    When I send a DELETE request to "/api/icd-cie/${deleteIcdCieID}"
    Then the response code should be 200
    And the JSON response should contain "message": "ICDCie record deleted successfully"

  Scenario: TC06 - Search ICD-CIE records paginated
    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=10"
    Then the response code should be 200
    And the JSON response should contain key "current_page"
    And the JSON response should contain key "records"
    And the JSON response should contain key "page_size"
    And the JSON response should contain key "total_pages"
    And the JSON response should contain key "total_records"

  Scenario: TC06.1 - Search ICD-CIE records with filters
    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=5&cie_version_like=CIE&description_like=test"
    Then the response code should be 200
    And the JSON response should contain key "current_page"
    And the JSON response should contain key "records"
    And the JSON response should contain key "page_size"
    And the JSON response should contain key "total_pages"
    And the JSON response should contain key "total_records"

  Scenario: TC06.2 - Search ICD-CIE records with exact matches
    Given I generate a unique alias as "exactMatchCieCode"
    And I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10-EXACT",
        "code": "${exactMatchCieCode}",
        "description": "ICD-CIE record for exact match test",
        "chapterNo": "X",
        "chapterTitle": "Diseases of the respiratory system"
      }
      """
    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=10&cie_version_match=CIE-10-EXACT&code_match=${exactMatchCieCode}"
    Then the response code should be 200
    And the JSON response should contain key "records"
    And the JSON response should contain key "total_records"

  Scenario: TC07 - Search ICD-CIE coincidences by property
    When I send a GET request to "/api/icd-cie/search-by-property?property=cie_version&search_text=CIE"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC07.1 - Search ICD-CIE coincidences by code property
    When I send a GET request to "/api/icd-cie/search-by-property?property=code&search_text=TEST"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC07.2 - Search ICD-CIE coincidences by description property
    When I send a GET request to "/api/icd-cie/search-by-property?property=description&search_text=test"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC07.3 - Search ICD-CIE coincidences by chapter number property
    When I send a GET request to "/api/icd-cie/search-by-property?property=chapter_no&search_text=I"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC07.4 - Search ICD-CIE coincidences by chapter title property
    When I send a GET request to "/api/icd-cie/search-by-property?property=chapter_title&search_text=diseases"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC07.5 - Attempt to search with invalid property
    When I send a GET request to "/api/icd-cie/search-by-property?property=invalid_property&search_text=test"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid property or search_text"

  Scenario: TC07.6 - Attempt to search with empty search text
    When I send a GET request to "/api/icd-cie/search-by-property?property=code&search_text="
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid property or search_text"

  Scenario: TC08 - Create multiple ICD-CIE records for comprehensive testing
    Given I generate a unique alias as "multiCieCode1"
    And I generate a unique alias as "multiCieCode2"
    And I generate a unique alias as "multiCieCode3"
    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${multiCieCode1}",
        "description": "First multi-test ICD-CIE record",
        "chapterNo": "XI",
        "chapterTitle": "Diseases of the digestive system"
      }
      """
    Then the response code should be 201

    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${multiCieCode2}",
        "description": "Second multi-test ICD-CIE record",
        "chapterNo": "XII",
        "chapterTitle": "Diseases of the skin and subcutaneous tissue"
      }
      """
    Then the response code should be 201

    When I send a POST request to "/api/icd-cie" with body:
      """
      {
        "cieVersion": "CIE-10",
        "code": "${multiCieCode3}",
        "description": "Third multi-test ICD-CIE record",
        "chapterNo": "XIII",
        "chapterTitle": "Diseases of the musculoskeletal system and connective tissue"
      }
      """
    Then the response code should be 201

  Scenario: TC09 - Test pagination with different page sizes
    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=1"
    Then the response code should be 200
    And the JSON response should contain "page_size": 1

    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=5"
    Then the response code should be 200
    And the JSON response should contain "page_size": 5

    When I send a GET request to "/api/icd-cie/search-paginated?page=2&limit=3"
    Then the response code should be 200
    And the JSON response should contain "current_page": 2
    And the JSON response should contain "page_size": 3

  Scenario: TC10 - Test edge cases for ICD-CIE operations
    When I send a GET request to "/api/icd-cie/search-paginated?page=0&limit=10"
    Then the response code should be 200
    And the JSON response should contain "current_page": 1

    When I send a GET request to "/api/icd-cie/search-paginated?page=1&limit=0"
    Then the response code should be 200
    And the JSON response should contain "page_size": 10

    When I send a GET request to "/api/icd-cie/search-paginated?page=-1&limit=-5"
    Then the response code should be 200
    And the JSON response should contain "current_page": 1
    And the JSON response should contain "page_size": 10 