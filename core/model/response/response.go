package response

import (
	"github.com/cyrildever/treee/core/exception"
)

// GetCodeFromError ...
func GetCodeFromError(err error) int {
	switch err.(type) {
	case *exception.AlreadyExistsInIndexError:
		return 303
	case *exception.EmptyItemError:
		return 400
	case *exception.InvalidHashStringError:
		return 400
	case *exception.LoopError:
		return 500
	case *exception.NotFoundError:
		return 404
	default:
		return 500
	}
}
