package api

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/cyrildever/treee/common/logger"
	"github.com/cyrildever/treee/config"
)

// InitHTTPServer ...
func InitHTTPServer(conf *config.Config) {
	log := logger.Init("api", "InitHTTPServer")

	router := routing.New()
	Routes(router)
	address := conf.Host + ":" + conf.HTTPPort
	log.Info("API started listening", "address", address)
	log.Crit(fasthttp.ListenAndServe(address, router.HandleRequest).Error())
}
