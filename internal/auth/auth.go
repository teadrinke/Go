package auth

import ("net/http"
	"errors"
	"strings"
)


//GetAPIKey extracts the API KEY from the headers of an http request
//Example: 
//Authorization : ApiKey {insert apikey here}
func GetAPIKey(headers http.Header) (string, error){
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2{
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first path of auth header")
	}
	return vals[1], nil
}