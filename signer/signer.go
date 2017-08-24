package signer

import (
	"crypto/x509"
	"encoding/base64"
	"errors"
	"os"

	"github.com/cloudflare/cfssl/auth"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type EncryptedSigner struct {
	cryptoKey string
	authKey   string

	key  string
	cert *x509.Certificate
}

func NewEncryptedSigner(cryptoKey, authKey, key, cert string) (*EncryptedSigner, error) {
	if cryptoKey == "" {
		return nil, errors.New("cryptoKey should not be empty")
	}
	if authKey == "" {
		return nil, errors.New("authKey should not be empty")
	}
	if key == "" {
		return nil, errors.New("key should not be empty")
	}
	if cert == "" {
		return nil, errors.New("cert should not be empty")
	}
	c, err := base64.StdEncoding.DecodeString(cert)
	if err != nil {
		return nil, err
	}
	parsedCert, err := helpers.ParseCertificatePEM(c)
	if err != nil {
		return nil, err
	}
	return &EncryptedSigner{
		cryptoKey: cryptoKey,
		authKey:   authKey,
		key:       key,
		cert:      parsedCert,
	}, nil
}

func (s *EncryptedSigner) Decrypt(ctx context.Context) (*local.Signer, error) {
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}
	kms, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}
	decryptRequest := cloudkms.DecryptRequest{Ciphertext: s.key}
	resp, err := kms.Projects.Locations.KeyRings.CryptoKeys.Decrypt(
		s.cryptoKey, &decryptRequest).Do()
	if err != nil {
		return nil, err
	}
	cakey, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		return nil, err
	}
	priv, err := helpers.ParsePrivateKeyPEMWithPassword(cakey, nil)
	if err != nil {
		return nil, err
	}
	policy := &config.Signing{
		Profiles: map[string]*config.SigningProfile{},
		Default:  config.DefaultConfig()}
	policy.Default.Provider, err = auth.New(s.authKey, nil)
	if err != nil {
		return nil, err
	}
	return local.NewSigner(priv, s.cert, signer.DefaultSigAlgo(priv), policy)
}

func NewEncryptedSignerFromEnv() (*EncryptedSigner, error) {
	cryptoKey := os.Getenv("KUBE_CA_CRYPTOKEY")
	authKey := os.Getenv("KUBE_CA_AUTHKEY")
	key := os.Getenv("KUBE_CA_KEY")
	cert := os.Getenv("KUBE_CA_CERT")

	return NewEncryptedSigner(cryptoKey, authKey, key, cert)
}
