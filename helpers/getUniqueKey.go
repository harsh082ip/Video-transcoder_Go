package helpers

import (
	"crypto/rand"
	"math/big"

	"github.com/harsh082ip/Video-transcoder_Go/consts"
)

func GetUniqueKey(length int) (string, error) {
	var sessionID string

	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(consts.LetterBytes))))
		if err != nil {
			return "", err
		}
		result[i] = consts.LetterBytes[num.Int64()]
	}
	sessionID = string(result)

	return sessionID, nil
}
