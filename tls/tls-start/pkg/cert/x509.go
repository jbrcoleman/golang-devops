package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
	"time"

	"github.com/jbrcoleman/golang-devops/tls/tls-start/pkg/key"
)

func CreateCACert(ca *CACert, keyFilePath, caCertFilePath string) error {
	template := &x509.Certificate{
		SerialNumber: ca.Serial,
		Subject: pkix.Name{
			Country:            RemoveEmptyString([]string{ca.Subject.Country}),
			Organization:       RemoveEmptyString([]string{ca.Subject.Organization}),
			OrganizationalUnit: RemoveEmptyString([]string{ca.Subject.OrganizationalUnit}),
			Locality:           RemoveEmptyString([]string{ca.Subject.Locality}),
			Province:           RemoveEmptyString([]string{ca.Subject.Province}),
			StreetAddress:      RemoveEmptyString([]string{ca.Subject.StreetAddress}),
			CommonName:         ca.Subject.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(ca.ValidForYears, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	keyBytes, certBytes, err := createCert(template, nil, nil)
	if err != nil {
		return err
	}

	if err := os.WriteFile(keyFilePath, keyBytes, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(caCertFilePath, certBytes, 0644); err != nil {
		return err
	}

	return nil
}

func CreateCert(cert *Cert, caKey []byte, caCert []byte, keyFilePath, certFilePath string) error {
	template := &x509.Certificate{
		SerialNumber: cert.Serial,
		Subject: pkix.Name{
			Country:            RemoveEmptyString([]string{cert.Subject.Country}),
			Organization:       RemoveEmptyString([]string{cert.Subject.Organization}),
			OrganizationalUnit: RemoveEmptyString([]string{cert.Subject.OrganizationalUnit}),
			Locality:           RemoveEmptyString([]string{cert.Subject.Locality}),
			Province:           RemoveEmptyString([]string{cert.Subject.Province}),
			StreetAddress:      RemoveEmptyString([]string{cert.Subject.StreetAddress}),
			PostalCode:         RemoveEmptyString([]string{cert.Subject.PostalCode}),
			CommonName:         cert.Subject.CommonName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(cert.ValidForYears, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		DNSNames:    RemoveEmptyString(cert.DNSNames),
	}

	caKeyParsed, err := key.PrivateKeyPemToRSA(caKey)
	if err != nil {
		return err
	}

	caCertParsed, err := PemToX509(caCert)
	if err != nil {
		return err
	}

	keyBytes, certBytes, err := createCert(template, caKeyParsed, caCertParsed)
	if err != nil {
		return err
	}

	if err := os.WriteFile(keyFilePath, keyBytes, 0600); err != nil {
		return err
	}
	if err := os.WriteFile(certFilePath, certBytes, 0644); err != nil {
		return err
	}

	return nil
}

func createCert(template *x509.Certificate, caKey *rsa.PrivateKey, caCert *x509.Certificate) ([]byte, []byte, error) {
	var (
		derBytes []byte
		certOut  bytes.Buffer
		keyOut   bytes.Buffer
	)

	privateKey, err := key.CreateRSAPrivateKey(4096)
	if err != nil {
		return nil, nil, err
	}
	if template.IsCA {
		derBytes, err = x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return nil, nil, err
		}
	} else {
		derBytes, err = x509.CreateCertificate(rand.Reader, template, caCert, &privateKey.PublicKey, caKey)
		if err != nil {
			return nil, nil, err
		}
	}
	if err = pem.Encode(&certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return nil, nil, err
	}
	if err = pem.Encode(&keyOut, key.RSAPrivateKeyToPEM(privateKey)); err != nil {
		return nil, nil, err
	}
	return keyOut.Bytes(), certOut.Bytes(), nil
}

func RemoveEmptyString(input []string) []string {
	if len(input) == 1 && input[0] == "" {
		return []string{}
	}
	return input
}
