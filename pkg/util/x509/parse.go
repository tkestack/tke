package x509

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func ParseBase64Encoded(encoded string) (*x509.Certificate, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	decode, _ := pem.Decode(decodedBytes)
	if decode == nil {
		return nil, fmt.Errorf("empty decode object")
	}
	certificate, err := x509.ParseCertificate(decode.Bytes)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}
