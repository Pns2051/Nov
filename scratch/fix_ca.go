package main
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)
func main() {
	certPath, keyPath := "/home/ela/.adblock-proxy/ca-cert.pem", "/home/ela/.adblock-proxy/ca-key.pem"
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	spki, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	skid := sha1.Sum(spki)
	subj := pkix.Name{Organization: []string{"System AdBlocker CA"}, CommonName: "System AdBlocker Root"}
	template := x509.Certificate{
		SerialNumber: serial, Subject: subj, Issuer: subj,
		NotBefore: time.Now().Add(-24 * time.Hour), NotAfter: time.Now().Add(10 * 365 * 24 * time.Hour),
		IsCA: true, BasicConstraintsValid: true,
		KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SubjectKeyId: skid[:], AuthorityKeyId: skid[:], // Points to itself as root
	}
	der, _ := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	cOut, _ := os.Create(certPath)
	pem.Encode(cOut, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	kOut, _ := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(kOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
}
