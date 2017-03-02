package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/wobeproject/handlers"
)

var rr *httptest.ResponseRecorder
var form url.Values
var r *http.Request

func Test_NewParser(t *testing.T) {
	p := handlers.NewParser(nil, "application/json")
	if p != nil {
		t.Fatalf("Parser Suppose to be nil")
	}
}

func Test_ReqParser(t *testing.T) {
	initTest()
	p := handlers.NewParser(rr, "application/x-www-form-urlencoded")

	form.Add("input", "abcdef")
	r, _ = http.NewRequest("POST", "/", strings.NewReader(form.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	p.RequestParse(r)
	if rr == nil {
		t.Fatal("ResponseWriter not suppose to be nil ", rr)
	}

	if rr.Code != 200 && strings.Contains(rr.Body.String(), "fedcba") {
		t.Fatal("output not matching")
	}
}

func Test_InputHandler(t *testing.T) {
	initTest()
	form.Add("input", "abcdef")
	r, _ = http.NewRequest("POST", "/", strings.NewReader(form.Encode()))

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(handlers.InputHandler)

	handler.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func Test_ValidationHandler(t *testing.T) {
	initTest()
	form.Add("input", "abcdef")

	r.Header.Add("Content-Type", "")
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	handler := handlers.ValidationHandler(testHandler)
	handler.ServeHTTP(rr, r)
	if rr.Code != 400 {
		t.Fatal("Suppose to be 400 error, got ", rr.Code)
	}
}

func Test_IPHandler400(t *testing.T) {
	initTest()
	form.Add("input", "")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	inHandler := http.HandlerFunc(handlers.InputHandler)
	handler := handlers.ValidationHandler(inHandler)
	handler.ServeHTTP(rr, r)
	if rr.Code != 422 {
		t.Fatal("Suppose to be 422 error, got ", rr.Code)
	}
}

func Test_ValidationHandler415(t *testing.T) {
	initTest()
	form.Add("input", "dabf")
	r.Header.Add("Content-Type", "application/json")
	inHandler := http.HandlerFunc(handlers.InputHandler)
	handler := handlers.ValidationHandler(inHandler)
	handler.ServeHTTP(rr, r)
	if rr.Code != 415 {
		t.Fatal("Suppose to be 415 error, got ", rr.Code)
	}
}

func Test_PanicRecover(t *testing.T) {
	initTest()
	form.Add("input", "dabef")

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	inHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Random panic")
	})
	handler := handlers.RecoverHandler(inHandler)
	handler = handlers.ValidationHandler(handler)
	handler.ServeHTTP(rr, r)
	if rr.Code != 500 {
		t.Fatal("Suppose to be 500 error, got ", rr.Code)
	}
}

func initTest() {
	rr = httptest.NewRecorder()
	form = url.Values{}
	r, _ = http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
}
