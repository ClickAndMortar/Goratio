package main

import (
	"encoding/json"
	"fmt"
	"github.com/ClickAndMortar/Goratio/zipcode"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/oschwald/geoip2-golang"
	"github.com/ttacon/libphonenumber"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	geoDB     *geoip2.Reader
	geoEnable bool
	err       error
)

type Location struct {
	Email   string `json:"email,omitempty"`
	Phone   Phone  `json:"phone,omitempty"`
	Zip     Zip    `json:"zip,omitempty"`
	IP      string `json:"ip,omitempty"`
	VatCode string `json:"vat_code,omitempty"`
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
	IP    ResponseIP    `json:"ip,omitempty"`
	Vat   ResponseVat   `json:"vat,omitempty"`
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

type ResponseIP struct {
	Address string        `json:"address,omitempty"`
	Valid   *bool         `json:"valid,omitempty"`
	Geo     ResponseIPGeo `json:"geo,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type ResponseIPGeo struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	City        string `json:"city"`
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
	geoEnable = true
	geodbPath := getEnvDefault("GEOIP_DB_PATH", "/var/GeoIP2.mmdb")
	if _, err := os.Stat(geodbPath); os.IsNotExist(err) {
		geoEnable = false
		fmt.Printf("GeoIP DB does not exist at [%s], feature disabled", geodbPath)
	}

	geoDB, err = geoip2.Open(geodbPath)
	if err != nil {
		geoEnable = false
		fmt.Printf("GeoIP DB loading error from [%s]: %s", geodbPath, err)
	}
	defer geoDB.Close()

	http.HandleFunc("/validate", ValidateLocationHandler)
	http.HandleFunc("/validate/zip", ValidateZipHandler)
	http.HandleFunc("/validate/phone", ValidatePhoneHandler)
	http.HandleFunc("/validate/ip", ValidateIPHandler)

	port := getEnvDefault("GORATIO_PORT", "8080")

	log.Printf("Listening on http://127.0.0.1:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func ValidateLocationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	location, errorResponse := ExtractLocation(r)
	response, errorResponse := ValidateLocation(location, errorResponse)

	if errorResponse.Error != "" {
		errorOutput, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(errorOutput)
		return
	}

	output, _ := json.Marshal(response)
	w.Write(output)
}

func ValidateZipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	location, errorResponse := ExtractLocation(r)
	response, errorResponse := ValidateLocation(location, errorResponse)

	if errorResponse.Error != "" {
		errorOutput, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(errorOutput)
		return
	}

	output, _ := json.Marshal(response.Zip)
	w.Write(output)
}

func ValidatePhoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	location, errorResponse := ExtractLocation(r)
	response, errorResponse := ValidateLocation(location, errorResponse)

	if errorResponse.Error != "" {
		errorOutput, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(errorOutput)
		return
	}

	output, _ := json.Marshal(response.Phone)
	w.Write(output)
}

func ValidateIPHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	location, errorResponse := ExtractLocation(r)
	response, errorResponse := ValidateLocation(location, errorResponse)

	if errorResponse.Error != "" {
		errorOutput, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(errorOutput)
		return
	}

	output, _ := json.Marshal(response.IP)
	w.Write(output)
}

func ExtractLocation(r *http.Request) (*Location, ErrorResponse) {
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

	return location, errorResponse
}

func ValidateLocation(location *Location, error ErrorResponse) (*ResponseLocation, ErrorResponse) {
	response := &ResponseLocation{}

	err := location.Validate()
	if err != nil {
		error.Error = err.Error()
		return response, error
	}

	response.Phone = ValidatePhone(location.Phone)
	response.Zip = ValidateZipCode(location.Zip)
	response.Email = ValidateEmail(location.Email)
	response.IP = ValidateIP(location.IP)
	response.Vat = ValidateVat(location.VatCode)

	return response, error
}

func ValidateZipCode(zip Zip) ResponseZip {
	response := ResponseZip{}
	if zip.Code == "" {
		return response
	}

	response.Code = zip.Code
	response.Country = zip.Country
	matched := zipcode.Validate(zip.Code, zip.Country)
	response.Valid = &matched

	return response
}

func ValidatePhone(phone Phone) ResponsePhone {
	response := ResponsePhone{}
	if phone.Number == "" {
		return response
	}

	response.Number = phone.Number
	response.Country = phone.Country
	num, err := libphonenumber.Parse(phone.Number, phone.Country)
	if err != nil {
		response.Error = err.Error()
	}

	valid := libphonenumber.IsValidNumberForRegion(num, phone.Country)
	response.Valid = &valid

	if valid {
		response.Formatted.E164 = libphonenumber.Format(num, libphonenumber.E164)
		response.Formatted.National = libphonenumber.Format(num, libphonenumber.NATIONAL)
		response.Formatted.International = libphonenumber.Format(num, libphonenumber.INTERNATIONAL)
	}

	return response
}

func ValidateEmail(email string) ResponseEmail {
	response := ResponseEmail{}
	if email == "" {
		return response
	}

	response.Address = email
	err := validation.Validate(email, is.Email)
	emailValid := err == nil
	response.Valid = &emailValid

	return response
}

func ValidateIP(ip string) ResponseIP {
	response := ResponseIP{}
	responseGeo := ResponseIPGeo{}
	if ip == "" {
		return response
	}

	isValid := false
	response.Address = ip
	response.Valid = &isValid
	ipObj := net.ParseIP(ip)
	if ipObj == nil {
		response.Error = fmt.Sprintf("Invalid IP address [%s]", ip)
		return response
	}

	isValid = true
	response.Valid = &isValid

	if !geoEnable {
		return response
	}

	record, err := geoDB.City(ipObj)
	if err != nil {
		response.Error = err.Error()
		return response
	}

	responseGeo.CountryCode = record.Country.IsoCode
	responseGeo.CountryName = record.Country.Names["en"]
	responseGeo.City = record.City.Names["en"]
	response.Geo = responseGeo

	return response
}

func ValidateVat(code string) ResponseVat {
	response := ResponseVat{}
	if code == "" {
		return response
	}

	response.Code = code
	response.Error = "Not yet implemented"

	return response
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
