package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) (string, error) {
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func GenerateShortLink(initialLink string, userID string) string {
	return GenerateShortLinkWithSalt(initialLink, userID, 0)
}

func GenerateShortLinkWithSalt(initialLink string, userID string, salt int) string {
	input := initialLink + userID
	if salt > 0 {
		input = fmt.Sprintf("%s:%d", input, salt)
	}

	urlHashBytes := sha256Of(input)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString, err := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	if err != nil {
		return ""
	}
	if len(finalString) < 8 {
		return finalString
	}
	return finalString[:8]
}
