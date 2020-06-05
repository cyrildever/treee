package api

import (
	"github.com/cyrildever/treee/api/handlers"
	routing "github.com/qiangxue/fasthttp-routing"
)

// Routes defines the routes under the /api group
func Routes(router *routing.Router) {
	apiRouter := router.Group("/api")
	apiRouter.Options("*", setCorsHeader)
	(*apiRouter).Get("/leaf", setCorsHeader, handlers.GetLeaf)
	(*apiRouter).Post("/leaf", setCorsHeader, handlers.PostLeaf)
}

func setCorsHeader(request *routing.Context) error {
	request.Response.Header.Set("Access-Control-Allow-Origin", "*")
	request.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	request.Response.Header.Set("Access-Control-Allow-Headers", "*")
	return nil
}
