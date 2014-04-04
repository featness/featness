package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func readHttpBody(response *http.Response) string {

	bodyBuffer := make([]byte, 1000)
	var str string

	count, err := response.Body.Read(bodyBuffer)

	for ; count > 0; count, err = response.Body.Read(bodyBuffer) {
		if err != nil {
		}
		str += string(bodyBuffer[:count])
	}

	return str
}

func getUncachedResponse(uri string) (*http.Response, error) {
	request, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Cache-Control", "no-cache")

	client := new(http.Client)

	return client.Do(request)
}

func getMe(token string) (map[string]interface{}, error) {
	response, err := getUncachedResponse("https://graph.facebook.com/me?access_token=" + token)

	if err != nil {
		return nil, err
	}

	var jsonBlob interface{}

	responseBody := readHttpBody(response)

	if responseBody != "" {
		err = json.Unmarshal([]byte(responseBody), &jsonBlob)

		if err != nil {
			return nil, err
		}

		jsonObj := jsonBlob.(map[string]interface{})
		return jsonObj, nil
	}

	return nil, fmt.Errorf("Empty response body received from Facebook.")
}

func FacebookAuthenticationProvider(token string, account string) (string, error) {
	facebookData, err := getMe(token)

	if err != nil {
		return "", fmt.Errorf("access token was invalid: %v.", err)
	}

	if facebookData["username"] != account {
		return "", fmt.Errorf("access token is for a different account.")
	}

	return token, nil
}

func AuthenticateWithFacebook(w http.ResponseWriter, r *http.Request) {
	AuthenticationRoute(w, r, "Facebook", FacebookAuthenticationProvider)
}
