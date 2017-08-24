package api

import (
	"errors"
	"net/http"

	"github.com/cloudflare/cfssl/api"
	"github.com/cloudflare/cfssl/api/signhandler"
	"github.com/kubeup/kube-ca/signer"
	"golang.org/x/net/context"
)

type SignHandler struct {
	signer *signer.EncryptedSigner
}

func NewSignHandler(signer *signer.EncryptedSigner) *SignHandler {
	return &SignHandler{
		signer: signer,
	}
}

func (s *SignHandler) ServeHTTPC(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if s.signer == nil {
		_, err := signer.NewEncryptedSignerFromEnv()
		api.HandleError(w, err)
		return

		api.HandleError(w, errors.New("No signer configured"))
		return
	}
	signer, err := s.signer.Decrypt(ctx)
	if err != nil {
		api.HandleError(w, err)
		return
	}
	handler, err := signhandler.NewAuthHandlerFromSigner(signer)
	if err != nil {
		api.HandleError(w, err)
		return
	}
	handler.ServeHTTP(w, r)
}
