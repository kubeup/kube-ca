package app

import (
	"github.com/kubeup/kube-ca/api"
	"github.com/kubeup/kube-ca/signer"
)

func init() {
	s, _ := signer.NewEncryptedSignerFromEnv()
	api.RegisterHandler(s)
}
