package tyrgin

import (
	"encoding/json"
	"testing"
)

func TestTraverse(t *testing.T) {

	traverseResponse := Traverse(testStatusEndpoints, []string{}, "", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	testAboutResponse := AboutResponse{}
	err := json.Unmarshal([]byte(traverseResponse), &testAboutResponse)
	if err != nil {
		t.Errorf("Response body is an invalid About format, was: `%s`", traverseResponse)
	}

	assertEqualAboutData(t, testAboutResponse, emptyCustomData, defaultServiceID)
}

func TestTraverseInvalidDependency(t *testing.T) {

	traverseResponse := Traverse(testStatusEndpoints, []string{"something"}, "", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	expected := `["CRIT",{"description":"Can't traverse","result":"CRIT","details":"Status path 'something' is not registered"}]`
	if traverseResponse != expected {
		t.Errorf("Response body should be `%s`, was: `%s`", expected, traverseResponse)
	}
}

func TestTraverseInvalidAction(t *testing.T) {

	traverseResponse := Traverse(testStatusEndpoints, []string{}, "something", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	expected := `["CRIT",{"description":"Unsupported action","result":"CRIT","details":"Unsupported traversal action 'something'"}]`
	if traverseResponse != expected {
		t.Errorf("Response body should be `%s`, was: `%s`", expected, traverseResponse)
	}
}

func TestTraverseNotTraversable(t *testing.T) {

	se := []StatusEndpoint{
		testStatusEndpointA,
		testStatusEndpointB,
		testStatusEndpointC,
		testStatusEndpointNotTraversable,
	}

	traverseResponse := Traverse(se, []string{"sss"}, "", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	expected := `["CRIT",{"description":"Can't traverse","result":"CRIT","details":"SSS is not traversable"}]`
	if traverseResponse != expected {
		t.Errorf("Response body should be `%s`, was: `%s`", expected, traverseResponse)
	}
}

func TestTraverseMissingTraverseCheck(t *testing.T) {

	se := []StatusEndpoint{
		testStatusEndpointA,
		testStatusEndpointB,
		testStatusEndpointC,
		testStatusEndpointMissingTraverseChecker,
	}

	traverseResponse := Traverse(se, []string{"ttt"}, "", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	expected := `["CRIT",{"description":"Can't traverse","result":"CRIT","details":"TTT does not have a TraverseCheck() function defined"}]`
	if traverseResponse != expected {
		t.Errorf("Response body should be `%s`, was: `%s`", expected, traverseResponse)
	}
}

func TestTraverseDependencyFound(t *testing.T) {

	se := []StatusEndpoint{
		testStatusEndpointA,
		testStatusEndpointB,
		testStatusEndpointTraversable,
	}

	traverseResponse := Traverse(se, []string{"uuu"}, "", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	expected := `{"Name":"UUU","Body":"Hello","Time":1294706395881547000}`
	if traverseResponse != expected {
		t.Errorf("Response body should be `%s`, was: `%s`", expected, traverseResponse)
	}
}

func TestTraverseDependencyFoundWithError(t *testing.T) {

	se := []StatusEndpoint{
		testStatusEndpointA,
		testStatusEndpointB,
		testStatusEndpointTraversableError,
	}

	traverseResponse := Traverse(se, []string{"vvv"}, "", AboutProtocolHTTP, "test/about.json", "test/version.txt", emptyCustomData)

	expected := `["CRIT",{"description":"Traverse","result":"CRIT","details":"Test Error"}]`
	if traverseResponse != expected {
		t.Errorf("Response body should be `%s`, was: `%s`", expected, traverseResponse)
	}
}
