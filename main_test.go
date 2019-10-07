package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidZipcodes(t *testing.T) {
	zip := Zip{
		Code:    "06000",
		Country: "FR",
	}

	response := ValidateZipCode(zip)
	if !*response.Valid {
		t.Errorf("Expecting valid FR zip code 06000")
	}

	zip = Zip{
		Code:    "3999",
		Country: "BE",
	}

	response = ValidateZipCode(zip)
	if !*response.Valid {
		t.Errorf("Expecting valid BE zipcode 3999")
	}
}

func TestInvalidZipcodes(t *testing.T) {
	zip := Zip{
		Code:    "1234",
		Country: "FR",
	}

	response := ValidateZipCode(zip)
	if *response.Valid {
		t.Errorf("Expecting invalid FR zip code 1234")
	}

	zip = Zip{
		Code:    "999",
		Country: "BE",
	}

	response = ValidateZipCode(zip)
	if *response.Valid {
		t.Errorf("Expecting invalid BE zip code 999")
	}
}

func TestValidPhone(t *testing.T) {
	phone := Phone{
		Number:  "0612346578",
		Country: "FR",
	}

	response := ValidatePhone(phone)
	if !*response.Valid {
		t.Errorf("Expecting valid FR phone number 0612345678")
	}

	phone = Phone{
		Number:  "0798139558‬",
		Country: "CH",
	}

	response = ValidatePhone(phone)
	if !*response.Valid {
		t.Errorf("Expecting valid CH phone number 0798139558‬")
	}
}

func TestInvalidPhone(t *testing.T) {
	phone := Phone{
		Number:  "06123465",
		Country: "FR",
	}

	response := ValidatePhone(phone)
	if *response.Valid {
		t.Errorf("Expecting invalid FR phone number 06123456")
	}

	phone = Phone{
		Number:  "07981395‬",
		Country: "CH",
	}

	response = ValidatePhone(phone)
	if *response.Valid {
		t.Errorf("Expecting invalid CH phone number 07981395‬")
	}
}

func TestHttpValidation(t *testing.T) {
	req, err := http.NewRequest("POST", "/validation", strings.NewReader(`{ "phone": { "country": "FR", "number": "0123456789" }, "zip": { "code": "06000", "country": "FR" }, "email": "john.doe@example.com" }`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ValidateLocationHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	rl := ResponseLocation{}
	err = json.Unmarshal(rr.Body.Bytes(), &rl)
	if err != nil {
		t.Errorf("Unable to JSON-decode response: %s", err)
	}

	if *rl.Phone.Valid != true {
		t.Errorf("Invalid phone number, expected valid")
	}

	if *rl.Zip.Valid != true {
		t.Errorf("Invalid zip code, expected valid")
	}

	if *rl.Email.Valid != true {
		t.Errorf("Invalid email address, expected valid")
	}
}

func TestHttpValidationZip(t *testing.T) {
	req, err := http.NewRequest("POST", "/validation/zip", strings.NewReader(`{ "zip": { "code": "0600", "country": "FR" } }`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ValidateZipHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	rz := ResponseZip{}
	err = json.Unmarshal(rr.Body.Bytes(), &rz)
	if err != nil {
		t.Errorf("Unable to JSON-decode response: %s", err)
	}

	if *rz.Valid == true {
		t.Errorf("Expected invalid zip code")
	}
}

func TestHttpValidationPhone(t *testing.T) {
	req, err := http.NewRequest("POST", "/validation/zip", strings.NewReader(`{ "phone": { "country": "FR", "number": "012345679" } }`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ValidatePhoneHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	rp := ResponsePhone{}
	err = json.Unmarshal(rr.Body.Bytes(), &rp)
	if err != nil {
		t.Errorf("Unable to JSON-decode response: %s", err)
	}

	if *rp.Valid == true {
		t.Errorf("Expected invalid phone number")
	}
}
