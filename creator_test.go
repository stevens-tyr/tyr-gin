package tyrgin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type respTest struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Uploaded   string `json:"uploaded"`
	Count      int    `json:"count"`
}

type reqTest struct {
	Upload string `json:"uploaded" binding:"required"`
	Count  int    `json:"count" binding:"required"`
}

var uploaded = ""
var count = 0

func testGetFunc(c *gin.Context) {
	count++

	c.JSON(
		http.StatusOK,
		gin.H{
			"status_code": http.StatusOK,
			"message":     "Hello World!",
			"uploaded":    uploaded,
			"count":       count,
		},
	)
}

func testDeleteFunc(c *gin.Context) {
	uploaded = ""
	count = 0

	c.JSON(
		http.StatusResetContent,
		gin.H{
			"status_code": http.StatusResetContent,
			"message":     "Hello Deleted",
			"uploaded":    uploaded,
			"count":       count,
		},
	)
}

func testPatchFunc(c *gin.Context) {
	var json reqTest

	err := c.ShouldBindJSON(&json)
	ErrorLogger(err, "Failed to bind json for testPatchFunc.")

	if json.Upload != "" {
		uploaded = json.Upload
	}

	count++

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status_code": http.StatusCreated,
			"message":     "Hello Patch",
			"uploaded":    uploaded,
			"count":       count,
		},
	)
}

func testPostFunc(c *gin.Context) {
	var rb reqTest

	c.ShouldBindJSON(&rb)

	uploaded = rb.Upload
	count = rb.Count

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status_code": http.StatusCreated,
			"message":     "Hello Post",
			"uploaded":    uploaded,
			"count":       count,
		},
	)
}

func testPutFunc(c *gin.Context) {
	var json reqTest

	err := c.ShouldBindJSON(&json)
	ErrorLogger(err, "Failed to bind json for testPutFunc.")

	uploaded = json.Upload

	count = json.Count

	c.JSON(
		http.StatusCreated,
		gin.H{
			"status_code": http.StatusCreated,
			"message":     "Hello Put",
			"uploaded":    uploaded,
			"count":       count,
		},
	)
}

var (
	testGet    = NewRoute(testGetFunc, "hello", false, GET)
	testDelete = NewRoute(testDeleteFunc, "hello", false, DELETE)
	testPatch  = NewRoute(testPatchFunc, "hello", false, PATCH)
	testPost   = NewRoute(testPostFunc, "hello", false, POST)
	testPut    = NewRoute(testPutFunc, "hello", false, PUT)
)

var testEndpoints = []APIAction{
	testGet,
	testDelete,
	testPatch,
	testPost,
	testPut,
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if method == "GET" {
		req, _ = http.NewRequest(method, path, nil)
	} else {
		bb := new(bytes.Buffer)
		json.NewEncoder(bb).Encode(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}

func TestNotFound(t *testing.T) {
	body := gin.H{
		"statusCode": http.StatusNotFound,
		"message":    NotFoundError,
	}

	router := SetupRouter()

	resp := performRequest(router, "GET", "/", nil)
	assert.Equal(t, http.StatusNotFound, resp.Code)
	var response respTest
	err := json.Unmarshal([]byte(resp.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, body["statusCode"], response.StatusCode)
	assert.Equal(t, body["message"], response.Message)
}

func TestCreatorRoutes(t *testing.T) {
	router := SetupRouter()
	jwt := AuthMid()

	AddRoutes(router, jwt, "1", "tester", testEndpoints)
	var response respTest

	// Test Get
	gresp := performRequest(router, "GET", "/api/v1/tester/hello", nil)
	assert.Equal(t, http.StatusOK, gresp.Code)
	err := json.Unmarshal([]byte(gresp.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, 200, response.StatusCode)
	assert.Equal(t, "Hello World!", response.Message)
	assert.Equal(t, "", response.Uploaded)
	assert.Equal(t, 1, response.Count)

	var request reqTest

	// Test Post
	request = reqTest{
		Upload: "New Hello",
		Count:  100,
	}
	jsb, _ := json.Marshal(request)

	poresp := performRequest(router, "POST", "/api/v1/tester/hello", jsb)
	assert.Equal(t, http.StatusCreated, poresp.Code)
	err = json.Unmarshal([]byte(poresp.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, 201, response.StatusCode)
	assert.Equal(t, "Hello Post", response.Message)
	assert.Equal(t, "New Hello", response.Uploaded)
	assert.Equal(t, 100, response.Count)

	// Test Put
	request = reqTest{
		Upload: "New New Hello",
		Count:  1000,
	}
	jsb, _ = json.Marshal(request)

	puresp := performRequest(router, "PUT", "/api/v1/tester/hello", jsb)
	assert.Equal(t, http.StatusCreated, puresp.Code)
	err = json.Unmarshal([]byte(puresp.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, 201, response.StatusCode)
	assert.Equal(t, "Hello Put", response.Message)
	assert.Equal(t, "New New Hello", response.Uploaded)
	assert.Equal(t, 1000, response.Count)

	// Test Patch
	request = reqTest{
		Upload: "New New New Hello",
	}
	jsb, _ = json.Marshal(request)

	paresp := performRequest(router, "PATCH", "/api/v1/tester/hello", jsb)
	assert.Equal(t, http.StatusCreated, paresp.Code)
	err = json.Unmarshal([]byte(paresp.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, 201, response.StatusCode)
	assert.Equal(t, "Hello Patch", response.Message)
	assert.Equal(t, "New New New Hello", response.Uploaded)
	assert.Equal(t, 1001, response.Count)

	// Test Delete
	dresp := performRequest(router, "DELETE", "/api/v1/tester/hello", nil)
	assert.Equal(t, http.StatusResetContent, dresp.Code)
	err = json.Unmarshal([]byte(dresp.Body.String()), &response)
	assert.Nil(t, err)
	assert.Equal(t, 205, response.StatusCode)
	assert.Equal(t, "Hello Deleted", response.Message)
	assert.Equal(t, "", response.Uploaded)
	assert.Equal(t, 0, response.Count)
}
