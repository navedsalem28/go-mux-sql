package main

import (
	"encoding/json"
	"net/http"
)

// LoginAPI for application users
func LoginAPI(w http.ResponseWriter, r *http.Request) {
	Debug(r.Proto + " Request From: " + r.RemoteAddr + "\tURL: " + r.URL.Path)
	addRestHeader(w)
	ResponseJSON := make(map[string]interface{})
	ResponseJSON["message"] = "invalid API"
	ResponseJSON["status"] = "407"
	RequestJSON, err := Input(r)
	if err != nil {
		json.NewEncoder(w).Encode(RequestJSON)
		return
	}
	addRestHeader(w)

	ok, _, userID := GetUserIdentification(RequestJSON, ResponseJSON)
	if !ok {
		json.NewEncoder(w).Encode(ResponseJSON)
		return
	}

	password, ok := RequestJSON["password"]
	if !ok {
		ResponseJSON["status"] = 407
		ResponseJSON["message"] = "password required"
		json.NewEncoder(w).Encode(ResponseJSON)
		return
	}

	if internalConfig.UserName != userID || internalConfig.Password != InterfaceToString(password) {
		ResponseJSON["status"] = 401
		ResponseJSON["message"] = "Invalid Credentials"
		json.NewEncoder(w).Encode(ResponseJSON)
		return
	}

	ResponseJSON["user_id"] = 1
	ResponseJSON["full_name"] = "admin"

	ResponseJSON["status"] = 200

	ResponseJSON["message"] = "Login Successful"
	json.NewEncoder(w).Encode(ResponseJSON)
}

// GetUserIdentification identify user identification field's name & value
func GetUserIdentification(RequestJSON map[string]interface{}, ResponseJSON map[string]interface{}) (bool, string, string) {

	userid, ok := RequestJSON["user_name"]
	if ok {
		return true, "`username`", InterfaceToString(userid)
	}
	if ResponseJSON["status"] == 405 {
		return false, "", ""
	}
	userid, ok = RequestJSON["phone_num"]
	if ok {
		return true, "`phone_num`", InterfaceToString(userid)
	}
	if ResponseJSON["status"] == 405 {
		return false, "", ""
	}
	userid, ok = RequestJSON["email"]
	if ok {
		return true, "`email`", InterfaceToString(userid)
	}
	if ResponseJSON["status"] != 405 {
		ResponseJSON["status"] = 406
		ResponseJSON["message"] = "No user identification available"
	}
	return false, "", ""
}
