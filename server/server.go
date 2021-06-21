package server

import (
	context "context"
	"encoding/json"
	"errors"

	cql "github.com/shreybatra/crankdb/cql"
)

type CrankServer struct {
	UnimplementedCrankDBServer
}

var Db *Database = NewDatabase()

func (s *CrankServer) Set(ctx context.Context, request *cql.DataPacket) (*cql.SetCommandResponse, error) {
	key := request.Key
	valueType := request.GetDataType()

	var value interface{}

	switch valueType {
	case cql.DataType_BOOL:
		value = request.GetBoolVal()
	case cql.DataType_BYTES:
		value = request.GetBytesVal()
	case cql.DataType_INT:
		value = request.GetS32IntVal()
	case cql.DataType_LONG:
		value = request.GetS64IntVal()
	case cql.DataType_FLOAT:
		value = request.GetFloatVal()
	case cql.DataType_DOUBLE:
		value = request.GetDoubleVal()
	case cql.DataType_STRING:
		value = request.GetStringVal()
	case cql.DataType_JSON:
		json.Unmarshal(request.GetJsonVal(), &value)
	default:
		return &cql.SetCommandResponse{Success: false}, errors.New("no value passed")
	}

	Db.Add(key, value, valueType)

	return &cql.SetCommandResponse{Success: true}, nil
}

func (s *CrankServer) Get(ctx context.Context, request *cql.GetCommandRequest) (*cql.DataPacket, error) {
	key := request.Key

	obj, ok := Db.Retrieve(key)

	if !ok {
		return &cql.DataPacket{}, errors.New("key not found")
	}

	response := &cql.DataPacket{}

	objType := obj.valType
	objValue := obj.value

	switch objType {
	case cql.DataType_BOOL:
		response.BoolVal = objValue.(bool)
	case cql.DataType_BYTES:
		response.BytesVal = objValue.([]byte)
	case cql.DataType_INT:
		response.S32IntVal = objValue.(int32)
	case cql.DataType_LONG:
		response.S64IntVal = objValue.(int64)
	case cql.DataType_FLOAT:
		response.FloatVal = objValue.(float32)
	case cql.DataType_DOUBLE:
		response.DoubleVal = objValue.(float64)
	case cql.DataType_STRING:
		response.StringVal = objValue.(string)
	case cql.DataType_JSON:
		response.JsonVal, _ = json.Marshal(objValue)
	}

	return response, nil
}

func (s *CrankServer) Find(request *cql.FindCommandRequest, stream CrankDB_FindServer) error {

	var queryObj interface{}
	json.Unmarshal(request.GetQuery(), &queryObj)

	resultStream := make(chan *dbObject)

	go searchStage(queryObj, resultStream)

	for result := range resultStream {
		jsonVal, _ := json.Marshal(result.value)
		newPacket := &cql.DataPacket{Key: result.key, JsonVal: jsonVal, DataType: cql.DataType_JSON}
		if err := stream.Send(newPacket); err != nil {
			return err
		}
	}

	return nil
}

func searchStage(query interface{}, resultStream chan *dbObject) {
	queryObj := query.(map[string]interface{})

	for _, object := range Db.store {
		if object.valType != cql.DataType_JSON {
			continue
		}
		objValue := object.value.(map[string]interface{})

		ok := true
		for key, value := range queryObj {
			if objValue[key] != value {
				ok = false
				break
			}
		}

		if ok {
			resultStream <- object
		}
	}
	close(resultStream)
}
