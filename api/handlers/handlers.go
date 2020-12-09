package handlers

import (
	"encoding/json"
	"time"

	apiErrors "github.com/cyrildever/treee/api/errors"
	"github.com/cyrildever/treee/common/http_errors"
	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/core/model"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/xeipuuv/gojsonschema"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func checkRequestSchema(body []byte, jsonSchema string) error {
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	requestLoader := gojsonschema.NewStringLoader(string(body))

	result, err := gojsonschema.Validate(schemaLoader, requestLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	return apiErrors.NewWrongSchemaError(result.Errors())
}

func createContext() (context.Context, context.CancelFunc, string, error) {
	requestID, err := model.GenerateUUID()
	if err != nil {
		return nil, nil, "", err
	}
	requestIDStr, _ := requestID.String()
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("requestID", requestIDStr))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // TODO Put in Config
	return ctx, cancel, requestIDStr, nil
}

func sendResponse(context string, request *routing.Context, requestID string, response interface{}, resErr error, statusCode ...int) error {
	log := logger.InitHandler("handlers", context, requestID)
	if resErr != nil {
		log.Error("Response issue", "error", resErr)
		return http_errors.SetRPCError(resErr, request, requestID)
	}
	var res []byte
	if response != nil {
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Error("Marshalling response issue", "error", err)
			return http_errors.SetMarshallingError(request, requestID)
		}
		res = jsonResponse
	} else {
		res = nil
	}
	request.Response.Header.Set("Content-Type", "application/json")
	sc := 200
	if len(statusCode) == 1 {
		sc = statusCode[0]
	}
	request.Response.SetStatusCode(sc)
	request.Response.SetBody(res)
	return nil
}
