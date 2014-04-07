package api

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

type RequestData struct {
	cookies string
	headers string
	ip      string
}

func execute(path string, obj string) bool {
	out, err := exec.Command("bash", path, obj).Output()
	if err != nil {
		return false
	}
	result, err := strconv.ParseBool(strings.Trim(string(out), "\n"))
	if err != nil {
		return false
	}
	return result
}

func Matches(matcherPaths []string, data *RequestData) (bool, error) {
	if len(matcherPaths) == 0 {
		return true, nil
	}

	var matcherChan chan bool = make(chan bool)

	dataJson, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	for _, path := range matcherPaths {
		go func(path string, obj string) {
			matcherChan <- execute(path, obj)
		}(path, string(dataJson))
	}

	numberOfMatchersReturned := 0
	for {
		result := <-matcherChan
		if result {
			return true, nil
		}
		numberOfMatchersReturned++

		if numberOfMatchersReturned >= len(matcherPaths) {
			break
		}
	}

	return false, nil
}
