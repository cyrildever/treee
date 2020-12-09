package handlers

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/cyrildever/treee/api/handlers/schema"
	"github.com/cyrildever/treee/common/http_errors"
	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/core/exception"
	"github.com/cyrildever/treee/core/index"
	"github.com/cyrildever/treee/core/index/branch"
	"github.com/cyrildever/treee/core/index/search"
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
	takeLast := request.QueryArgs().GetBool("takeLast")

	if len(ids) == 0 {
		log.Info("Empty query string")
		return http_errors.SetInvalidParam(request, requestID, "missing at least one leaf id")
	}

	var res []branch.Leaf

	if len(ids) == 1 {
		if takeLast {
			if leaf, err := index.Current.Last(ids[0]); err == nil {
				res = append(res, *leaf)
			}
		} else {
			if leaf, err := index.Current.Search(ids[0]); err == nil {
				res = append(res, *leaf)
			}
		}
	} else {
		var engine search.Engine
		if takeLast {
			engine = index.Current.Last
		} else {
			engine = index.Current.Search
		}
		var resp = make(chan branch.Leaf, len(ids))
		var wg sync.WaitGroup
		for _, id := range ids {
			wg.Add(1)
			go func(wg *sync.WaitGroup, engine search.Engine, id model.Hash, resp chan branch.Leaf) {
				defer wg.Done()
				if leaf, err := engine(id); err == nil {
					resp <- *leaf
				} else {
					resp <- branch.Leaf{}
				}
			}(&wg, engine, id, resp)
		}
		wg.Wait()
		close(resp)
		for r := range resp {
			if r.ID.NonEmpty() {
				res = append(res, r)
			}
		}
	}

	if len(res) == 0 {
		return http_errors.SetNotFoundError(request, requestID)
	}

	return sendResponse("GetLeaf", request, requestID, res, nil)
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
	log.Debug("Receiving leaf...", "id", leaf.ID) // Remove in production

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

func DeleteLeaf(request *routing.Context) error {
	_, cancel, requestID, err := createContext()
	log := logger.InitHandler("handlers", "DeleteLeaf", requestID)
	if err != nil {
		log.Error("Creating context error", "error", err)
		return http_errors.SetInternalError(request, requestID)
	}
	defer cancel()

	var ids model.Hashes
	request.QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "ids" {
			if string(value) != "" {
				ids = append(ids, model.Hash(string(value)))
			}
		}
	})

	if len(ids) == 0 {
		log.Info("Empty query string")
		return http_errors.SetInvalidParam(request, requestID, "missing at least one leaf id")
	}

	var errs []string
	for _, id := range ids {
		if err := index.Current.Remove(id); err != nil {
			if _, ok := err.(*exception.NotFoundError); !ok {
				idStr, _ := id.String()
				errs = append(errs, idStr)
			}
		}
	}

	if len(errs) > 0 {
		log.Info("Unable to remove all passed leaves", "ids", strings.Join(errs, ","))
		res := response.UndeletedResponse{
			List: errs,
		}
		return sendResponse("DeleteLeaf", request, requestID, res, nil)
	}

	return sendResponse("DeleteLeaf", request, requestID, nil, nil, 204)
}
