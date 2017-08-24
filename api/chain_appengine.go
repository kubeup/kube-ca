// +build appengine

package api

import (
	"net/http"
	"time"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

func NewChain() xhandler.Chain {
	c := xhandler.Chain{}
	c.UseC(AppengineHandler)
	return c
}

func AppengineHandler(next xhandler.HandlerC) xhandler.HandlerC {
	return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = appengine.NewContext(r)
		ctx, _ = context.WithTimeout(ctx, time.Duration(30)*time.Second)
		next.ServeHTTPC(ctx, w, r)
	})
}
