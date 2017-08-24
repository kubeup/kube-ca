package api

import (
	"net/http"

	"github.com/kubeup/kube-ca/signer"
	"github.com/rs/xmux"
)

func RegisterHandler(signer *signer.EncryptedSigner) {
	c := NewChain()
	mux := xmux.New()
	group := mux.NewGroup("/api/v1/cfssl")

	signHandler := NewSignHandler(signer)
	group.POST("/authsign", signHandler)

	http.Handle("/api/v1/cfssl/", c.Handler(mux))
}
