package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
Note:
We got to decide and design how do we want the authorization
information exist in the http header
*/
// Example:
// Authorization: ApiKey {some api key here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}

	tokens := strings.Split(val, " ")
	if len(tokens) != 2 {
		return "", errors.New("malformed auth header")
	}
	if tokens[0] != "ApiKey" {
		return "", errors.New("first part of auth header is malformed")
	}

	return tokens[1], nil
}
