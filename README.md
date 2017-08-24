# kube-ca

kube-ca is an experimental CA based on Cloudflare's CFSSL project. It is built together with [kube-remote-signer][kube-remote-signer]
as a proof of concept about using external CA for CSR signing to enhance security for Kubernetes clusters.

Currently we have to put the CA private key and certificate on master nodes and pass them to the builtin
certificate controller running in [kube-controller-manager][kube-controller-manager] to support the token
based node bootstrapping process. It is a burden to manage the CA private key properly and there are risks
about key leaking which would leads to critical security incidents.

By moving the signer out of the Kubernetes cluster, we could reduce security risk and simplify the
configuration process for master servers.

## Features

- Support running as a standalone process or running in Google App Engine
- CA private key protection by using Google Cloud KMS for encryption and decryption
- HMAC authentication to avoid unauthorized access
- Compatible with CFSSL api server but only provide `authsign` api endpoint to lower attack surface

## Installation

For better protection, you could run `kube-ca` in Google App Engine which is secure sandbox environment.
Follow these steps to setup `kube-ca` in Google App Engine.

### Create KMS crypto key

First create a KeyRing:

```
gcloud kms keyrings create KEYRING_NAME --location LOCATION
```

Then create a CryptoKey in the KeyRing:

```
gcloud kms keys create CRYPTOKEY_NAME --location LOCATION --keyring KEYRING_NAME --purpose encryption
```

Check this [official document][kms-create-keys] for more detail.

### Generate a self-signed root CA

Create `csr.json` with the following content:

```
{
    "hosts": [
        "example.com",
        "www.example.com"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "US",
            "L": "San Francisco",
            "O": "Internet Widgets, Inc.",
            "OU": "WWW",
            "ST": "California"
        }
    ]
}
```

Create the CA private key and certificate:

```
cfssl genkey -initca csr.json | cfssljson -bare ca
```

### Encrypt CA private key with KMS

Use `gcloud` to encrypt the CA private key:

```
gcloud kms encrypt --key CRYPTOKEY_NAME --keyring KEYRING_NAME --location LOCATION --plaintext-file ca-key.pem --ciphertext-file ca-key.pem.enc
```

### Update app.yaml and deploy

Update `appengine/app.yaml` file.

- Update `PROJECT` with your Google Cloud project name
- Update `LOCATION`, `KEYRING_NAME` and `CRYPTOKEY_NAME` with the settings which is used to create the key
- Update `AUTHENTICATION_KEY` with a random hex-encoded string (e.g. "000102030405060708")
- Update `BASE64_ENCRYPTED_PRIVATE_KEY` with the output of `cat ca-key.pem.enc | base64`
- Update `BASE64_CERTIFICATE` with the output of `cat ca.pem | base64`

Deploy the application with `goapp deploy`.

## License

Apache Version 2.0

[kube-remote-signer]: https://github.com/kubeup/kube-remote-signer
[kube-controller-manager]: https://kubernetes.io/docs/admin/kube-controller-manager/
[kms-create-keys]: https://cloud.google.com/kms/docs/creating-keys
