# Goratio

Goratio is a location validation microservice for:

* Postal codes / zip codes (using [Unicode Common Locale Data Repository](http://cldr.unicode.org/) regexes)
* Phone numbers (using [Google's libphonenumber](https://github.com/google/libphonenumber))
* Email addresses
* VAT numbers (coming)

## Usage

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
