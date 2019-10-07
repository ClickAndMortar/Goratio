# Goratio

[![Go Report Card](https://goreportcard.com/badge/github.com/ClickAndMortar/Goratio)](https://goreportcard.com/report/github.com/ClickAndMortar/Goratio) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-70%25-brightgreen.svg?longCache=true&style=flat)</a>

Goratio is a location validation microservice for:

* Postal codes / zip codes (using [Unicode Common Locale Data Repository](http://cldr.unicode.org/) regexes)
* Phone numbers (using [Google's libphonenumber](https://github.com/google/libphonenumber))
* Email addresses
* VAT numbers (coming)

## Usage

See [OpenAPI specification](https://app.swaggerhub.com/apis-docs/Click-and-Mortar/Goratio/1.0.0#/) for full API doc.

With the application running, post your query to the `/validation` endpoint:

```json
{
    "phone": {
        "number": "0612345678",
        "country": "FR"
    },
    "zip": {
        "code": "06000",
        "country": "FR"
    },
    "email": "john.doe@example.com"
}
```

Will output the following result:

```json
{
    "phone": {
        "number": "0612345678",
        "country": "FR",
        "valid": true,
        "formatted": {
            "E164": "+33612345678",
            "national": "06 12 34 56 78",
            "international": "+33 6 12 34 56 78"
        }
    },
    "zip": {
        "code": "06000",
        "country": "FR",
        "valid": true
    },
    "email": {
        "address": "john.doe@example.com",
        "valid": true
    }
}
```

## Testing

```
make test
```
