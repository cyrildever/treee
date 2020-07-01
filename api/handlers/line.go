package handlers

import (
	"github.com/cyrildever/treee/common/http_errors"
	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/core/index"
	"github.com/cyrildever/treee/core/model"
	routing "github.com/qiangxue/fasthttp-routing"
)

// GetLine ...
func GetLine(request *routing.Context) error {
	_, cancel, requestID, err := createContext()
	log := logger.InitHandler("handlers", "GetLine", requestID)
	if err != nil {
		log.Error("Creating context error", "error", err)
		return http_errors.SetInternalError(request, requestID)
	}
	defer cancel()

	var id model.Hash
	request.QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "id" {
			id = model.Hash(string(value))
		}
	})

	if id.IsEmpty() {
		log.Info("Empty query string")
		return http_errors.SetInvalidParam(request, requestID, "missing the leaf id")
	}

	line, err := index.Current.Line(id)
	if err != nil {
		log.Error("Impossible to find a line", "error", err)
		return http_errors.SetInternalError(request, requestID)
	}
	res := make([]string, len(line))
	for i, leaf := range line {
		idStr, _ := leaf.ID.String()
		res[i] = idStr
	}
	if len(res) == 0 {
		return http_errors.SetNotFoundError(request, requestID)
	}

	return sendResponse("GetLine", request, requestID, res, nil)
}
