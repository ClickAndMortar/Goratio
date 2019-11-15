# Goratio

[![Go Report Card](https://goreportcard.com/badge/github.com/ClickAndMortar/Goratio)](https://goreportcard.com/report/github.com/ClickAndMortar/Goratio) <a href='https://github.com/jpoles1/gopherbadger' target='_blank'>![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-70%25-brightgreen.svg?longCache=true&style=flat)</a>

Goratio is a location validation microservice for:

* Postal codes / zip codes (using [Unicode Common Locale Data Repository](http://cldr.unicode.org/) regexes)
* Phone numbers (using [Google's libphonenumber](https://github.com/google/libphonenumber))
* Email addresses
* IP addresses with GeoIP (using [Maxmind's GeoLite2 DB](https://dev.maxmind.com/geoip/geoip2/geolite2/))
* VAT numbers (coming)

## Usage

See [OpenAPI specification](https://app.swaggerhub.com/apis-docs/Click-and-Mortar/Goratio/1.1.0#/) for full API doc.

With the application running, post your query to the `/validate` endpoint:

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
    "email": "john.doe@example.com",
    "ip": "3.3.3.3"
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
    },
    "ip": {
        "address": "3.3.3.3",
        "valid": true,
        "geo": {
            "country_code": "US",
            "country_name": "United States",
            "city": "Seattle"
        }
    }
}
```

### GeoIP

Path to a GeoLite2 DB (`.mmdb` format) defaults to `/var/GeoLite2.mmdb`.

It's path may be defined explicitely using `GEOIP_DB_PATH` environment variable.

If the file at the given path does not exist or is not valid, GeoIP feature is disabled.

ℹ️  Note that GeoLite2 Country database might be slightly faster than the City DB. Consider using it if you don't need a city-level accuracy.

#### Downloading the DB

```bash
wget https://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz -O GeoLite2.tgz
tar --strip-components=1 -zxf GeoLite2.tar.gz *.mmdb
```

## Testing

```
make test
```
