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
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

var (
	base      = "http://localhost:8080"
	resp      *http.Response
	body      []byte
	err       error
	savedVars = make(map[string]string)
	logger    = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
)

func TestMain(m *testing.M) {
	flag.Parse()
	if flag.Lookup("test.v") == nil || flag.Lookup("test.v").Value.String() != "true" {
		logger.SetOutput(os.Stderr)
	}
	os.Exit(m.Run())
}

func init() {
	rand.New(rand.NewSource(12345))
}

// Generic service initialization step
func theServiceIsInitialized() error {
	logger.Println("Service initialized.")
	return nil
}

// Generic HTTP request step - supports any method and path
func iSendARequestTo(method, path string) error {
	fullURL := substitute(path)
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
	}
	return err
}

// Generic HTTP request with body step - supports any method, path and JSON body
func iSendARequestWithBody(method, path string, payload *godog.DocString) error {
	fullURL := substitute(path)
	bodyContent := replaceVars(payload.Content)
	logger.Printf("Sending %s request to: %s\n", method, fullURL)
	logger.Printf("Request body: %s\n", bodyContent)
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
	} else {
		logger.Printf("Received response status: %s\n", resp.Status)
	}
	return err
}

// Generic response status code validation
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
	logger.Printf("Response body: %s\n", string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

// Generic JSON key existence validation - supports nested keys with dot notation
func theJSONResponseShouldContainKey(key string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check for key %q", key)
	}
	logger.Printf("Checking if JSON response contains key: %q\n", key)

	// Handle nested keys with dot notation
	if strings.Contains(key, ".") {
		return validateNestedKey(key)
	}

	var obj map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &obj); errUnmarshal == nil {
		if _, ok := obj[key]; ok {
			return nil
		}
		return fmt.Errorf("expected JSON object to contain key %q", key)
	}

	var arr []interface{}
	if errUnmarshal := json.Unmarshal(body, &arr); errUnmarshal == nil {
		idx, errConv := strconv.Atoi(key)
		if errConv != nil {
			return fmt.Errorf("invalid array index %q", key)
		}
		if idx < 0 || idx >= len(arr) {
			return fmt.Errorf("array index %d out of range (len=%d)", idx, len(arr))
		}
		return nil
	}

	return fmt.Errorf("response is neither JSON object nor array; cannot check key %q", key)
}

// Generic JSON field value validation - supports nested fields, wildcards, and pattern matching
func theJSONResponseShouldContain(field, value string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check field %q", field)
	}
	expectedValue := replaceVars(value)
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}

	// Get the actual value using nested field support
	actualValue, err := getNestedValue(data, field)
	if err != nil {
		return err
	}

	// Handle different validation types
	switch {
	case expectedValue == "*":
		// Wildcard validation - just check that field exists and is not empty
		if actualValue == nil || actualValue == "" {
			return fmt.Errorf("expected field %q to be non-empty, but got nil or empty", field)
		}
		return nil
	case strings.HasPrefix(expectedValue, "regex:") && strings.HasSuffix(expectedValue, ":"):
		// Regex pattern validation
		pattern := strings.TrimPrefix(strings.TrimSuffix(expectedValue, ":"), "regex:")
		regex, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern %q: %v", pattern, err)
		}
		actualStr := fmt.Sprintf("%v", actualValue)
		if !regex.MatchString(actualStr) {
			return fmt.Errorf("field %q value %q does not match regex pattern %q", field, actualStr, pattern)
		}
		return nil
	case expectedValue == "null":
		// Null validation
		if actualValue != nil {
			return fmt.Errorf("expected field %q to be null, but got %v", field, actualValue)
		}
		return nil
	case expectedValue == "not_null":
		// Not null validation
		if actualValue == nil {
			return fmt.Errorf("expected field %q to be not null", field)
		}
		return nil
	default:
		// Exact value validation
		actualStr := fmt.Sprintf("%v", actualValue)
		if actualStr != expectedValue {
			return fmt.Errorf("expected %q = %q, but got %v", field, expectedValue, actualStr)
		}
		return nil
	}
}

// Generic JSON field type validation
func theJSONResponseFieldShouldBeOfType(fieldName, expectedType string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check field type %q", fieldName)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}

	actualValue, err := getNestedValue(data, fieldName)
	if err != nil {
		return err
	}

	actualType := getTypeName(actualValue)
	if actualType != expectedType {
		return fmt.Errorf("expected field %q to be of type %q, but got %q", fieldName, expectedType, actualType)
	}
	return nil
}

// Generic JSON array validation
func theJSONResponseShouldBeAnArray() error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check if it's an array")
	}
	var arr []interface{}
	if errUnmarshal := json.Unmarshal(body, &arr); errUnmarshal != nil {
		return fmt.Errorf("expected JSON array, but got: %v. Body: %s", errUnmarshal, string(body))
	}
	return nil
}

// Generic JSON array length validation
func theJSONResponseArrayShouldHaveLength(length int) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check array length")
	}
	var arr []interface{}
	if errUnmarshal := json.Unmarshal(body, &arr); errUnmarshal != nil {
		return fmt.Errorf("expected JSON array, but got: %v. Body: %s", errUnmarshal, string(body))
	}
	if len(arr) != length {
		return fmt.Errorf("expected array length %d, but got %d", length, len(arr))
	}
	return nil
}

// Generic variable saving from JSON response
func iSaveTheJSONResponseKeyAs(key, varName string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot save key %q", key)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON for saving: %v. Body: %s", errUnmarshal, string(body))
	}

	actualValue, err := getNestedValue(data, key)
	if err != nil {
		return err
	}

	var valStr string
	switch v := actualValue.(type) {
	case string:
		valStr = v
	case float64:
		if v == float64(int64(v)) {
			valStr = strconv.FormatInt(int64(v), 10)
		} else {
			valStr = strconv.FormatFloat(v, 'f', -1, 64)
		}
	case int:
		valStr = strconv.Itoa(v)
	case bool:
		valStr = strconv.FormatBool(v)
	case nil:
		valStr = ""
	default:
		valStr = fmt.Sprintf("%v", v)
	}
	savedVars[varName] = valStr
	logger.Printf("Saved %q as %q\n", key, varName)
	return nil
}

// Generic array element saving
func iSaveFirstArrayElementKeyAs(key, arrayKey, varName string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot save array element key %q", key)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON for saving array element: %v. Body: %s", errUnmarshal, string(body))
	}
	array, ok := data[arrayKey]
	if !ok {
		return fmt.Errorf("array key %q not found in response", arrayKey)
	}
	arrayData, ok := array.([]interface{})
	if !ok {
		return fmt.Errorf("key %q is not an array", arrayKey)
	}
	if len(arrayData) == 0 {
		return fmt.Errorf("array %q is empty", arrayKey)
	}
	firstElement, ok := arrayData[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("first element of array %q is not an object", arrayKey)
	}
	raw, ok := firstElement[key]
	if !ok {
		return fmt.Errorf("key %q not found in first element of array %q", key, arrayKey)
	}
	var valStr string
	switch v := raw.(type) {
	case string:
		valStr = v
	case float64:
		if v == float64(int64(v)) {
			valStr = strconv.FormatInt(int64(v), 10)
		} else {
			valStr = strconv.FormatFloat(v, 'f', -1, 64)
		}
	case int:
		valStr = strconv.Itoa(v)
	case bool:
		valStr = strconv.FormatBool(v)
	case nil:
		valStr = ""
	default:
		valStr = fmt.Sprintf("%v", v)
	}
	savedVars[varName] = valStr
	logger.Printf("Saved first element key %q from array %q as %q\n", key, arrayKey, varName)
	return nil
}

// Generic PATCH request step
func iSendAPatchRequestTo(path string) error {
	return iSendARequestTo("PATCH", path)
}

// Generic substring validation
func theJSONResponseFieldShouldContainString(fieldName, expectedSubstring string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check field substring %q", fieldName)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}

	actualValue, err := getNestedValue(data, fieldName)
	if err != nil {
		return err
	}

	actualStr, ok := actualValue.(string)
	if !ok {
		return fmt.Errorf("field %q is not a string, got %T", fieldName, actualValue)
	}
	if !strings.Contains(actualStr, expectedSubstring) {
		return fmt.Errorf("field %q value %q does not contain substring %q", fieldName, actualStr, expectedSubstring)
	}
	return nil
}

// Generic boolean validation
func theJSONResponseShouldContainBoolean(fieldName string, expectedValue bool) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check boolean field %q", fieldName)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}

	actualValue, err := getNestedValue(data, fieldName)
	if err != nil {
		return err
	}

	actualBool, ok := actualValue.(bool)
	if !ok {
		return fmt.Errorf("field %q is not a boolean, got %T", fieldName, actualValue)
	}
	if actualBool != expectedValue {
		return fmt.Errorf("expected field %q to be %t, but got %t", fieldName, expectedValue, actualBool)
	}
	return nil
}

// Generic multiple status code validation
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
	logger.Printf("Response body: %s\n", string(body))
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

// Generic data generation steps
func iGenerateAUniqueEANCodeAs(varName string) error {
	// Generate a unique EAN code (13 digits)
	ean := fmt.Sprintf("%013d", rand.Int63n(10000000000000))
	savedVars[varName] = ean
	logger.Printf("Generated EAN code: %s\n", ean)
	return nil
}

func iGenerateAUniqueRFCAs(varName string) error {
	// Generate a unique RFC (Mexican tax ID)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rfc := ""
	for i := 0; i < 4; i++ {
		rfc += string(letters[rand.Intn(len(letters))])
	}
	rfc += fmt.Sprintf("%06d", rand.Intn(1000000))
	savedVars[varName] = rfc
	logger.Printf("Generated RFC: %s\n", rfc)
	return nil
}

// Helper functions
func substitute(path string) string {
	if strings.HasPrefix(path, "http") {
		return path
	}
	return base + path
}

func replaceVars(s string) string {
	for key, value := range savedVars {
		s = strings.ReplaceAll(s, "{"+key+"}", value)
	}
	return s
}

func addAuthHeader(req *http.Request) {
	if token, ok := savedVars["authToken"]; ok {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

// Helper function to get nested values from JSON
func getNestedValue(data map[string]interface{}, field string) (interface{}, error) {
	parts := strings.Split(field, ".")
	var current interface{} = data

	for _, part := range parts {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("field path %q not found in JSON response", field)
		}
		current, ok = m[part]
		if !ok {
			return nil, fmt.Errorf("field %q not found in JSON response", field)
		}
	}
	return current, nil
}

// Helper function to validate nested keys
func validateNestedKey(key string) error {
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}

	_, err := getNestedValue(data, key)
	return err
}

// Helper function to get type name
func getTypeName(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case int:
		return "integer"
	case bool:
		return "boolean"
	case nil:
		return "null"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

// Test suite initialization
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		logger.Println("Starting test suite...")
	})
	ctx.AfterSuite(func() {
		logger.Println("Test suite completed.")
	})
}

// Scenario initialization
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		logger.Printf("Starting scenario: %s\n", sc.Name)
		// Reset saved variables for each scenario
		savedVars = make(map[string]string)
		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		logger.Printf("Completed scenario: %s\n", sc.Name)
		return ctx, nil
	})

	// Register all steps
	ctx.Step(`^the service is initialized$`, theServiceIsInitialized)
	ctx.Step(`^I send a (GET|POST|PUT|DELETE|PATCH) request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^I send a (GET|POST|PUT|DELETE|PATCH) request to "([^"]*)" with body:$`, iSendARequestWithBody)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response code should be (\d+) or (\d+)$`, theResponseCodeShouldBeOr)
	ctx.Step(`^the JSON response should contain key "([^"]*)"$`, theJSONResponseShouldContainKey)
	ctx.Step(`^the JSON response should contain "([^"]*)" "([^"]*)"$`, theJSONResponseShouldContain)
	ctx.Step(`^the JSON response field "([^"]*)" should be of type "([^"]*)"$`, theJSONResponseFieldShouldBeOfType)
	ctx.Step(`^the JSON response should be an array$`, theJSONResponseShouldBeAnArray)
	ctx.Step(`^the JSON response array should have length (\d+)$`, theJSONResponseArrayShouldHaveLength)
	ctx.Step(`^I save the JSON response key "([^"]*)" as "([^"]*)"$`, iSaveTheJSONResponseKeyAs)
	ctx.Step(`^I save first array element key "([^"]*)" as "([^"]*)"$`, iSaveFirstArrayElementKeyAs)
	ctx.Step(`^I send a PATCH request to "([^"]*)"$`, iSendAPatchRequestTo)
	ctx.Step(`^the JSON response field "([^"]*)" should contain string "([^"]*)"$`, theJSONResponseFieldShouldContainString)
	ctx.Step(`^the JSON response should contain boolean "([^"]*)" (true|false)$`, theJSONResponseShouldContainBoolean)
	ctx.Step(`^I generate a unique EAN code as "([^"]*)"$`, iGenerateAUniqueEANCodeAs)
	ctx.Step(`^I generate a unique RFC as "([^"]*)"$`, iGenerateAUniqueRFCAs)
}
