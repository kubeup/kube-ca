// +build !appengine

package api

import (
	"github.com/rs/xhandler"
)

func NewChain() xhandler.Chain {
	return xhandler.Chain{}
}
