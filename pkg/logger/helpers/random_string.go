package helpers

import cryptoRand "crypto/rand"

const Letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := cryptoRand.Read(b)
	return b, err
}

func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)

	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = Letters[b%byte(len(Letters))]
	}

	return string(bytes), nil
}
