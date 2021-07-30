package server

import (
	"encoding/json"
	"reflect"

	"github.com/ahsanbarkati/crankdb/cql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func exists(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	switch queryValue := queryValue.(type) {
	case bool:
		if queryValue {
			return objectValue != nil, nil
		}

		return objectValue == nil, nil
	default:
		return false, status.Error(codes.InvalidArgument, "exists equality operator supports only boolean.")
	}
}

func eq(key string, queryValue interface{}, objectValue interface{}) (bool, error) {
	switch objectValue := objectValue.(type) {
	case float64, string, bool, nil:
		return objectValue == queryValue, nil
	default:
		return false, status.Error(codes.InvalidArgument, "eq/eq equality operator supports only numerical, string, boolean and nil values.")
	}
}

func neq(key string, queryValue interface{}, objectValue interface{}) (bool, error) {
	ok, err := eq(key, queryValue, objectValue)

	return !ok, err
}

func lt(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	if reflect.TypeOf(objectValue).Kind() != reflect.Float64 {
		return false, nil
	}

	switch queryValue := queryValue.(type) {
	case float64:
		return objectValue.(float64)-queryValue < 0, nil
	default:
		return false, status.Error(codes.InvalidArgument, "lt operator only supports numerical values.")
	}
}

func lte(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	if reflect.TypeOf(objectValue).Kind() != reflect.Float64 {
		return false, nil
	}

	switch queryValue := queryValue.(type) {
	case float64:
		return objectValue.(float64)-queryValue <= 0, nil
	default:
		return false, status.Error(codes.InvalidArgument, "lte operator only supports numerical values.")
	}
}

func gt(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	if reflect.TypeOf(objectValue).Kind() != reflect.Float64 {
		return false, nil
	}

	switch queryValue := queryValue.(type) {
	case float64:
		return objectValue.(float64)-queryValue > 0, nil
	default:
		return false, status.Error(codes.InvalidArgument, "gt operator only supports numerical values.")
	}
}

func gte(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	if reflect.TypeOf(objectValue).Kind() != reflect.Float64 {
		return false, nil
	}

	switch queryValue := queryValue.(type) {
	case float64:
		return objectValue.(float64)-queryValue >= 0, nil
	default:
		return false, status.Error(codes.InvalidArgument, "gte operator only supports numerical values.")
	}
}

func in(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	queryMap := queryValue.(map[interface{}]struct{})

	_, found := queryMap[objectValue]

	if found {
		return true, nil
	}

	return false, nil

}

func nin(key string, queryValue interface{}, objectValue interface{}) (bool, error) {

	ok, err := in(key, queryValue, objectValue)

	return !ok, err

}

func handleOperatorQuery(key string, queryObj map[string]interface{}, objectValue interface{}) (bool, error) {

	var ok bool = true
	var err error = nil

	for operator, args := range queryObj {

		switch operator {
		case "eq":
			ok, err = eq(key, args, objectValue)
		case "neq":
			ok, err = neq(key, args, objectValue)
		case "lt":
			ok, err = lt(key, args, objectValue)
		case "lte":
			ok, err = lte(key, args, objectValue)
		case "gt":
			ok, err = gt(key, args, objectValue)
		case "gte":
			ok, err = gte(key, args, objectValue)
		case "exists":
			ok, err = exists(key, args, objectValue)
		case "in":
			ok, err = in(key, args, objectValue)
		case "nin":
			ok, err = nin(key, args, objectValue)
		}

		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func optimiseQuery(queryObj *map[string]interface{}) (*map[string]interface{}, error) {
	/* Changes lists to sets (maps with 0 byte structs) */

	for _, value := range *queryObj {

		switch value := value.(type) {
		case float64, string, bool, nil:
			continue
		case map[string]interface{}:
			// update "in" operator list to map[interface{}]interface{}
			arrList, found := value["in"]
			if !found {
				continue
			}
			arrMap := make(map[interface{}]struct{})
			exists := struct{}{}
			for _, elem := range arrList.([]interface{}) {
				arrMap[elem] = exists
			}
			value["in"] = arrMap
		default:
			return nil, status.Errorf(codes.InvalidArgument, "Unsupported query value - %v", value)
		}
	}

	return queryObj, nil
}

func searchStage(query interface{}, resultStream chan *dbObject, done chan bool, errorStream chan *error) {
	queryObj := query.(map[string]interface{})

	_, err := optimiseQuery(&queryObj)

	if err != nil {
		errorStream <- &err
		close(resultStream)
		close(errorStream)
		close(done)
	}

	// Collection scan
	for _, dbobject := range Db.store {

		// ignore non JSON values.
		if dbobject.valType != cql.DataType_JSON {
			continue
		}

		dbObjValue, er := dbobject.value.(map[string]interface{})

		// Ignoring array based JSON values.
		if !er {
			continue
		}

		ok := true
		var err error = nil

		for qkey, qvalue := range queryObj {

			objValue := dbObjValue[qkey]

			switch qvalue := qvalue.(type) {
			case float64, string, bool, nil:
				ok, err = eq(qkey, qvalue, objValue)
			case map[string]interface{}:
				ok, err = handleOperatorQuery(qkey, qvalue, objValue)
			default:
				ok, err = false, status.Errorf(codes.InvalidArgument, "Unsupported query value - %v", qvalue)
			}

			if !ok {
				break
			}
		}

		if err != nil {
			errorStream <- &err
			close(resultStream)
			close(errorStream)
			close(done)
			return
		}

		if ok {
			resultStream <- dbobject
		}
	}
	done <- true
	close(resultStream)
}

func (s *CrankServer) Find(request *cql.FindCommandRequest, stream CrankDB_FindServer) error {
	/* Method call for Find command. Returns streaming output to client. */

	var queryObj interface{}
	json.Unmarshal(request.GetQuery(), &queryObj)

	resultStream := make(chan *dbObject)
	done := make(chan bool)
	errorStream := make(chan *error)

	go searchStage(queryObj, resultStream, done, errorStream)

	for {
		select {
		case result := <-resultStream:
			jsonVal, _ := json.Marshal(result.value)
			newPacket := &cql.DataPacket{Key: result.key, JsonVal: jsonVal, DataType: cql.DataType_JSON}
			if err := stream.Send(newPacket); err != nil {
				return err
			}
		case err := <-errorStream:
			return *err
		case <-done:
			return nil
		}

	}

}
