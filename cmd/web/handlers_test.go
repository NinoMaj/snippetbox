package main

import (
	"net/http"
	"testing"
)

// Unit test
// func TestPing(t *testing.T) {
// 	t.Parallel() // Mark that it is OK to run this test in parallel
// 	// Initialize a new httptest.ResponseRecorder.
// 	rr := httptest.NewRecorder()

// 	// Initialize a new dummy http.Request.
// 	r, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Call the ping handler function, passing in the
// 	// httptest.ResponseRecorder and http.Request.
// 	ping(rr, r)

// 	// Call the Result() method on the http.ResponseRecorder to get the
// 	// http.Response generated by the ping handler.
// 	rs := rr.Result()

// 	// We can then examine the http.Response to check that the status code
// 	// written by the ping handler was 200.
// 	if rs.StatusCode != http.StatusOK {
// 		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
// 	}

// 	// And we can check that the response
// 	// body written by the ping handler equals "OK".
// 	defer rs.Body.Close()
// 	body, err := ioutil.ReadAll(rs.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if string(body) != "OK" {
// 		t.Errorf("want body to equal %q", "OK")
// 	}
// }

// E2E test
func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
