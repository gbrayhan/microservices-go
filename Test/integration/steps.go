//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/google/uuid"
)

var (
	base              = "http://localhost:8080"
	resp              *http.Response
	body              []byte
	err               error
	savedVars         = make(map[string]string)
	logger            = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
	createdResources  []string
	scenarioResources []string
	deletedResources  []string // Track resources that have been manually deleted
	lastRequestMethod string
	lastRequestPath   string
	skipNextTracking  bool // Flag to skip tracking for manually deleted resources
)

// ===== AUTONOMOUS RESOURCE MANAGEMENT FUNCTIONS =====

// Global variables to track autonomous resources
var (
	autonomousResources = make(map[string][]string) // resourceType -> []resourceIDs
	currentScenarioID   string
)

func TestMain(m *testing.M) {
	flag.Parse()
	if flag.Lookup("test.v") == nil || flag.Lookup("test.v").Value.String() != "true" {
		logger.SetOutput(os.Stderr)
	}
	os.Exit(m.Run())
}

// generateUniqueValue generates a unique value with a prefix
func generateUniqueValue(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, generateUUID())
}

// generateUniqueRFC generates a unique RFC
func generateUniqueRFC() string {
	return generateUniqueValue("RFC")
}

// generateUniqueAlias generates a unique alias
func generateUniqueAlias(prefix string) string {
	return generateUniqueValue(prefix)
}

// generateUniqueLegalName generates a unique legal name
func generateUniqueLegalName(prefix string) string {
	return generateUniqueValue(prefix)
}

// generateUUID genera un string UUID
func generateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
}

func theServiceIsInitialized() error {
	logger.Println("Service initialized.")
	return nil
}

func iSendARequestTo(method, path string) error {
	substitutedPath := substitute(path)
	fullURL := base + substitutedPath
	lastRequestMethod = method
	lastRequestPath = substitutedPath
	if strings.Contains(substitutedPath, "${") {
		logger.Printf("ERROR: Unsubstituted variable found in request URL: %s", substitutedPath)
		logger.Printf("Available variables: %v", savedVars)
		return fmt.Errorf("unsubstituted variable found in request URL: %s", substitutedPath)
	}
	logger.Printf("Sending %s request to: %s\n", method, fullURL)
	req, e := http.NewRequest(method, fullURL, nil)
	if e != nil {
		logger.Printf("Error creating request: %v\n", e)
		return e
	}
	addAuthHeader(req)
	logger.Printf("Request headers: %v\n", req.Header)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("Error sending request: %v\n", err)
	} else {
		logger.Printf("Received response status: %s\n", resp.Status)
		// If this is a successful DELETE request, mark the resource as manually deleted
		if method == "DELETE" && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			logger.Printf("MANUAL DELETE DETECTED: %s - Resource will be skipped in cleanup", fullURL)
			// Remove this resource from tracking since it was manually deleted
			removeResourceFromTracking(fullURL)
		}
	}
	return err
}

func iSendARequestWithBody(method, path string, payload *godog.DocString) error {
	substitutedPath := substitute(path)
	fullURL := base + substitutedPath
	lastRequestMethod = method
	lastRequestPath = substitutedPath
	bodyContent := replaceVars(payload.Content)
	if strings.Contains(substitutedPath, "${") || strings.Contains(bodyContent, "${") {
		logger.Printf("ERROR: Unsubstituted variable found in request. URL: %s, Body: %s", substitutedPath, bodyContent)
		logger.Printf("Available variables: %v", savedVars)
		return fmt.Errorf("unsubstituted variable found in request: URL or body contains ${...}")
	}
	logger.Printf("Sending %s request to: %s\n", method, fullURL)
	logger.Printf("Request body: %s\n", bodyContent)
	// Validate JSON before sending. If invalid, log but continue so tests can send malformed payloads
	var jsonTest interface{}
	if err := json.Unmarshal([]byte(bodyContent), &jsonTest); err != nil {
		logger.Printf("Warning: Invalid JSON in request body: %v", err)
		logger.Printf("Original body: %s", payload.Content)
		logger.Printf("Substituted body: %s", bodyContent)
	}
	req, e := http.NewRequest(method, fullURL, bytes.NewBufferString(bodyContent))
	if e != nil {
		logger.Printf("Error creating request with body: %v\n", e)
		return e
	}
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req)
	logger.Printf("Request headers: %v\n", req.Header)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("Error sending request with body: %v\n", err)
		return err
	}
	logger.Printf("Received response status: %s\n", resp.Status)

	// Read and save the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Error reading response body: %v\n", err)
		resp.Body.Close()
		return err
	}
	body = bodyBytes
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}

func theResponseCodeShouldBe(code int) error {
	if resp == nil {
		return fmt.Errorf("response is nil, cannot check status code")
	}
	logger.Printf("Validating response code. Expected: %d, Got: %d\n", code, resp.StatusCode)
	if resp.StatusCode != code {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			logger.Printf("Error closing response body: %v\n", err)
			return err
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		logger.Printf("Response body on error: %s\n", string(bodyBytes))
		return fmt.Errorf("expected status code %d but got %d. Body: %s", code, resp.StatusCode, string(bodyBytes))
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Error reading response body: %v\n", err)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

func getNestedValue(m map[string]interface{}, path string) (interface{}, bool) {
	parts := strings.Split(path, ".")
	var current interface{} = m
	for _, p := range parts {
		if mp, ok := current.(map[string]interface{}); ok {
			if val, ok := mp[p]; ok {
				current = val
			} else {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	return current, true
}

func theJSONResponseShouldContainKey(key string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if _, exists := getNestedValue(response, key); !exists {
		return fmt.Errorf("expected key '%s' not found in response", key)
	}

	return nil
}

func theJSONResponseShouldContain(field, value string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	// Substitute any saved or environment variables in the expected value
	value = replaceVars(value)

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := getNestedValue(response, field)
	if !exists {
		return fmt.Errorf("field '%s' not found in response", field)
	}

	// Convert field value to string for comparison
	var fieldStr string
	switch v := fieldValue.(type) {
	case string:
		fieldStr = v
	case float64:
		fieldStr = fmt.Sprintf("%.0f", v)
	case int:
		fieldStr = fmt.Sprintf("%d", v)
	case bool:
		fieldStr = fmt.Sprintf("%t", v)
	default:
		fieldStr = fmt.Sprintf("%v", v)
	}

	if fieldStr != value {
		return fmt.Errorf("expected field '%s' to be '%s' but got '%s'", field, value, fieldStr)
	}

	return nil
}

func theJSONResponseShouldContainError(field, expectedError string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Check if the response contains an error field
	if errors, exists := response["errors"]; exists {
		if errorsArray, ok := errors.([]interface{}); ok {
			for _, err := range errorsArray {
				if errMap, ok := err.(map[string]interface{}); ok {
					if fieldValue, exists := errMap[field]; exists {
						if fieldStr, ok := fieldValue.(string); ok && fieldStr == expectedError {
							return nil
						}
					}
				}
			}
		}
	}

	// Also check direct field access
	if fieldValue, exists := response[field]; exists {
		if fieldStr, ok := fieldValue.(string); ok && fieldStr == expectedError {
			return nil
		}
	}

	return fmt.Errorf("expected error '%s' for field '%s' not found in response", expectedError, field)
}

func theJSONResponseShouldContainErrorMessage(expectedError string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Check for error message in various possible locations
	possibleFields := []string{"error", "message", "errorMessage", "msg"}
	for _, field := range possibleFields {
		if fieldValue, exists := response[field]; exists {
			if fieldStr, ok := fieldValue.(string); ok && strings.Contains(fieldStr, expectedError) {
				return nil
			}
		}
	}

	// Check in errors array
	if errors, exists := response["errors"]; exists {
		if errorsArray, ok := errors.([]interface{}); ok {
			for _, err := range errorsArray {
				if errMap, ok := err.(map[string]interface{}); ok {
					for _, field := range possibleFields {
						if fieldValue, exists := errMap[field]; exists {
							if fieldStr, ok := fieldValue.(string); ok && strings.Contains(fieldStr, expectedError) {
								return nil
							}
						}
					}
				}
			}
		}
	}

	return fmt.Errorf("expected error message containing '%s' not found in response", expectedError)
}

func theJSONResponseShouldContainNumeric(field string, expectedValue int) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := getNestedValue(response, field)
	if !exists {
		return fmt.Errorf("field '%s' not found in response", field)
	}

	var actualValue int
	switch v := fieldValue.(type) {
	case float64:
		actualValue = int(v)
	case int:
		actualValue = v
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			actualValue = parsed
		} else {
			return fmt.Errorf("field '%s' is a string that cannot be converted to int: %s", field, v)
		}
	default:
		return fmt.Errorf("field '%s' has unexpected type: %T", field, fieldValue)
	}

	if actualValue != expectedValue {
		return fmt.Errorf("expected field '%s' to be %d but got %d", field, expectedValue, actualValue)
	}

	return nil
}

func theJSONResponseShouldContainWithNumericValueUserID(field string) error {
	idStr, ok := savedVars["userID"]
	if !ok {
		return fmt.Errorf("userID not found in saved variables")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("userID is not numeric: %v", err)
	}
	return theJSONResponseShouldContainNumeric(field, id)
}

func iSaveTheJSONResponseKeyAs(key, varName string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := getNestedValue(response, key)
	if !exists {
		return fmt.Errorf("key '%s' not found in response", key)
	}

	// Convert the value to string
	var valueStr string
	switch v := fieldValue.(type) {
	case string:
		valueStr = v
	case float64:
		valueStr = fmt.Sprintf("%.0f", v)
	case int:
		valueStr = fmt.Sprintf("%d", v)
	case bool:
		valueStr = fmt.Sprintf("%t", v)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	savedVars[varName] = valueStr
	logger.Printf("Saved '%s' as '%s' with value: %s", key, varName, valueStr)

	// Track the resource if it's an ID field
	if strings.HasSuffix(key, "Id") || strings.HasSuffix(key, "ID") || key == "id" {
		// Determine the resource type based on the context
		resourceType := "unknown"
		// Match more specific paths first to ensure correct resource type
		if strings.Contains(lastRequestPath, "/medicine") {
			resourceType = "medicine"
		} else if strings.Contains(lastRequestPath, "/user") {
			resourceType = "user"
		}

		if resourceType != "unknown" {
			trackAutonomousResource(resourceType, valueStr)
		}
	}

	return nil
}

func substitute(path string) string {
	result := path

	// First substitute saved variables
	for key, value := range savedVars {
		placeholder := "${" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Then substitute environment variables
	for _, envVar := range os.Environ() {
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			value := pair[1]
			placeholder := "${" + key + "}"
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	return result
}

func replaceVars(s string) string {
	result := s

	// First substitute saved variables
	for key, value := range savedVars {
		placeholder := "${" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}

	// Then substitute environment variables
	for _, envVar := range os.Environ() {
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) == 2 {
			key := pair[0]
			value := pair[1]
			placeholder := "${" + key + "}"
			result = strings.ReplaceAll(result, placeholder, value)
		}
	}

	return result
}

func addAuthHeader(req *http.Request) {
	if token, exists := savedVars["accessToken"]; exists && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		logger.Printf("DEBUG: No accessToken found in savedVars. exists: %t, token: %v", exists, token)
		// Check if authentication failed during setup
		if authFailed, exists := savedVars["auth_failed"]; exists && authFailed == "true" {
			logger.Printf("WARNING: Authentication failed during test setup. This request will likely fail with 401.")
		}
	}
}

func trackResource(path string) {
	if !skipNextTracking {
		createdResources = append(createdResources, path)
		logger.Printf("Tracking resource: %s", path)
	}
}

func removeResourceFromTracking(fullURL string) {
	// Extract the path from the full URL
	path := strings.TrimPrefix(fullURL, base)

	// Remove from createdResources
	for i, resource := range createdResources {
		if resource == path {
			createdResources = append(createdResources[:i], createdResources[i+1:]...)
			logger.Printf("Removed from tracking: %s", path)
			break
		}
	}
}

func iGenerateAUniqueEANCodeAs(varName string) error {
	uniqueEAN := generateUniqueValue("EAN")
	savedVars[varName] = uniqueEAN
	logger.Printf("Generated unique EAN code: %s", uniqueEAN)
	return nil
}

func iGenerateAUniqueLoteAs(varName string) error {
	uniqueLote := generateUniqueValue("LOTE")
	savedVars[varName] = uniqueLote
	logger.Printf("Generated unique lote: %s", uniqueLote)
	return nil
}

func iGenerateAUniqueRFCAs(varName string) error {
	uniqueRFC := generateUniqueRFC()
	savedVars[varName] = uniqueRFC
	logger.Printf("Generated unique RFC: %s", uniqueRFC)
	return nil
}

func iSaveFirstArrayElementKeyAs(key, arrayKey, varName string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	arrayValue, exists := response[arrayKey]
	if !exists {
		return fmt.Errorf("array key '%s' not found in response", arrayKey)
	}

	array, ok := arrayValue.([]interface{})
	if !ok {
		return fmt.Errorf("field '%s' is not an array", arrayKey)
	}

	if len(array) == 0 {
		return fmt.Errorf("array '%s' is empty", arrayKey)
	}

	firstElement, ok := array[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("first element of array '%s' is not an object", arrayKey)
	}

	fieldValue, exists := firstElement[key]
	if !exists {
		return fmt.Errorf("key '%s' not found in first element of array '%s'", key, arrayKey)
	}

	// Convert the value to string
	var valueStr string
	switch v := fieldValue.(type) {
	case string:
		valueStr = v
	case float64:
		valueStr = fmt.Sprintf("%.0f", v)
	case int:
		valueStr = fmt.Sprintf("%d", v)
	case bool:
		valueStr = fmt.Sprintf("%t", v)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	savedVars[varName] = valueStr
	logger.Printf("Saved first element key '%s' from array '%s' as '%s' with value: %s", key, arrayKey, varName, valueStr)

	return nil
}

func iSendAPatchRequestTo(path string) error {
	return iSendARequestTo("PATCH", path)
}

func theJSONResponseFieldShouldContainString(fieldName, expectedSubstring string) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := getNestedValue(response, fieldName)
	if !exists {
		return fmt.Errorf("field '%s' not found in response", fieldName)
	}

	fieldStr, ok := fieldValue.(string)
	if !ok {
		return fmt.Errorf("field '%s' is not a string", fieldName)
	}

	if !strings.Contains(fieldStr, expectedSubstring) {
		return fmt.Errorf("field '%s' does not contain substring '%s'. Actual value: '%s'", fieldName, expectedSubstring, fieldStr)
	}

	return nil
}

func theJSONResponseShouldContainBoolean(fieldName string, expectedValue bool) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := response[fieldName]
	if !exists {
		return fmt.Errorf("field '%s' not found in response", fieldName)
	}

	fieldBool, ok := fieldValue.(bool)
	if !ok {
		return fmt.Errorf("field '%s' is not a boolean", fieldName)
	}

	if fieldBool != expectedValue {
		return fmt.Errorf("expected field '%s' to be %t but got %t", fieldName, expectedValue, fieldBool)
	}

	return nil
}

func theResponseCodeShouldBeOr(code1, code2 int) error {
	if resp == nil {
		return fmt.Errorf("response is nil, cannot check status code")
	}
	logger.Printf("Validating response code. Expected: %d or %d, Got: %d\n", code1, code2, resp.StatusCode)
	if resp.StatusCode != code1 && resp.StatusCode != code2 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			logger.Printf("Error closing response body: %v\n", err)
			return err
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		logger.Printf("Response body on error: %s\n", string(bodyBytes))
		return fmt.Errorf("expected status code %d or %d but got %d. Body: %s", code1, code2, resp.StatusCode, string(bodyBytes))
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Error reading response body: %v\n", err)
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

func theJSONResponseShouldBeAnArray() error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response []interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("response is not a valid JSON array: %v", err)
	}

	return nil
}

func theJSONResponseShouldBeAnObject() error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("response is not a valid JSON object: %v", err)
	}

	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		logger.Println("Setting up test suite...")

		// Get the initial user credentials from environment variables
		startUserEmail := os.Getenv("START_USER_EMAIL")
		startUserPw := os.Getenv("START_USER_PW")

		if startUserEmail == "" || startUserPw == "" {
			logger.Printf("Warning: START_USER_EMAIL or START_USER_PW not set, using default test credentials")
			startUserEmail = "test@test.com"
			startUserPw = "test123"
		}

		// Try to authenticate with the initial user credentials
		loginData := map[string]interface{}{
			"email":    startUserEmail,
			"password": startUserPw,
		}

		jsonData, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", base+"/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Printf("Failed to authenticate: %v", err)
			// If authentication fails, try to create a test user
			createTestUserDirectly()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			var loginResponse map[string]interface{}
			if json.Unmarshal(bodyBytes, &loginResponse) == nil {
				if sec, ok := loginResponse["security"].(map[string]interface{}); ok {
					if accessToken, ok := sec["jwtAccessToken"].(string); ok {
						savedVars["accessToken"] = accessToken
						logger.Printf("Authentication successful with %s, access token saved", startUserEmail)
						return
					}
				}
			}
		} else {
			logger.Printf("Authentication with %s failed with status: %d", startUserEmail, resp.StatusCode)
			// Try to create a test user if authentication fails
			createTestUserDirectly()
		}
	})

	ctx.AfterSuite(func() {
		logger.Println("Cleaning up test suite...")
		// Clean up any remaining resources
		for _, resource := range createdResources {
			logger.Printf("Cleaning up resource: %s", resource)
			req, _ := http.NewRequest("DELETE", base+resource, nil)
			addAuthHeader(req)
			http.DefaultClient.Do(req)
		}

		for resourceType, resourceIDs := range autonomousResources {
			for _, resourceID := range resourceIDs {
				logger.Printf("Cleaning up autonomous resource: %s/%s", resourceType, resourceID)
				deleteAutonomousResource(resourceType, resourceID)
			}
		}
	})
}

// createTestUserDirectly creates a test user by making direct API calls
// This function assumes that the seeding mechanism has already created the admin role
func createTestUserDirectly() {
	logger.Println("Attempting to create test user directly...")

	// First, try to authenticate with the seeded user credentials
	startUserEmail := os.Getenv("START_USER_EMAIL")
	startUserPw := os.Getenv("START_USER_PW")

	if startUserEmail == "" || startUserPw == "" {
		logger.Println("START_USER_EMAIL or START_USER_PW not set, using default test credentials")
		startUserEmail = "test@test.com"
		startUserPw = "test123"
	}

	logger.Printf("🔍 Attempting authentication with email: %s", startUserEmail)

	loginData := map[string]interface{}{
		"email":    startUserEmail,
		"password": startUserPw,
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", base+"/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("❌ Network error during authentication: %v", err)
		savedVars["auth_failed"] = "true"
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	logger.Printf("🔍 Authentication response status: %d", resp.StatusCode)
	logger.Printf("🔍 Authentication response body: %s", string(bodyBytes))

	if resp.StatusCode == 200 {
		var loginResponse map[string]interface{}
		if json.Unmarshal(bodyBytes, &loginResponse) == nil {
			if sec, ok := loginResponse["security"].(map[string]interface{}); ok {
				if accessToken, ok := sec["jwtAccessToken"].(string); ok {
					savedVars["accessToken"] = accessToken
					logger.Printf("✅ Authentication successful with seeded user, access token saved: %s", accessToken)
					return
				}
			}
			logger.Printf("❌ No accessToken found in successful response: %v", loginResponse)
		} else {
			logger.Printf("❌ Failed to parse authentication response: %v", err)
		}
	} else {
		logger.Printf("❌ Authentication failed with status: %d", resp.StatusCode)
	}

	// If authentication fails, we cannot proceed with tests
	logger.Println("❌ Authentication failed with seeded user credentials")
	logger.Printf("   Email: %s", startUserEmail)
	logger.Println("   This usually means:")
	logger.Println("   1. The application is not properly seeded")
	logger.Println("   2. The database is not initialized correctly")
	logger.Println("   3. The START_USER_EMAIL and START_USER_PW environment variables are incorrect")
	logger.Println("   4. The application failed to start or migrate")
	logger.Println("")
	logger.Println("   Please ensure:")
	logger.Println("   1. The database is running and accessible")
	logger.Println("   2. All environment variables are set correctly")
	logger.Println("   3. The application starts without errors")
	logger.Println("   4. The seeding process completes successfully")
	logger.Println("")
	logger.Println("   You can check the application logs for more details.")

	// Set a flag to indicate authentication failure
	savedVars["auth_failed"] = "true"
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		logger.Printf("Starting scenario: %s", sc.Name)
		currentScenarioID = generateUUID()
		scenarioResources = []string{}
		skipNextTracking = false

		// Ensure we have a valid authentication token for each scenario
		if token, exists := savedVars["accessToken"]; !exists || token == "" {
			createTestUserDirectly()
		}

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		logger.Printf("Ending scenario: %s", sc.Name)

		// Clean up scenario-specific resources
		for _, resource := range scenarioResources {
			logger.Printf("Cleaning up scenario resource: %s", resource)
			req, _ := http.NewRequest("DELETE", base+resource, nil)
			addAuthHeader(req)
			http.DefaultClient.Do(req)
		}

		// Clear scenario-specific variables
		for key := range savedVars {
			if strings.HasPrefix(key, "scenario_") {
				delete(savedVars, key)
			}
		}

		return ctx, nil
	})

	// Basic HTTP request steps
	ctx.Step(`^the service is initialized$`, theServiceIsInitialized)
	ctx.Step(`^I send a (GET|POST|PUT|DELETE|PATCH) request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^I send a (GET|POST|PUT|DELETE|PATCH) request to "([^"]*)" with body:$`, iSendARequestWithBody)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response code should be (\d+) or (\d+)$`, theResponseCodeShouldBeOr)

	// JSON response validation steps
	ctx.Step(`^the JSON response should contain key "([^"]*)"$`, theJSONResponseShouldContainKey)
	ctx.Step(`^the JSON response should contain "([^"]*)" with value "([^"]*)"$`, theJSONResponseShouldContain)
	ctx.Step(`^the JSON response should contain error "([^"]*)" for field "([^"]*)"$`, theJSONResponseShouldContainError)
	ctx.Step(`^the JSON response should contain error message "([^"]*)"$`, theJSONResponseShouldContainErrorMessage)
	ctx.Step(`^the JSON response should contain "([^"]*)" with numeric value (\d+)$`, theJSONResponseShouldContainNumeric)
	ctx.Step(`^the JSON response should contain "([^"]*)" with numeric value  \$\{userID\}$`, theJSONResponseShouldContainWithNumericValueUserID)
	ctx.Step(`^the JSON response should contain "([^"]*)" with boolean value (true|false)$`, theJSONResponseShouldContainBoolean)
	ctx.Step(`^the JSON response should be an array$`, theJSONResponseShouldBeAnArray)
	ctx.Step(`^the JSON response field "([^"]*)" should contain string "([^"]*)"$`, theJSONResponseFieldShouldContainString)

	// Variable management steps
	ctx.Step(`^I save the JSON response key "([^"]*)" as "([^"]*)"$`, iSaveTheJSONResponseKeyAs)
	ctx.Step(`^I save the first array element key "([^"]*)" from array "([^"]*)" as "([^"]*)"$`, iSaveFirstArrayElementKeyAs)

	// Unique value generation steps
	ctx.Step(`^I generate a unique EAN code as "([^"]*)"$`, iGenerateAUniqueEANCodeAs)
	ctx.Step(`^I generate a unique lote as "([^"]*)"$`, iGenerateAUniqueLoteAs)
	ctx.Step(`^I generate a unique RFC as "([^"]*)"$`, iGenerateAUniqueRFCAs)
	ctx.Step(`^I generate a unique alias as "([^"]*)"$`, iGenerateAUniqueAliasAs)
	ctx.Step(`^I generate a unique legal name as "([^"]*)"$`, iGenerateAUniqueLegalNameAs)

	// Additional JSON response validation steps
	ctx.Step(`^the JSON response should contain "([^"]*)": "([^"]*)"$`, theJSONResponseShouldContain)
	ctx.Step(`^the JSON response should contain "([^"]*)": (\d+)$`, theJSONResponseShouldContainNumeric)
	ctx.Step(`^the JSON response should contain "([^"]*)": (\d+)\.(\d+)$`, theJSONResponseShouldContainFloatWithTwoParams)
	ctx.Step(`^the JSON response should contain "([^"]*)": true$`, theJSONResponseShouldContainTrue)
	ctx.Step(`^the JSON response should contain "([^"]*)": false$`, theJSONResponseShouldContainFalse)
	ctx.Step(`^the JSON response should contain error "([^"]*)": "([^"]*)"$`, theJSONResponseShouldContainError)
	ctx.Step(`^the JSON response should be an object$`, theJSONResponseShouldBeAnObject)

	// Authentication steps
	ctx.Step(`^I clear the authentication token$`, iClearTheAuthenticationToken)
}

func iClearTheAuthenticationToken() error {
	delete(savedVars, "accessToken")
	logger.Println("Authentication token cleared")
	return nil
}

func theJSONResponseShouldContainTrue(field string) error {
	return theJSONResponseShouldContainBoolean(field, true)
}

func theJSONResponseShouldContainFalse(field string) error {
	return theJSONResponseShouldContainBoolean(field, false)
}

func theJSONResponseShouldContainFloat(field string, expectedValue float64) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := getNestedValue(response, field)
	if !exists {
		return fmt.Errorf("field '%s' not found in response", field)
	}

	var actualValue float64
	switch v := fieldValue.(type) {
	case float64:
		actualValue = v
	case int:
		actualValue = float64(v)
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			actualValue = parsed
		} else {
			return fmt.Errorf("field '%s' is a string that cannot be converted to float: %s", field, v)
		}
	default:
		return fmt.Errorf("field '%s' has unexpected type: %T", field, fieldValue)
	}

	if actualValue != expectedValue {
		return fmt.Errorf("expected field '%s' to be %f but got %f", field, expectedValue, actualValue)
	}

	return nil
}

func theJSONResponseShouldContainFloatWithTwoParams(field string, wholePart, decimalPart int) error {
	expectedValue := float64(wholePart) + float64(decimalPart)/100.0
	return theJSONResponseShouldContainFloat(field, expectedValue)
}

func iGenerateAUniqueAliasAs(varName string) error {
	uniqueAlias := generateUniqueAlias("ALIAS")
	savedVars[varName] = uniqueAlias
	logger.Printf("Generated unique alias: %s", uniqueAlias)
	return nil
}

func iGenerateAUniqueLegalNameAs(varName string) error {
	uniqueLegalName := generateUniqueLegalName("LEGAL")
	savedVars[varName] = uniqueLegalName
	logger.Printf("Generated unique legal name: %s", uniqueLegalName)
	return nil
}

func createUniqueMedicine(prefix string) (int, error) {
	medicineData := map[string]interface{}{
		"name":        fmt.Sprintf("%s Medicine", prefix),
		"description": fmt.Sprintf("Description for %s medicine", prefix),
		"eanCode":     generateUniqueValue("EAN"),
		"laboratory":  "TestLab",
	}

	jsonData, _ := json.Marshal(medicineData)
	req, _ := http.NewRequest("POST", base+"/v1/medicine", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		var medicine map[string]interface{}
		if json.Unmarshal(bodyBytes, &medicine) == nil {
			if id, ok := medicine["id"]; ok {
				if idFloat, ok := id.(float64); ok {
					medicineID := int(idFloat)
					trackAutonomousResource("medicine", fmt.Sprintf("%d", medicineID))
					return medicineID, nil
				}
			}
		}
	}

	return 0, fmt.Errorf("failed to create medicine: %d", resp.StatusCode)
}

func iEnsureATestMedicineExists() error {
	medicineID, err := createUniqueMedicine("TestMedicine")
	if err != nil {
		return err
	}
	savedVars["testMedicineID"] = fmt.Sprintf("%d", medicineID)
	logger.Printf("Ensured test medicine exists with ID: %d", medicineID)
	return nil
}

func iCreateATestMedicine() error {
	return iEnsureATestMedicineExists()
}

func theJSONResponseFieldShouldContainArrayWithElement(field string, expectedCount int) error {
	if body == nil {
		return fmt.Errorf("response body is nil")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fieldValue, exists := response[field]
	if !exists {
		return fmt.Errorf("field '%s' not found in response", field)
	}

	array, ok := fieldValue.([]interface{})
	if !ok {
		return fmt.Errorf("field '%s' is not an array", field)
	}

	if len(array) != expectedCount {
		return fmt.Errorf("expected array '%s' to have %d elements but got %d", field, expectedCount, len(array))
	}

	return nil
}

func theJSONResponseShouldContainWithValue(field, expectedValue string) error {
	return theJSONResponseShouldContain(field, expectedValue)
}

func iGetAValidICDCIECodeForTheClient() error {
	// Create a test ICD-CIE code
	icdData := map[string]interface{}{
		"code":        generateUniqueValue("ICD"),
		"description": "Test ICD-CIE code for integration tests",
		"active":      true,
	}

	jsonData, _ := json.Marshal(icdData)
	req, _ := http.NewRequest("POST", base+"/api/icd-cie", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		var icd map[string]interface{}
		if json.Unmarshal(bodyBytes, &icd) == nil {
			if id, ok := icd["id"]; ok {
				if idFloat, ok := id.(float64); ok {
					icdID := int(idFloat)
					savedVars["testICDID"] = fmt.Sprintf("%d", icdID)
					trackAutonomousResource("icd-cie", fmt.Sprintf("%d", icdID))
					logger.Printf("Created test ICD-CIE with ID: %d", icdID)
					return nil
				}
			}
		}
	}

	return fmt.Errorf("failed to create ICD-CIE code: %d", resp.StatusCode)
}

func iGetAValidZipCodeAndNeighborhood() error {
	// For testing purposes, we'll use a static zip code and neighborhood
	savedVars["testZipCode"] = "12345"
	savedVars["testNeighborhood"] = "Test Neighborhood"
	logger.Printf("Set test zip code: %s and neighborhood: %s", savedVars["testZipCode"], savedVars["testNeighborhood"])
	return nil
}

func trackAutonomousResource(resourceType, resourceID string) {
	autonomousResources[resourceType] = append(autonomousResources[resourceType], resourceID)
	logger.Printf("Tracking autonomous resource: %s/%s", resourceType, resourceID)
}

func createTestUser() {
	logger.Println("Attempting to create test user...")

	// First create a role
	roleData := map[string]interface{}{
		"name":        "test_role",
		"description": "Test role for integration tests",
		"enabled":     true,
	}

	jsonData, _ := json.Marshal(roleData)
	req, _ := http.NewRequest("POST", base+"/api/users/roles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("Failed to create test role: %v", err)
		return
	}
	defer resp.Body.Close()

	var roleID int
	if resp.StatusCode == 201 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		var roleResponse map[string]interface{}
		if json.Unmarshal(bodyBytes, &roleResponse) == nil {
			if id, ok := roleResponse["id"]; ok {
				if idFloat, ok := id.(float64); ok {
					roleID = int(idFloat)
					logger.Printf("Created test role with ID: %d", roleID)
				}
			}
		}
	} else {
		logger.Printf("Failed to create test role, status: %d", resp.StatusCode)
		return
	}

	// Get credentials from environment variables or use defaults
	startUserEmail := os.Getenv("START_USER_EMAIL")
	startUserPw := os.Getenv("START_USER_PW")

	if startUserEmail == "" {
		startUserEmail = "test@test.com"
	}
	if startUserPw == "" {
		startUserPw = "test123"
	}

	// Now create a test user
	userData := map[string]interface{}{
		"username":    startUserEmail,
		"firstName":   "Test",
		"lastName":    "User",
		"email":       startUserEmail,
		"password":    startUserPw,
		"jobPosition": "Tester",
		"roleId":      roleID,
		"enabled":     true,
	}

	jsonData, _ = json.Marshal(userData)
	req, _ = http.NewRequest("POST", base+"/api/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("Failed to create test user: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		logger.Printf("Created test user successfully")
		// Try to authenticate with the new user
		loginData := map[string]interface{}{
			"email":    startUserEmail,
			"password": startUserPw,
		}

		jsonData, _ = json.Marshal(loginData)
		req, _ = http.NewRequest("POST", base+"/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			logger.Printf("Failed to authenticate with new test user: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			var loginResponse map[string]interface{}
			if json.Unmarshal(bodyBytes, &loginResponse) == nil {
				if sec, ok := loginResponse["security"].(map[string]interface{}); ok {
					if accessToken, ok := sec["jwtAccessToken"].(string); ok {
						savedVars["accessToken"] = accessToken
						logger.Printf("Authentication successful with new test user, access token saved")
					}
				}
			}
		}
	} else {
		logger.Printf("Failed to create test user, status: %d", resp.StatusCode)
	}
}

func deleteAutonomousResource(resourceType, resourceID string) error {
	var endpoint string
	switch resourceType {
	case "user":
		endpoint = fmt.Sprintf("/v1/user/%s", resourceID)
	case "medicine":
		endpoint = fmt.Sprintf("/v1/medicine/%s", resourceID)
	default:
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}

	req, _ := http.NewRequest("DELETE", base+endpoint, nil)
	addAuthHeader(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Printf("Successfully deleted autonomous resource: %s/%s", resourceType, resourceID)
		return nil
	}

	return fmt.Errorf("failed to delete autonomous resource %s/%s: %d", resourceType, resourceID, resp.StatusCode)
}
