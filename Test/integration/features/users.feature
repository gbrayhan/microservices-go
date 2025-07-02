Feature: User Management
  As an API consumer
  I want to manage users, roles and devices
  So that I can perform CRUD operations and search for them independently.

  Background:
    # Login to obtain accessToken is handled globally by InitializeScenario
    # and the token is automatically added to headers by the addAuthHeader function.
    # All resources created in scenarios are automatically tracked and cleaned up
    # by the test framework's teardown mechanism.

  # ===== ROLES MANAGEMENT =====

  Scenario: TC01 - Create a new role successfully
    Given I generate a unique alias as "newRoleName"
    When I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${newRoleName}",
        "description": "Test role for integration testing",
        "enabled": true
      }
      """
    Then the response code should be 201
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "roleID"
    And the JSON response should contain "name": "${newRoleName}"
    And the JSON response should contain "description": "Test role for integration testing"
    And the JSON response should contain "enabled": true

  Scenario: TC01.1 - Attempt to create a role with missing required fields
    When I send a POST request to "/api/users/roles" with body:
      """
      {
        "description": "Role without name"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC02 - Retrieve all roles
    When I send a GET request to "/api/users/roles"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC03 - Retrieve a specific role
    Given I generate a unique alias as "retrieveRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${retrieveRoleName}",
        "description": "Role for retrieval test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "retrieveRoleID"
    When I send a GET request to "/api/users/roles/${retrieveRoleID}"
    Then the response code should be 200
    And the JSON response should contain "name": "${retrieveRoleName}"
    And the JSON response should contain key "id"

  Scenario: TC03.1 - Attempt to retrieve a non-existent role
    When I send a GET request to "/api/users/roles/999999"
    Then the response code should be 404
    And the JSON response should contain error "error": "Role not found"

  Scenario: TC03.2 - Attempt to retrieve a role with invalid ID format
    When I send a GET request to "/api/users/roles/invalidID"
    Then the response code should be 400
    And the JSON response should contain error "error": "Invalid ID"

  Scenario: TC04 - Update an existing role
    Given I generate a unique alias as "updateRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${updateRoleName}",
        "description": "Role for update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "updateRoleID"
    And I generate a unique alias as "updatedRoleName"
    When I send a PUT request to "/api/users/roles/${updateRoleID}" with body:
      """
      {
        "name": "${updatedRoleName}",
        "description": "Updated role description",
        "enabled": false
      }
      """
    Then the response code should be 200
    And the JSON response should contain "name": "${updatedRoleName}"
    And the JSON response should contain "description": "Updated role description"
    And the JSON response should contain "enabled": false

  Scenario: TC04.1 - Attempt to update a non-existent role
    When I send a PUT request to "/api/users/roles/999999" with body:
      """
      {
        "name": "NonExistentRole",
        "description": "This role doesn't exist",
        "enabled": true
      }
      """
    Then the response code should be 404
    And the JSON response should contain error "error": "Role not found"

  Scenario: TC05 - Delete a role
    Given I generate a unique alias as "deleteRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${deleteRoleName}",
        "description": "Role to be deleted",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deleteRoleID"
    When I send a DELETE request to "/api/users/roles/${deleteRoleID}"
    Then the response code should be 200
    And the JSON response should contain "message": "Role deleted successfully"

  # ===== USERS MANAGEMENT =====

  Scenario: TC06 - Create a new user successfully
    Given I generate a unique alias as "newUserUsername"
    And I generate a unique alias as "newUserEmail"
    And I generate a unique alias as "userRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${userRoleName}",
        "description": "Role for user creation test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "userRoleID"
    When I send a POST request to "/api/users" with body:
      """
      {
        "username": "${newUserUsername}",
        "firstName": "John",
        "lastName": "Doe",
        "email": "${newUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Developer",
        "roleId": ${userRoleID},
        "enabled": true
      }
      """
    Then the response code should be 201
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "userID"
    And the JSON response should contain "username": "${newUserUsername}"
    And the JSON response should contain "email": "${newUserEmail}@test.com"
    And the JSON response should contain "firstName": "John"
    And the JSON response should contain "lastName": "Doe"

  Scenario: TC06.1 - Attempt to create a user with missing required fields
    When I send a POST request to "/api/users" with body:
      """
      {
        "firstName": "John",
        "lastName": "Doe"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC07 - Retrieve all users
    When I send a GET request to "/api/users"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC08 - Retrieve a specific user
    Given I generate a unique alias as "retrieveUserUsername"
    And I generate a unique alias as "retrieveUserEmail"
    And I generate a unique alias as "retrieveUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${retrieveUserRoleName}",
        "description": "Role for user retrieval test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "retrieveUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${retrieveUserUsername}",
        "firstName": "Jane",
        "lastName": "Smith",
        "email": "${retrieveUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Manager",
        "roleId": ${retrieveUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "retrieveUserID"
    When I send a GET request to "/api/users/${retrieveUserID}"
    Then the response code should be 200
    And the JSON response should contain "username": "${retrieveUserUsername}"
    And the JSON response should contain "email": "${retrieveUserEmail}@test.com"
    And the JSON response should contain key "role"

  Scenario: TC08.1 - Attempt to retrieve a non-existent user
    When I send a GET request to "/api/users/999999"
    Then the response code should be 404
    And the JSON response should contain error "error": "User not found"

  Scenario: TC09 - Update an existing user
    Given I generate a unique alias as "updateUserUsername"
    And I generate a unique alias as "updateUserEmail"
    And I generate a unique alias as "updateUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${updateUserRoleName}",
        "description": "Role for user update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "updateUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${updateUserUsername}",
        "firstName": "Bob",
        "lastName": "Johnson",
        "email": "${updateUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Analyst",
        "roleId": ${updateUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "updateUserID"
    And I generate a unique alias as "updatedUserUsername"
    And I generate a unique alias as "updatedUserEmail"
    When I send a PUT request to "/api/users/${updateUserID}" with body:
      """
      {
        "username": "${updatedUserUsername}",
        "firstName": "Robert",
        "lastName": "Johnson",
        "email": "${updatedUserEmail}@test.com",
        "password": "newSecurePassword123",
        "jobPosition": "Senior Analyst",
        "roleId": ${updateUserRoleID},
        "enabled": false
      }
      """
    Then the response code should be 200
    And the JSON response should contain "username": "${updatedUserUsername}"
    And the JSON response should contain "email": "${updatedUserEmail}@test.com"
    And the JSON response should contain "firstName": "Robert"
    And the JSON response should contain "lastName": "Johnson"
    And the JSON response should contain "jobPosition": "Senior Analyst"
    And the JSON response should contain "enabled": false

  Scenario: TC09.1 - Update user with partial fields (only firstName)
    Given I generate a unique alias as "partialUpdateUserUsername"
    And I generate a unique alias as "partialUpdateUserEmail"
    And I generate a unique alias as "partialUpdateUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${partialUpdateUserRoleName}",
        "description": "Role for partial user update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "partialUpdateUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${partialUpdateUserUsername}",
        "firstName": "Original",
        "lastName": "User",
        "email": "${partialUpdateUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Original Position",
        "roleId": ${partialUpdateUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "partialUpdateUserID"
    When I send a PUT request to "/api/users/${partialUpdateUserID}" with body:
      """
      {
        "firstName": "Updated"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "firstName": "Updated"
    And the JSON response should contain "lastName": "User"
    And the JSON response should contain "jobPosition": "Original Position"
    And the JSON response should contain "enabled": true

  Scenario: TC09.2 - Update user with multiple partial fields
    Given I generate a unique alias as "multiPartialUpdateUserUsername"
    And I generate a unique alias as "multiPartialUpdateUserEmail"
    And I generate a unique alias as "multiPartialUpdateUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${multiPartialUpdateUserRoleName}",
        "description": "Role for multi partial user update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${multiPartialUpdateUserUsername}",
        "firstName": "Original",
        "lastName": "User",
        "email": "${multiPartialUpdateUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Original Position",
        "roleId": ${multiPartialUpdateUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateUserID"
    When I send a PUT request to "/api/users/${multiPartialUpdateUserID}" with body:
      """
      {
        "firstName": "Updated",
        "lastName": "UpdatedUser",
        "jobPosition": "Updated Position",
        "enabled": false
      }
      """
    Then the response code should be 200
    And the JSON response should contain "firstName": "Updated"
    And the JSON response should contain "lastName": "UpdatedUser"
    And the JSON response should contain "jobPosition": "Updated Position"
    And the JSON response should contain "enabled": false
    And the JSON response should contain "email": "${multiPartialUpdateUserEmail}@test.com"

  Scenario: TC09.3 - Update user with empty request
    Given I generate a unique alias as "emptyUpdateUserUsername"
    And I generate a unique alias as "emptyUpdateUserEmail"
    And I generate a unique alias as "emptyUpdateUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${emptyUpdateUserRoleName}",
        "description": "Role for empty user update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "emptyUpdateUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${emptyUpdateUserUsername}",
        "firstName": "Empty",
        "lastName": "User",
        "email": "${emptyUpdateUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Empty Position",
        "roleId": ${emptyUpdateUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "emptyUpdateUserID"
    When I send a PUT request to "/api/users/${emptyUpdateUserID}" with body:
      """
      {
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "No fields to update"

  Scenario: TC10 - Delete a user
    Given I generate a unique alias as "deleteUserUsername"
    And I generate a unique alias as "deleteUserEmail"
    And I generate a unique alias as "deleteUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${deleteUserRoleName}",
        "description": "Role for user deletion test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deleteUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${deleteUserUsername}",
        "firstName": "Alice",
        "lastName": "Brown",
        "email": "${deleteUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Tester",
        "roleId": ${deleteUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deleteUserID"
    When I send a DELETE request to "/api/users/${deleteUserID}"
    Then the response code should be 200
    And the JSON response should contain "message": "user deleted successfully"

  # ===== DEVICES MANAGEMENT =====

  Scenario: TC11 - Create a new device successfully
    Given I generate a unique alias as "deviceUserUsername"
    And I generate a unique alias as "deviceUserEmail"
    And I generate a unique alias as "deviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${deviceUserRoleName}",
        "description": "Role for device test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${deviceUserUsername}",
        "firstName": "Device",
        "lastName": "User",
        "email": "${deviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Device Tester",
        "roleId": ${deviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deviceUserID"
    When I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${deviceUserID},
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "device_type": "desktop",
        "browser": "Chrome",
        "browser_version": "120.0.0.0",
        "os": "Windows",
        "language": "en-US"
      }
      """
    Then the response code should be 201
    And the JSON response should contain key "id"
    And I save the JSON response key "id" as "deviceID"
    And the JSON response should contain "ip_address": "192.168.1.100"
    And the JSON response should contain "device_type": "desktop"
    And the JSON response should contain "browser": "Chrome"

  Scenario: TC11.1 - Attempt to create a device with missing required fields
    When I send a POST request to "/api/users/devices" with body:
      """
      {
        "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        "device_type": "desktop"
      }
      """
    Then the response code should be 400
    And the JSON response should contain key "error"

  Scenario: TC12 - Retrieve devices by user ID
    Given I generate a unique alias as "devicesUserUsername"
    And I generate a unique alias as "devicesUserEmail"
    And I generate a unique alias as "devicesUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${devicesUserRoleName}",
        "description": "Role for devices retrieval test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "devicesUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${devicesUserUsername}",
        "firstName": "Devices",
        "lastName": "User",
        "email": "${devicesUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Devices Tester",
        "roleId": ${devicesUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "devicesUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${devicesUserID},
        "ip_address": "192.168.1.101",
        "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
        "device_type": "desktop",
        "browser": "Safari",
        "browser_version": "17.0",
        "os": "macOS",
        "language": "en-US"
      }
      """
    When I send a GET request to "/api/users/devices/user-id/${devicesUserID}"
    Then the response code should be 200
    And the JSON response should be an array

  Scenario: TC13 - Retrieve a specific device
    Given I generate a unique alias as "specificDeviceUserUsername"
    And I generate a unique alias as "specificDeviceUserEmail"
    And I generate a unique alias as "specificDeviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${specificDeviceUserRoleName}",
        "description": "Role for specific device test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "specificDeviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${specificDeviceUserUsername}",
        "firstName": "Specific",
        "lastName": "DeviceUser",
        "email": "${specificDeviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Specific Device Tester",
        "roleId": ${specificDeviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "specificDeviceUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${specificDeviceUserID},
        "ip_address": "192.168.1.102",
        "user_agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
        "device_type": "desktop",
        "browser": "Firefox",
        "browser_version": "115.0",
        "os": "Linux",
        "language": "en-US"
      }
      """
    And I save the JSON response key "id" as "specificDeviceID"
    When I send a GET request to "/api/users/devices/${specificDeviceID}"
    Then the response code should be 200
    And the JSON response should contain "ip_address": "192.168.1.102"
    And the JSON response should contain "browser": "Firefox"
    And the JSON response should contain "os": "Linux"

  Scenario: TC14 - Update a device
    Given I generate a unique alias as "updateDeviceUserUsername"
    And I generate a unique alias as "updateDeviceUserEmail"
    And I generate a unique alias as "updateDeviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${updateDeviceUserRoleName}",
        "description": "Role for device update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "updateDeviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${updateDeviceUserUsername}",
        "firstName": "Update",
        "lastName": "DeviceUser",
        "email": "${updateDeviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Update Device Tester",
        "roleId": ${updateDeviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "updateDeviceUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${updateDeviceUserID},
        "ip_address": "192.168.1.103",
        "user_agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15",
        "device_type": "mobile",
        "browser": "Safari",
        "browser_version": "17.0",
        "os": "iOS",
        "language": "en-US"
      }
      """
    And I save the JSON response key "id" as "updateDeviceID"
    When I send a PUT request to "/api/users/devices/${updateDeviceID}" with body:
      """
      {
        "ip_address": "192.168.1.104",
        "user_agent": "Updated User Agent String",
        "device_type": "tablet",
        "browser": "Chrome Mobile",
        "browser_version": "120.0.0.0",
        "os": "Android",
        "language": "es-ES"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "ip_address": "192.168.1.104"
    And the JSON response should contain "device_type": "tablet"
    And the JSON response should contain "browser": "Chrome Mobile"
    And the JSON response should contain "os": "Android"
    And the JSON response should contain "language": "es-ES"

  Scenario: TC14.1 - Update device with partial fields (only ip_address)
    Given I generate a unique alias as "partialUpdateDeviceUserUsername"
    And I generate a unique alias as "partialUpdateDeviceUserEmail"
    And I generate a unique alias as "partialUpdateDeviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${partialUpdateDeviceUserRoleName}",
        "description": "Role for partial device update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "partialUpdateDeviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${partialUpdateDeviceUserUsername}",
        "firstName": "Partial",
        "lastName": "DeviceUser",
        "email": "${partialUpdateDeviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Partial Device Tester",
        "roleId": ${partialUpdateDeviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "partialUpdateDeviceUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${partialUpdateDeviceUserID},
        "ip_address": "192.168.1.105",
        "user_agent": "Original User Agent",
        "device_type": "desktop",
        "browser": "Original Browser",
        "browser_version": "1.0.0",
        "os": "Original OS",
        "language": "en-US"
      }
      """
    And I save the JSON response key "id" as "partialUpdateDeviceID"
    When I send a PUT request to "/api/users/devices/${partialUpdateDeviceID}" with body:
      """
      {
        "ip_address": "192.168.1.106"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "ip_address": "192.168.1.106"
    And the JSON response should contain "user_agent": "Original User Agent"
    And the JSON response should contain "device_type": "desktop"
    And the JSON response should contain "browser": "Original Browser"

  Scenario: TC14.2 - Update device with multiple partial fields
    Given I generate a unique alias as "multiPartialUpdateDeviceUserUsername"
    And I generate a unique alias as "multiPartialUpdateDeviceUserEmail"
    And I generate a unique alias as "multiPartialUpdateDeviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${multiPartialUpdateDeviceUserRoleName}",
        "description": "Role for multi partial device update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateDeviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${multiPartialUpdateDeviceUserUsername}",
        "firstName": "Multi",
        "lastName": "DeviceUser",
        "email": "${multiPartialUpdateDeviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Multi Device Tester",
        "roleId": ${multiPartialUpdateDeviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateDeviceUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${multiPartialUpdateDeviceUserID},
        "ip_address": "192.168.1.107",
        "user_agent": "Original User Agent",
        "device_type": "desktop",
        "browser": "Original Browser",
        "browser_version": "1.0.0",
        "os": "Original OS",
        "language": "en-US"
      }
      """
    And I save the JSON response key "id" as "multiPartialUpdateDeviceID"
    When I send a PUT request to "/api/users/devices/${multiPartialUpdateDeviceID}" with body:
      """
      {
        "ip_address": "192.168.1.108",
        "device_type": "mobile",
        "browser": "Updated Browser",
        "language": "es-ES"
      }
      """
    Then the response code should be 200
    And the JSON response should contain "ip_address": "192.168.1.108"
    And the JSON response should contain "device_type": "mobile"
    And the JSON response should contain "browser": "Updated Browser"
    And the JSON response should contain "language": "es-ES"
    And the JSON response should contain "user_agent": "Original User Agent"
    And the JSON response should contain "os": "Original OS"

  Scenario: TC14.3 - Update device with empty request
    Given I generate a unique alias as "emptyUpdateDeviceUserUsername"
    And I generate a unique alias as "emptyUpdateDeviceUserEmail"
    And I generate a unique alias as "emptyUpdateDeviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${emptyUpdateDeviceUserRoleName}",
        "description": "Role for empty device update test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "emptyUpdateDeviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${emptyUpdateDeviceUserUsername}",
        "firstName": "Empty",
        "lastName": "DeviceUser",
        "email": "${emptyUpdateDeviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Empty Device Tester",
        "roleId": ${emptyUpdateDeviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "emptyUpdateDeviceUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${emptyUpdateDeviceUserID},
        "ip_address": "192.168.1.109",
        "user_agent": "Empty Update Test",
        "device_type": "desktop",
        "browser": "Test Browser",
        "browser_version": "1.0.0",
        "os": "Test OS",
        "language": "en-US"
      }
      """
    And I save the JSON response key "id" as "emptyUpdateDeviceID"
    When I send a PUT request to "/api/users/devices/${emptyUpdateDeviceID}" with body:
      """
      {
      }
      """
    Then the response code should be 400
    And the JSON response should contain error "error": "No fields to update"

  Scenario: TC15 - Delete a device
    Given I generate a unique alias as "deleteDeviceUserUsername"
    And I generate a unique alias as "deleteDeviceUserEmail"
    And I generate a unique alias as "deleteDeviceUserRoleName"
    And I send a POST request to "/api/users/roles" with body:
      """
      {
        "name": "${deleteDeviceUserRoleName}",
        "description": "Role for device deletion test",
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deleteDeviceUserRoleID"
    And I send a POST request to "/api/users" with body:
      """
      {
        "username": "${deleteDeviceUserUsername}",
        "firstName": "Delete",
        "lastName": "DeviceUser",
        "email": "${deleteDeviceUserEmail}@test.com",
        "password": "securePassword123",
        "jobPosition": "Delete Device Tester",
        "roleId": ${deleteDeviceUserRoleID},
        "enabled": true
      }
      """
    And I save the JSON response key "id" as "deleteDeviceUserID"
    And I send a POST request to "/api/users/devices" with body:
      """
      {
        "userId": ${deleteDeviceUserID},
        "ip_address": "192.168.1.105",
        "user_agent": "Device to be deleted",
        "device_type": "desktop",
        "browser": "Edge",
        "browser_version": "120.0.0.0",
        "os": "Windows",
        "language": "en-US"
      }
      """
    And I save the JSON response key "id" as "deleteDeviceID"
    When I send a DELETE request to "/api/users/devices/${deleteDeviceID}"
    Then the response code should be 200
    And the JSON response should contain "message": "Device deleted successfully"

  Scenario: TC16 - Search devices paginated
    When I send a GET request to "/api/users/devices/search-paginated?page=1&limit=5"
    Then the response code should be 200
    And the JSON response should contain key "current_page"
    And the JSON response should contain key "records"
    And the JSON response should contain key "page_size"
    And the JSON response should contain key "total_pages"
    And the JSON response should contain key "total_records"

  Scenario: TC17 - Search device coincidences by property
    When I send a GET request to "/api/users/devices/search-by-property?property=device_type&search_text=desktop"
    Then the response code should be 200
    And the JSON response should be an array 