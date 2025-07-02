Feature: Medicine Management
  As an API consumer
  I want to manage medicines
  So that I can perform CRUD operations and search for them independently.

  Background:
    # Login to obtain accessToken is handled globally by InitializeScenario
    # and the token is automatically added to headers by the addAuthHeader function.
    # All resources created in scenarios are automatically tracked and cleaned up
    # by the test framework's teardown mechanism.

  Scenario: TC01 - Create a new medicine successfully
    Given I generate a unique EAN code as "newMedicineEan"
    When I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${newMedicineEan}",
        "description": "Ibuprofen 600mg Test Tablets",
        "laboratory": "PharmaTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Ibuprofen",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 30.0,
        "unitType": "tablet",
        "rxCode": "RXP001"
      }
      """
    Then the response code should be 201
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "medicineID"
    And the JSON response should contain "eanCode": "${newMedicineEan}"
    And the JSON response should contain "description": "Ibuprofen 600mg Test Tablets"

  Scenario: TC01.1 - Attempt to create a medicine with a duplicate EAN code
    Given I generate a unique EAN code as "duplicateTestEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${duplicateTestEan}",
        "description": "First Ibuprofen for duplicate test",
        "laboratory": "DuplicateTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Ibuprofen",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 10.0,
        "unitType": "tablet",
        "rxCode": "RXP001D-First"
      }
      """
    Then the response code should be 201

    When I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${duplicateTestEan}",
        "description": "Another Ibuprofen with same EAN",
        "laboratory": "DuplicateTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Ibuprofen",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 10.0,
        "unitType": "tablet",
        "rxCode": "RXP001D-Second"
      }
      """
    Then the response code should be 500 or 409
    And the JSON response field "error" should contain string "Could not create medicine"
    And the JSON response field "error" should contain string "duplicate code"

  Scenario: TC01.2 - Attempt to create a medicine with missing required fields
    When I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "MEDTESTXYZ002-MISSING",
        "laboratory": "MissingFields Labs"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC02 - Retrieve the created medicine
    Given I generate a unique EAN code as "retrieveMedicineEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${retrieveMedicineEan}",
        "description": "Medicine for retrieval test",
        "laboratory": "RetrieveTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Acetaminophen",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 20.0,
        "unitType": "tablet",
        "rxCode": "RXP002"
      }
      """
    And I save the JSON response key "id" as "retrieveMedicineID"
    When I send a GET request to "/api/medicines/${retrieveMedicineID}"
    Then the response code should be 200
    And the JSON response should contain key "description"
    And the JSON response should contain key "eanCode"

  Scenario: TC02.1 - Attempt to retrieve a non-existent medicine
    When I send a GET request to "/api/medicines/999999"
    Then the response code should be 404
    And the JSON response should contain error "error": "Medicine not found"

  Scenario: TC02.2 - Attempt to retrieve a medicine with an invalid ID format
    When I send a GET request to "/api/medicines/invalidIDFormat"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid ID"

  Scenario: TC03 - Update the existing medicine
    Given I generate a unique EAN code as "updateMedicineEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${updateMedicineEan}",
        "description": "Medicine for update test",
        "laboratory": "UpdateTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Aspirin",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 25.0,
        "unitType": "tablet",
        "rxCode": "RXP003"
      }
      """
    And I save the JSON response key "id" as "updateMedicineID"
    And I generate a unique EAN code as "updatedMedicineEan"
    When I send a PUT request to "/api/medicines/${updateMedicineID}" with body:
      """
      {
        "eanCode": "${updatedMedicineEan}",
        "description": "Ibuprofen 600mg Tablets - Updated",
        "laboratory": "PharmaTest Labs Inc.",
        "type": "tablet",
        "iva": "16UPD",
        "satKey": "51182201_UPD",
        "activeIngredient": "Ibuprofen Forte",
        "temperatureControl": "refrigerated",
        "isControlled": true,
        "unitQuantity": 32.0,
        "unitType": "tablet",
        "rxCode": "RXP002_UPD"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "description": "Ibuprofen 600mg Tablets - Updated"
    And the JSON response should contain "eanCode": "${updatedMedicineEan}"
    And the JSON response should contain "isControlled": true

  Scenario: TC03.1 - Attempt to update a non-existent medicine
    When I send a PUT request to "/api/medicines/999999" with body:
      """
      {
        "eanCode": "MEDNONEXIST001",
        "description": "Non Existent Update",
        "laboratory": "NoLab",
        "type": "capsule",
        "unitQuantity": 10.0,
        "unitType": "capsule"
      }
      """
    Then the response code should be 404
    And the JSON response should contain error "error": "Medicine not found"

  Scenario: TC03.2 - Update medicine with partial fields (only description)
    Given I generate a unique EAN code as "partialUpdateEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${partialUpdateEan}",
        "description": "Original description for partial update test",
        "laboratory": "PartialUpdateTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Acetaminophen",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 20.0,
        "unitType": "tablet",
        "rxCode": "RXP008"
      }
      """
    And I save the JSON response key "id" as "partialUpdateMedicineID"
    When I send a PUT request to "/api/medicines/${partialUpdateMedicineID}" with body:
      """
      {
        "description": "Updated description only"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "description": "Updated description only"
    And the JSON response should contain "eanCode": "${partialUpdateEan}"
    And the JSON response should contain "laboratory": "PartialUpdateTest Labs"
    And the JSON response should contain "isControlled": false

  Scenario: TC03.3 - Update medicine with multiple partial fields
    Given I generate a unique EAN code as "multiPartialUpdateEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${multiPartialUpdateEan}",
        "description": "Original description for multi partial update",
        "laboratory": "OriginalLab",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Original Ingredient",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 15.0,
        "unitType": "tablet",
        "rxCode": "RXP009"
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateMedicineID"
    When I send a PUT request to "/api/medicines/${multiPartialUpdateMedicineID}" with body:
      """
      {
        "description": "Updated description",
        "laboratory": "UpdatedLab",
        "isControlled": true,
        "unitQuantity": 25.0
      }
      """
    Then the response code should be 200
    And the JSON response should contain "description": "Updated description"
    And the JSON response should contain "laboratory": "UpdatedLab"
    And the JSON response should contain "isControlled": true
    And the JSON response should contain "unitQuantity": 25.0
    And the JSON response should contain "eanCode": "${multiPartialUpdateEan}"
    And the JSON response should contain "type": "tablet"
    And the JSON response should contain "activeIngredient": "Original Ingredient"

  Scenario: TC03.4 - Attempt to update medicine with empty body
    Given I generate a unique EAN code as "emptyUpdateEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${emptyUpdateEan}",
        "description": "Medicine for empty update test",
        "laboratory": "EmptyUpdateTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Test Ingredient",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 10.0,
        "unitType": "tablet",
        "rxCode": "RXP010"
      }
      """
    And I save the JSON response key "id" as "emptyUpdateMedicineID"
    When I send a PUT request to "/api/medicines/${emptyUpdateMedicineID}" with body:
      """
      {
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "No fields to update"

  Scenario: TC03.5 - Update medicine with invalid type in partial update
    Given I generate a unique EAN code as "invalidTypeUpdateEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${invalidTypeUpdateEan}",
        "description": "Medicine for invalid type update test",
        "laboratory": "InvalidTypeTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Test Ingredient",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 10.0,
        "unitType": "tablet",
        "rxCode": "RXP011"
      }
      """
    And I save the JSON response key "id" as "invalidTypeUpdateMedicineID"
    When I send a PUT request to "/api/medicines/${invalidTypeUpdateMedicineID}" with body:
      """
      {
        "type": "invalid_type"
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid medicine type, must be one of: injection, tablet, capsule"

  Scenario: TC03.6 - Update medicine with invalid temperature control in partial update
    Given I generate a unique EAN code as "invalidTempUpdateEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${invalidTempUpdateEan}",
        "description": "Medicine for invalid temp update test",
        "laboratory": "InvalidTempTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Test Ingredient",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 10.0,
        "unitType": "tablet",
        "rxCode": "RXP012"
      }
      """
    And I save the JSON response key "id" as "invalidTempUpdateMedicineID"
    When I send a PUT request to "/api/medicines/${invalidTempUpdateMedicineID}" with body:
      """
      {
        "temperatureControl": "invalid_temperature"
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid temperature control, must be one of: room, refrigerated, frozen"

  Scenario: TC03.7 - Update medicine with invalid unit type in partial update
    Given I generate a unique EAN code as "invalidUnitUpdateEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${invalidUnitUpdateEan}",
        "description": "Medicine for invalid unit update test",
        "laboratory": "InvalidUnitTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Test Ingredient",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 10.0,
        "unitType": "tablet",
        "rxCode": "RXP013"
      }
      """
    And I save the JSON response key "id" as "invalidUnitUpdateMedicineID"
    When I send a PUT request to "/api/medicines/${invalidUnitUpdateMedicineID}" with body:
      """
      {
        "unitType": "invalid_unit"
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid unit type, must be one of: ml, g, piece, tablet, capsule"

  Scenario: TC04 - Search for medicines by description (paginated)
    Given I generate a unique EAN code as "searchMedicineEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${searchMedicineEan}",
        "description": "Search Test Ibuprofen",
        "laboratory": "SearchTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Ibuprofen",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 15.0,
        "unitType": "tablet",
        "rxCode": "RXP004"
      }
      """
    When I send a GET request to "/api/medicines/search-paginated?description_like=Ibuprofen&page=1&limit=5"
    Then the response code should be 200
    And the JSON response should contain key "medicines"
    And the JSON response should contain key "total_records"

  Scenario: TC04.1 - Search for medicines by EAN code (paginated)
    Given I generate a unique EAN code as "searchEanCode"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${searchEanCode}",
        "description": "EAN Search Test Medicine",
        "laboratory": "EANSearchTest Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Paracetamol",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 18.0,
        "unitType": "tablet",
        "rxCode": "RXP005"
      }
      """
    When I send a GET request to "/api/medicines/search-paginated?ean_code_like=${searchEanCode}"
    Then the response code should be 200
    And the JSON response should contain key "medicines"

  Scenario: TC05 - Search for medicine property coincidences
    Given I generate a unique EAN code as "propertySearchEan"
    And I send a POST request to "/api/medicines" with body:
      """
      {
        "eanCode": "${propertySearchEan}",
        "description": "Property Search Test Medicine",
        "laboratory": "PharmaTest Property Labs",
        "type": "tablet",
        "iva": "16",
        "satKey": "51182200",
        "activeIngredient": "Omeprazole",
        "temperatureControl": "room",
        "isControlled": false,
        "unitQuantity": 12.0,
        "unitType": "tablet",
        "rxCode": "RXP006"
      }
      """
    When I send a GET request to "/api/medicines/search-by-property?property=laboratory&search_text=PharmaTest"
    Then the response code should be 200

  Scenario: TC05.1 - Search for medicine property coincidences with invalid property
    When I send a GET request to "/api/medicines/search-by-property?property=invalidProp&search_text=Test"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid property"
