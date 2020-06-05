package logger

import (
	"os"
	"path/filepath"

	"github.com/inconshreveable/log15"
)

// Init ...
func Init(serviceName, contextName string) log15.Logger {
	logger := log15.New("service", serviceName, "context", contextName)
	return setHandlers(logger)
}

// InitHandler ...
func InitHandler(serviceName, contextName, requestID string) log15.Logger {
	logger := log15.New("service", serviceName, "context", contextName, "request_id", requestID)
	return setHandlers(logger)
}

func setHandlers(logger log15.Logger) log15.Logger {
	var handlers []log15.Handler

	// stdOut
	stdOutHandler := log15.CallerStackHandler("%+v", log15.StdoutHandler)
	handlers = append(handlers, stdOutHandler)

	// file
	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	fileHandler := log15.Must.FileHandler(pwd+"/treee.log", log15.LogfmtFormat())
	handlers = append(handlers, fileHandler)

	logger.SetHandler(log15.MultiHandler(handlers...))
	return logger
}
