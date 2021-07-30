package server

import (
	context "context"
	"encoding/json"
	"log"

	cql "github.com/ahsanbarkati/crankdb/cql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return &cql.SetCommandResponse{Success: false}, status.Error(codes.InvalidArgument, "no value passed")
	}

	log.Printf("key: %v , valType: %v , value: %v", key, valueType, value)
	Db.Add(key, value, valueType)

	return &cql.SetCommandResponse{Success: true}, nil
}

func (s *CrankServer) Get(ctx context.Context, request *cql.GetCommandRequest) (*cql.DataPacket, error) {
	key := request.Key

	obj, ok := Db.Retrieve(key)

	if !ok {
		return &cql.DataPacket{}, status.Error(codes.NotFound, "key not found")
	}

	objType := obj.valType
	objValue := obj.value

	response := &cql.DataPacket{DataType: objType}

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
