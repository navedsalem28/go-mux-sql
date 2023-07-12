package main

import (
	"fmt"
	"log"
	"strconv"
)

var Config map[string]string
var Countries []map[string]string

// LoadConfiguration API
func LoadConfiguration() bool {
	log.Println("********Loading Configuration*****")
	params := make([]interface{}, 0)
	//params[0] = 1
	ConfigData, ok := GetAllRows("SELECT * FROM `configurations`", params, "default")
	if !ok {
		panic("Can't load configuration from database")
		return false
	}
	//for _, ThisParam := range ConfigData {
	for i := 0; i < len(ConfigData); i++ {
		ThisParam := ConfigData[i]
		log.Println(fmt.Sprintf("This Parameter Name: %s\tValue: %s", ThisParam["param"], ThisParam["value"]))
		Config[ThisParam["param"]] = ThisParam["value"]
	}

	log.Println(fmt.Sprintf("Configurations : %v", Config))
	return true
}

// GetKeyFromJSON (RequestJSON map[string]interface{}, ResJSON map[string]interface{}, KeyName string) try to get value of key from json and return (bool status, string value)
func GetKeyFromJSON(RequestJSON map[string]interface{}, ResJSON map[string]interface{}, KeyName string) (bool, string) {
	value, ok := RequestJSON[KeyName]
	if !ok {
		ResJSON["status"] = 407
		ResJSON["message"] = KeyName + " required"
		return false, ""
	}
	KeyValue := InterfaceToString(value)
	if IsSQLSafe(KeyValue) {
		return true, KeyValue
	}
	ResJSON["status"] = 405
	ResJSON["message"] = "prohibited chracters in key " + KeyName
	return false, ""
}

// InterfaceToString (interface{}) get value of an interface return as string {
func InterfaceToString(ThisInterface interface{}) string {
	if ThisInterface == nil {
		log.Println("InterfaceToString :: NIL value passed for conversion")
		return ""
	}
	switch ThisInterface.(type) {
	case int:
		return strconv.Itoa(ThisInterface.(int))
	case int64:
		return strconv.FormatInt(ThisInterface.(int64), 10)
	case float64:
		TmpStr := strconv.FormatFloat(ThisInterface.(float64), 'f', 10, 64)
		if TmpStr[(len(TmpStr)-10):] == "0000000000" {
			return TmpStr[:(len(TmpStr) - 11)]
		}
		return TmpStr
	default:
		return ThisInterface.(string)
	}
}

// isNumeric (string ) check if provided string is numeric
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// InterfaceToInt (interface{}) get value of interface and return as int64
func InterfaceToInt(ThisInterface interface{}) int64 {
	switch ThisInterface.(type) {
	case int:
		return ThisInterface.(int64)
	case int64:
		return ThisInterface.(int64)
	case float64:
		//return ThisInterface.(int64)
		return int64(ThisInterface.(float64))
	default:
		val, _ := strconv.ParseInt(ThisInterface.(string), 10, 64)
		return val
	}
}
