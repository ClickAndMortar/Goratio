package main

import (
	"encoding/json"
	"fmt"
	"github.com/clickandmortar/goratio/zipcode"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/ttacon/libphonenumber"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Location struct {
	Email string `json:"email,omitempty"`
	Phone Phone  `json:"phone,omitempty"`
	Zip   Zip    `json:"zip,omitempty"`
	VatCode   string    `json:"vat_code,omitempty"`
}

type Phone struct {
	Number  string `json:"number,omitempty"`
	Country string `json:"country,omitempty"`
}

type Zip struct {
	Code    string `json:"code,omitempty"`
	Country string `json:"country,omitempty"`
}

type ResponseLocation struct {
	Phone ResponsePhone `json:"phone,omitempty"`
	Zip   ResponseZip   `json:"zip,omitempty"`
	Email ResponseEmail `json:"email,omitempty"`
	Vat  ResponseVat   `json:"vat,omitempty"`
}

type ResponsePhone struct {
	Number    string         `json:"number,omitempty"`
	Country   string         `json:"country,omitempty"`
	Valid     *bool          `json:"valid,omitempty"`
	Error     string         `json:"error,omitempty"`
	Formatted FormattedPhone `json:"formatted,omitempty"`
}

type ResponseZip struct {
	Code    string `json:"code,omitempty"`
	Country string `json:"country,omitempty"`
	Valid   *bool  `json:"valid,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ResponseEmail struct {
	Address string `json:"address,omitempty"`
	Valid   *bool  `json:"valid,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ResponseVat struct {
	Code  string `json:"code,omitempty"`
	Valid *bool  `json:"valid,omitempty"`
	Error string `json:"error,omitempty"`
}

type FormattedPhone struct {
	E164          string `json:"E164,omitempty"`
	National      string `json:"national,omitempty"`
	International string `json:"international,omitempty"`
}

type ErrorResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error"`
}

func main() {
	http.HandleFunc("/validate", validateHandler)

	port := getEnvDefault("GORATIO_PORT", "8080")

	log.Printf("Listening on http://127.0.0.1:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	errorResponse := ErrorResponse{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorResponse.Error = err.Error()
	}

	location := &Location{}
	err = json.Unmarshal(body, location)
	if err != nil {
		errorResponse.Error = err.Error()
	}

	err = location.Validate()
	if err != nil {
		errorResponse.Error = err.Error()
	}

	if errorResponse.Error != "" {
		errorOutput, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(errorOutput)
		return
	}

	response := &ResponseLocation{}

	if location.Phone.Number != "" {
		response.Phone.Number = location.Phone.Number
		response.Phone.Country = location.Phone.Country
		num, err := libphonenumber.Parse(location.Phone.Number, location.Phone.Country)
		if err != nil {
			response.Phone.Error = err.Error()
		}

		valid := libphonenumber.IsValidNumberForRegion(num, location.Phone.Country)
		response.Phone.Valid = &valid

		if valid {
			response.Phone.Formatted.E164 = libphonenumber.Format(num, libphonenumber.E164)
			response.Phone.Formatted.National = libphonenumber.Format(num, libphonenumber.NATIONAL)
			response.Phone.Formatted.International = libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
		}
	}

	if location.Zip.Code != "" {
		response.Zip.Code = location.Zip.Code
		response.Zip.Country = location.Zip.Country
		matched := zipcode.Validate(location.Zip.Code, location.Zip.Country)
		response.Zip.Valid = &matched
	}

	if location.Email != "" {
		response.Email.Address = location.Email
		err = validation.Validate(location.Email, is.Email)
		emailValid := err == nil
		response.Email.Valid = &emailValid
	}

	if location.VatCode != "" {
		response.Vat.Code = location.VatCode
		response.Vat.Error = "Not yet implemented"
	}

	output, _ := json.Marshal(response)
	w.Write(output)
}

func (a Zip) Validate() error {
	if a.Country == "" && a.Code == "" {
		return nil
	}

	return validation.ValidateStruct(&a,
		validation.Field(&a.Country, validation.Required, is.CountryCode2),
		validation.Field(&a.Code, validation.Required, validation.Length(2, 20)),
	)
}

func (p Phone) Validate() error {
	if p.Number == "" && p.Country == "" {
		return nil
	}

	return validation.ValidateStruct(&p,
		validation.Field(&p.Number, validation.Required, validation.Length(4, 20)),
		validation.Field(&p.Country, validation.Required, is.CountryCode2),
	)
}

func (l Location) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Zip),
		validation.Field(&l.Phone),
	)
}

func getEnvDefault(name string, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" && defaultValue != "" {
		log.Printf("Environment variable %s not set or empty, using default %s", name, defaultValue)
		variable = defaultValue
	}

	return variable
}
