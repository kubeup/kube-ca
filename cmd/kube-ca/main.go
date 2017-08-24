package main

import (
	"log"
	"net/http"

	"github.com/kubeup/kube-ca/api"
	"github.com/kubeup/kube-ca/signer"
)

func main() {
	s, _ := signer.NewEncryptedSignerFromEnv()
	api.RegisterHandler(s)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
