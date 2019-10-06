package main

import "testing"

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
