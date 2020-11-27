package util

import "encoding/base64"

// VerifyDecodedPassword verifies password.
func VerifyDecodedPassword(decodedPasswd string) (string, error) {
	if decodedPasswd == "" {
		return "", nil
	}

	dec, err := base64.StdEncoding.DecodeString(decodedPasswd)
	if err != nil {
		return "", err
	}

	return string(dec), nil
}
