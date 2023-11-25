package helper

import (
	"net/url"
	"strings"
)

func GetQueryParameters(queryString string) map[string]string {
	queryString = strings.ToLower(queryString)
	values, _ := url.ParseQuery(queryString)
	parameters := make(map[string]string)

	for key, val := range values {
		parameters[key] = strings.Join(val, ",")
	}

	return parameters
}
