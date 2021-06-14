package server

import (
	"encoding/json"
)

func extractNextArgument(arguments string) (string, string) {

	var i int
	length := len(arguments)

	for i = range arguments {
		if arguments[i] == ' ' {
			break
		}
	}

	if i == length-1 && arguments[i] != ' ' {
		i++
	}

	extracted := string(arguments[0:i])

	var remaining string

	if i+1 <= length {
		remaining = string(arguments[i+1:])
	} else {
		remaining = ""
	}

	return extracted, remaining
}

func set(arguments string) (response string) {

	key, arguments := extractNextArgument(arguments)

	var data interface{}

	if err := json.Unmarshal([]byte(arguments), &data); err != nil {
		return "Invalid value"
	}

	database[key] = data
	return key + " set"
}

func get(arguments string) (response interface{}) {
	key, _ := extractNextArgument(arguments)
	value, ok := database[key]
	if !ok {
		return "Not found"
	}
	return value
}

func del(arguments string) (response interface{}) {
	key, _ := extractNextArgument(arguments)
	_, ok := database[key]
	if !ok {
		return "Not found"
	}
	delete(database, key)
	return key + " deleted."
}

func find(arguments string) (response []interface{}) {

	var data interface{}

	if err := json.Unmarshal([]byte(arguments), &data); err != nil {
		panic(err)
	}

	queryObj := data.(map[string]interface{})

	var result []interface{} = make([]interface{}, 0)

	for _, item := range database {
		item, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		found := true
		for key, value := range queryObj {
			if item[key] != value {
				found = false
				break
			}
		}
		if found {
			result = append(result, item)
		}
	}
	return result
}
