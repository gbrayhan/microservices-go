//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cucumber/godog"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
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

func theServiceIsInitialized() error {
	logger.Println("Service initialized.")
	return nil
}

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

func theJSONResponseShouldContainKey(key string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check for key %q", key)
	}
	logger.Printf("Checking if JSON response contains key: %q\n", key)

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

func theJSONResponseShouldContain(field, value string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check field %q", field)
	}
	expectedValue := replaceVars(value)
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}
	rawValue, ok := data[field]
	if !ok {
		return fmt.Errorf("field %q not found in JSON response", field)
	}
	actualValue := fmt.Sprintf("%v", rawValue)
	if actualValue != expectedValue {
		return fmt.Errorf("expected %q = %q, but got %v", field, expectedValue, actualValue)
	}
	return nil
}

func iSaveTheJSONResponseKeyAs(key, varName string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot save key %q", key)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON for saving: %v. Body: %s", errUnmarshal, string(body))
	}
	raw, ok := data[key]
	if !ok {
		return fmt.Errorf("key %q not found in response for saving", key)
	}
	var valStr string
	switch v := raw.(type) {
	case string:
		valStr = v
	case float64:
		if v == float64(int64(v)) {
			valStr = fmt.Sprintf("%.0f", v)
		} else {
			valStr = fmt.Sprintf("%f", v)
		}
	case bool:
		valStr = fmt.Sprintf("%t", v)
	case nil:
		valStr = ""
	default:
		valStr = fmt.Sprintf("%v", v)
	}
	savedVars[varName] = valStr
	return nil
}

func substitute(path string) string {
	substituted := path
	for k, v := range savedVars {
		substituted = strings.ReplaceAll(substituted, "${"+k+"}", v)
	}
	return base + substituted
}

func replaceVars(s string) string {
	res := s
	for k, v := range savedVars {
		res = strings.ReplaceAll(res, "${"+k+"}", v)
	}
	return res
}

func addAuthHeader(req *http.Request) {
	if token, ok := savedVars["accessToken"]; ok && token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

func iGenerateAUniqueEANCodeAs(varName string) error {
	uniqueEAN := fmt.Sprintf("EAN%d%d", time.Now().UnixNano(), rand.Intn(1000))
	savedVars[varName] = uniqueEAN
	return nil
}

func iGenerateAUniqueRFCAs(varName string) error {
	uniqueRFC := fmt.Sprintf("RFC%d%d", time.Now().UnixNano()/int64(time.Millisecond), rand.Intn(1000))
	savedVars[varName] = uniqueRFC
	return nil
}

func iSaveFirstArrayElementKeyAs(key, arrayKey, varName string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot save key %q from array %q", key, arrayKey)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON for saving from array: %v. Body: %s", errUnmarshal, string(body))
	}
	arrRaw, ok := data[arrayKey]
	if !ok {
		return fmt.Errorf("array %q not found in JSON response", arrayKey)
	}
	arr, ok := arrRaw.([]interface{})
	if !ok || len(arr) == 0 {
		return fmt.Errorf("array %q is empty or not an array", arrayKey)
	}
	obj, ok := arr[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("first element in %q is not an object", arrayKey)
	}
	val, ok := obj[key]
	if !ok {
		return fmt.Errorf("key %q not found in first element of %q", key, arrayKey)
	}
	var valStr string
	switch v := val.(type) {
	case string:
		valStr = v
	case float64:
		if v == float64(int64(v)) {
			valStr = fmt.Sprintf("%.0f", v)
		} else {
			valStr = fmt.Sprintf("%f", v)
		}
	case bool:
		valStr = fmt.Sprintf("%t", v)
	case nil:
		valStr = ""
	default:
		valStr = fmt.Sprintf("%v", v)
	}
	savedVars[varName] = valStr
	return nil
}

func iSendAPatchRequestTo(path string) error {
	fullURL := substitute(path)
	req, e := http.NewRequest("PATCH", fullURL, nil)
	if e != nil {
		return e
	}
	addAuthHeader(req)
	resp, err = http.DefaultClient.Do(req)
	return err
}

func theJSONResponseFieldShouldContainString(fieldName, expectedSubstring string) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check field %q for substring %q", fieldName, expectedSubstring)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}
	raw, ok := data[fieldName]
	if !ok {
		return fmt.Errorf("field %q not found in JSON response", fieldName)
	}
	strVal, ok := raw.(string)
	if !ok {
		return fmt.Errorf("field %q is not a string (got %T)", fieldName, raw)
	}
	if !strings.Contains(strVal, expectedSubstring) {
		return fmt.Errorf("expected field %q to contain %q, but it did not (value: %q)", fieldName, expectedSubstring, strVal)
	}
	return nil
}

func theJSONResponseShouldContainBoolean(fieldName string, expectedValue bool) error {
	if body == nil {
		return fmt.Errorf("response body is nil, cannot check field %q for boolean %t", fieldName, expectedValue)
	}
	var data map[string]interface{}
	if errUnmarshal := json.Unmarshal(body, &data); errUnmarshal != nil {
		return fmt.Errorf("error unmarshalling JSON: %v. Body: %s", errUnmarshal, string(body))
	}
	raw, ok := data[fieldName]
	if !ok {
		return fmt.Errorf("field %q not found in JSON response", fieldName)
	}
	boolVal, ok := raw.(bool)
	if !ok {
		return fmt.Errorf("field %q is not a boolean (got %T)", fieldName, raw)
	}
	if boolVal != expectedValue {
		return fmt.Errorf("expected field %q to be %t, but got %t", fieldName, expectedValue, boolVal)
	}
	return nil
}

func theResponseCodeShouldBeOr(code1, code2 int) error {
	if resp == nil {
		return fmt.Errorf("response is nil, cannot check status code")
	}
	if resp.StatusCode != code1 && resp.StatusCode != code2 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			logger.Printf("Error closing response body: %v\n", err)
			return err
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return fmt.Errorf("expected status code %d or %d but got %d. Body: %s", code1, code2, resp.StatusCode, string(bodyBytes))
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return nil
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	rand.New(rand.NewSource(12345))

	ctx.BeforeSuite(func() {
		payload := `{"email":"gbrayhan@gmail.com","password":"qweqwe"}`
		req, errLogin := http.NewRequest("POST", base+"/login", bytes.NewBufferString(payload))
		if errLogin != nil {
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")
		loginResp, errLogin := http.DefaultClient.Do(req)
		if errLogin != nil {
			os.Exit(1)
		}
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				logger.Printf("Error closing login response body: %v\n", err)
				os.Exit(1)
			}
		}(loginResp.Body)
		loginBodyBytes, errRead := io.ReadAll(loginResp.Body)
		if errRead != nil {
			os.Exit(1)
		}
		var data map[string]interface{}
		if errJson := json.Unmarshal(loginBodyBytes, &data); errJson != nil {
			os.Exit(1)
		}
		token, ok := data["accessToken"].(string)
		if !ok || token == "" {
			os.Exit(1)
		}
		savedVars["accessToken"] = token
	})
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctxHook context.Context, sc *godog.Scenario) (context.Context, error) {
		resp = nil
		body = nil
		return ctxHook, nil
	})
	ctx.Step(`^the service is initialized$`, theServiceIsInitialized)
	ctx.Step(`^I send a (GET|DELETE) request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^I send a (POST|PUT) request to "([^"]*)" with body:$`, iSendARequestWithBody)
	ctx.Step(`^the response code should be (\d+)$`, theResponseCodeShouldBe)
	ctx.Step(`^the response code should be (\d+) or (\d+)$`, theResponseCodeShouldBeOr)
	ctx.Step(`^the JSON response should contain key "([^"]*)"$`, theJSONResponseShouldContainKey)
	ctx.Step(`^the JSON response should contain "([^"]*)": "([^"]*)"$`, theJSONResponseShouldContain)
	ctx.Step(`^I save the JSON response key "([^"]*)" as "([^"]*)"$`, iSaveTheJSONResponseKeyAs)
	ctx.Step(`^the JSON response field "([^"]*)" should contain string "([^"]*)"$`, theJSONResponseFieldShouldContainString)
	ctx.Step(`^the JSON response should contain "([^"]*)": true$`, func(field string) error {
		return theJSONResponseShouldContainBoolean(field, true)
	})
	ctx.Step(`^the JSON response should contain "([^"]*)": false$`, func(field string) error {
		return theJSONResponseShouldContainBoolean(field, false)
	})
	ctx.Step(`^I generate a unique EAN code as "([^"]*)"$`, iGenerateAUniqueEANCodeAs)
	ctx.Step(`^I generate a unique RFC as "([^"]*)"$`, iGenerateAUniqueRFCAs)
	ctx.Step(`^I save the first element key "([^"]*)" from array "([^"]*)" as "([^"]*)"$`, iSaveFirstArrayElementKeyAs)
	ctx.Step(`^I send a PATCH request to "([^"]*)"$`, iSendAPatchRequestTo)
}
