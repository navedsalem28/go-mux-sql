package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type Route struct {
	Name    string
	Path    string
	Method  string
	Handler http.HandlerFunc
}

var Routes = []Route{
	{
		"Home Page",
		"/",
		"GET",
		home,
	},
	{
		"Login",
		"/Login",
		"POST",
		LoginAPI,
	},
}

func Input(r *http.Request) (map[string]interface{}, error) {
	ResponseJSON := make(map[string]interface{})
	BodyText, err := ioutil.ReadAll(r.Body)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)
	if err != nil {
		ResponseJSON["message"] = "Invalid Request Body"
		ResponseJSON["status"] = 406
		Error("controlCommand: can't ready request body text, Error: " + err.Error())
		return ResponseJSON, err
	}
	Log("Request Body Text:\t" + string(BodyText))
	err = json.Unmarshal(BodyText, &ResponseJSON)
	if err != nil {
		ResponseJSON["message"] = "Invalid JSON"
		ResponseJSON["status"] = 406
		Error("controlCommand: Failed Parsing JSON: " + err.Error())
		return ResponseJSON, err
	}

	return ResponseJSON, err
}
