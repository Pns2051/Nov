package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/elazarl/goproxy"
)

func LoadOrCreateCA(certPath, keyPath string) (*tls.Certificate, error) {
	certReq, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err == nil {
		return &certReq, nil
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	spki, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	skid := sha1.Sum(spki)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"System AdBlocker CA"},
			CommonName:   "System AdBlocker Root",
		},
		Issuer: pkix.Name{
			Organization: []string{"System AdBlocker CA"},
			CommonName:   "System AdBlocker Root",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLenZero:        false,
		MaxPathLen:            2,
		SubjectKeyId:          skid[:],
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s for writing: %w", certPath, err)
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s for writing: %w", keyPath, err)
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	certReq, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load generated key pair: %w", err)
	}
	return &certReq, nil
}

func SetCAForGoproxy(proxy *goproxy.ProxyHttpServer, cert *tls.Certificate) {
	goproxy.GoproxyCa = *cert
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
}
