package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	certPath := "/home/ela/.adblock-proxy/ca-cert.pem"
	keyPath := "/home/ela/.adblock-proxy/ca-key.pem"

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	notBefore := time.Now().Add(-24 * time.Hour)
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	subject := pkix.Name{
		Organization: []string{"System AdBlocker CA"},
		CommonName:   "System AdBlocker Root",
	}

	spki, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	skid := sha1.Sum(spki)

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		Issuer:                subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		BasicConstraintsValid: true,
		MaxPathLenZero:        false,
		MaxPathLen:            2,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SubjectKeyId:          skid[:],
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	certOut, _ := os.Create(certPath)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()

	fmt.Println("Perfect CA generated.")
}
