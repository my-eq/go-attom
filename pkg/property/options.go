package property

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Option configures optional query parameters for Property API requests.
type Option func(values url.Values)

// applyOptions builds a url.Values map from the supplied options.
func applyOptions(opts []Option) url.Values {
	values := url.Values{}
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	return values
}

// WithString sets an arbitrary string parameter when the value is not empty.
func WithString(key, value string) Option {
	return func(values url.Values) {
		if key == "" || value == "" {
			return
		}
		values.Set(key, value)
	}
}

// WithStringSlice joins a slice of strings with the provided separator.
func WithStringSlice(key string, valuesList []string, separator string) Option {
	return func(values url.Values) {
		if key == "" || len(valuesList) == 0 {
			return
		}
		sep := separator
		if sep == "" {
			sep = "|"
		}
		values.Set(key, strings.Join(valuesList, sep))
	}
}

// WithAttomID sets the attomid query parameter.
func WithAttomID(attomID string) Option {
	return WithString("attomid", attomID)
}

// WithPropertyID sets the id query parameter for legacy property identifiers.
func WithPropertyID(id string) Option {
	return WithString("id", id)
}

// WithFIPSAndAPN sets the fips and apn parameters used for property lookups.
func WithFIPSAndAPN(fips, apn string) Option {
	return func(values url.Values) {
		if fips != "" {
			values.Set("fips", fips)
		}
		if apn != "" {
			values.Set("APN", apn)
		}
	}
}

// WithAddress sets the address parameter using a single formatted string.
func WithAddress(address string) Option {
	return WithString("address", address)
}

// WithAddressLines sets address1 and address2 query parameters.
func WithAddressLines(address1, address2 string) Option {
	return func(values url.Values) {
		if address1 != "" {
			values.Set("address1", address1)
		}
		if address2 != "" {
			values.Set("address2", address2)
		}
	}
}

// WithLatitudeLongitude adds latitude and longitude parameters.
func WithLatitudeLongitude(latitude, longitude float64) Option {
	return func(values url.Values) {
		values.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
		values.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	}
}

// WithRadius sets radius parameter expressed in miles.
func WithRadius(radiusMiles float64) Option {
	return func(values url.Values) {
		if radiusMiles <= 0 {
			return
		}
		values.Set("radius", strconv.FormatFloat(radiusMiles, 'f', -1, 64))
	}
}

// WithPostalCode sets the postalCode query parameter.
func WithPostalCode(code string) Option {
	return WithString("postalCode", code)
}

// WithCityName sets the cityname parameter.
func WithCityName(city string) Option {
	return WithString("cityname", city)
}

// WithGeoID sets the geoid parameter.
func WithGeoID(geoID string) Option {
	return WithString("geoid", geoID)
}

// WithGeoIDV4 sets the geoIdV4 parameter.
func WithGeoIDV4(geoID string) Option {
	return WithString("geoIdV4", geoID)
}

// WithPropertyType sets the propertytype parameter.
func WithPropertyType(propertyType string) Option {
	return WithString("propertytype", propertyType)
}

// WithPropertyIndicator sets the propertyIndicator parameter.
func WithPropertyIndicator(indicator int) Option {
	return func(values url.Values) {
		if indicator <= 0 {
			return
		}
		values.Set("propertyIndicator", strconv.Itoa(indicator))
	}
}

// WithBedsRange sets minimum and maximum beds filters.
func WithBedsRange(minBeds, maxBeds int) Option {
	return func(values url.Values) {
		if minBeds > 0 {
			values.Set("minBeds", strconv.Itoa(minBeds))
		}
		if maxBeds > 0 {
			values.Set("maxBeds", strconv.Itoa(maxBeds))
		}
	}
}

// WithBathsRange sets minimum and maximum baths filters.
func WithBathsRange(minBaths, maxBaths float64) Option {
	return func(values url.Values) {
		if minBaths > 0 {
			values.Set("minBathsTotal", strconv.FormatFloat(minBaths, 'f', -1, 64))
		}
		if maxBaths > 0 {
			values.Set("maxBathsTotal", strconv.FormatFloat(maxBaths, 'f', -1, 64))
		}
	}
}

// WithSaleAmountRange sets minimum and maximum sale amount filters.
func WithSaleAmountRange(minAmt, maxAmt float64) Option {
	return func(values url.Values) {
		if minAmt > 0 {
			values.Set("minSaleAmt", strconv.FormatFloat(minAmt, 'f', -1, 64))
		}
		if maxAmt > 0 {
			values.Set("maxSaleAmt", strconv.FormatFloat(maxAmt, 'f', -1, 64))
		}
	}
}

// WithUniversalSizeRange filters by the universal size in square feet.
func WithUniversalSizeRange(minSize, maxSize int) Option {
	return func(values url.Values) {
		if minSize > 0 {
			values.Set("minUniversalSize", strconv.Itoa(minSize))
		}
		if maxSize > 0 {
			values.Set("maxUniversalSize", strconv.Itoa(maxSize))
		}
	}
}

// WithYearBuiltRange filters by year built range.
func WithYearBuiltRange(minYear, maxYear int) Option {
	return func(values url.Values) {
		if minYear > 0 {
			values.Set("minYearBuilt", strconv.Itoa(minYear))
		}
		if maxYear > 0 {
			values.Set("maxYearBuilt", strconv.Itoa(maxYear))
		}
	}
}

// WithLotSize1Range filters by lot size in acres.
func WithLotSize1Range(minSize, maxSize float64) Option {
	return func(values url.Values) {
		if minSize > 0 {
			values.Set("minLotSize1", strconv.FormatFloat(minSize, 'f', -1, 64))
		}
		if maxSize > 0 {
			values.Set("maxLotSize1", strconv.FormatFloat(maxSize, 'f', -1, 64))
		}
	}
}

// WithLotSize2Range filters by lot size in square feet.
func WithLotSize2Range(minSize, maxSize int) Option {
	return func(values url.Values) {
		if minSize > 0 {
			values.Set("minLotSize2", strconv.Itoa(minSize))
		}
		if maxSize > 0 {
			values.Set("maxLotSize2", strconv.Itoa(maxSize))
		}
	}
}

// WithDateRange sets a start and end date for parameters with the provided prefix.
// The ATTOM Property API accepts dates formatted as YYYY/MM/DD for most filters.
func WithDateRange(prefix string, start, end time.Time) Option {
	return func(values url.Values) {
		layout := "2006/01/02"
		if !start.IsZero() {
			values.Set("start"+prefix, start.Format(layout))
		}
		if !end.IsZero() {
			values.Set("end"+prefix, end.Format(layout))
		}
	}
}

// WithISODateRange uses ISO8601 format (YYYY-MM-DD) for start/end parameters.
func WithISODateRange(prefix string, start, end time.Time) Option {
	return func(values url.Values) {
		layout := "2006-01-02"
		if !start.IsZero() {
			values.Set("start"+prefix, start.Format(layout))
		}
		if !end.IsZero() {
			values.Set("end"+prefix, end.Format(layout))
		}
	}
}

// WithPage sets the page index for paginated responses.
func WithPage(page int) Option {
	return func(values url.Values) {
		if page > 0 {
			values.Set("page", strconv.Itoa(page))
		}
	}
}

// WithPageSize sets the pagesize parameter when greater than zero.
func WithPageSize(p int) Option {
	return func(values url.Values) {
		if p > 0 {
			values.Set("pagesize", strconv.Itoa(p))
		}
	}
}

// WithOrderBy sets the orderby parameter.
func WithOrderBy(field string) Option {
	return WithString("orderby", field)
}

// WithAdditionalParam allows callers to supply custom string parameters.
func WithAdditionalParam(key, value string) Option {
	return WithString(key, value)
}

// WithWKTString sets the WKTString parameter.
func WithWKTString(wktString string) Option {
	return WithString("WKTString", wktString)
}

// WithStateID sets the StateId parameter.
func WithStateID(stateID string) Option {
	return WithString("StateId", stateID)
}

// WithFIPS sets the fips parameter.
func WithFIPS(fips string) Option {
	return WithString("fips", fips)
}

// WithAPN sets the APN parameter.
func WithAPN(apn string) Option {
	return WithString("APN", apn)
}
