package handlers

import (
	"encoding/json"

	"github.com/cyrildever/treee/api/handlers/schema"
	"github.com/cyrildever/treee/common/http_errors"
	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/core/index"
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/model"
	"github.com/cyrildever/treee/core/model/response"
	routing "github.com/qiangxue/fasthttp-routing"
)

// GetLeaf ...
func GetLeaf(request *routing.Context) error {
	_, cancel, requestID, err := createContext()
	log := logger.InitHandler("handlers", "GetLeaf", requestID)
	if err != nil {
		log.Error("Creating context error", "error", err)
		return http_errors.SetInternalError(request, requestID)
	}
	defer cancel()

	var ids []model.Hash
	request.QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "ids" {
			ids = append(ids, model.Hash(string(value)))
		}
	})

	if len(ids) == 0 {
		log.Info("Empty query string")
		return http_errors.SetInvalidParam(request, requestID, "missing at least one leaf id")
	}

	var resp []branch.Leaf
	for _, id := range ids {
		if leaf, err := index.Current.Search(id); err == nil {
			resp = append(resp, *leaf)
		}
	}

	return sendResponse("GetLeaf", request, requestID, resp, nil)
}

// PostLeaf ...
func PostLeaf(request *routing.Context) error {
	_, cancel, requestID, err := createContext()
	log := logger.InitHandler("handlers", "PostLeaf", requestID)
	if err != nil {
		log.Error("Creating context error", "error", err)
		return http_errors.SetInternalError(request, requestID)
	}
	defer cancel()

	if err = checkRequestSchema(request.PostBody(), schema.Leaf); err != nil {
		log.Error("Wrong leaf format", "error", err)
		return http_errors.SetInvalidParam(request, requestID, err.Error())
	}

	leaf := branch.Leaf{}
	err = json.Unmarshal(request.PostBody(), &leaf)
	if err != nil {
		log.Warn("Unmarshalling error", "error", err)
		return http_errors.SetMarshallingError(request, requestID)
	}

	save := false
	var resp response.PostLeaf
	err = index.Current.Add(leaf)
	if err != nil {
		resp = response.PostLeaf{
			Code:  response.GetCodeFromError(err),
			Error: err.Error(),
		}
	} else {
		idLeaf, _ := leaf.ID.String()
		resp = response.PostLeaf{
			Code:   200,
			Result: idLeaf,
		}
		save = true
	}

	err = sendResponse("PostLeaf", request, requestID, resp, nil)
	if save {
		go index.Current.Save()
	}
	return err
}
