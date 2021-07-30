package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	cql "github.com/ahsanbarkati/crankdb/cql"
	client "github.com/ahsanbarkati/crankdb/server"
	"google.golang.org/grpc"
)

type GoCrank struct {
	_conn   *grpc.ClientConn
	_client client.CrankDBClient
}

func NewCrankConnection(hostport string) (*GoCrank, error) {
	conn, err := grpc.Dial(hostport, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := client.NewCrankDBClient(conn)
	newObj := &GoCrank{_conn: conn, _client: client}
	return newObj, nil
}

func (gc *GoCrank) Set(key string, value interface{}) (*cql.SetCommandResponse, error) {

	var valType cql.DataType

	dataPacket := cql.DataPacket{
		Key: key,
	}

	switch value := value.(type) {

	case string:
		valType = cql.DataType_STRING
		dataPacket.DataType = valType
		dataPacket.StringVal = value
	case int32:
		valType = cql.DataType_INT
		dataPacket.DataType = valType
		dataPacket.S32IntVal = value
	case int64:
		valType = cql.DataType_LONG
		dataPacket.DataType = valType
		dataPacket.S64IntVal = value
	case float32:
		valType = cql.DataType_FLOAT
		dataPacket.DataType = valType
		dataPacket.FloatVal = value
	case float64:
		valType = cql.DataType_DOUBLE
		dataPacket.DataType = valType
		dataPacket.DoubleVal = value
	case bool:
		valType = cql.DataType_BOOL
		dataPacket.DataType = valType
		dataPacket.BoolVal = value
	case []byte:
		valType = cql.DataType_BYTES
		dataPacket.DataType = valType
		dataPacket.BytesVal = value
	default:
		jsonVal, err := json.Marshal(value)
		if err != nil || string(jsonVal) == "null" {
			errorMsg := fmt.Sprintf("unsupported type value - %T", value)
			return nil, errors.New(errorMsg)
		}

		valType = cql.DataType_JSON
		dataPacket.DataType = valType
		dataPacket.JsonVal = jsonVal
	}

	response, err := gc._client.Set(context.Background(), &dataPacket)
	return response, err
}

func (gc *GoCrank) Get(key string) (*cql.DataPacket, error) {
	query := cql.GetCommandRequest{Key: key}
	response, err := gc._client.Get(context.Background(), &query)
	if err != nil {
		return nil, err
	}

	return response, err
}

type FindDataPacket struct {
	key   string
	value interface{}
}

func (gc *GoCrank) Find(query map[string]interface{}) ([]FindDataPacket, error) {
	byteData, err := json.Marshal(query)
	if err != nil {
		return make([]FindDataPacket, 0), err
	}
	queryObj := cql.FindCommandRequest{Query: byteData}
	response, err := gc._client.Find(context.Background(), &queryObj)
	if err != nil {
		return make([]FindDataPacket, 0), err
	}

	var result []FindDataPacket

	for {
		dataPacket, err := response.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return make([]FindDataPacket, 0), err
		}
		var jsonVal interface{}
		json.Unmarshal(dataPacket.GetJsonVal(), &jsonVal)
		data := FindDataPacket{
			key:   dataPacket.GetKey(),
			value: jsonVal,
		}
		result = append(result, data)

	}

	return result, nil
}
