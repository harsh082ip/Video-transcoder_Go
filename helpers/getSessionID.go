package helpers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/harsh082ip/Video-transcoder_Go/consts"
	"github.com/harsh082ip/Video-transcoder_Go/db"
)

func CreateSeessionID(length int, ctx context.Context) (string, error) {
	var sessionID string
	l := length
	for {
		result := make([]byte, length)
		for i := range result {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(consts.LetterBytes))))
			if err != nil {
				return "", err
			}
			result[i] = consts.LetterBytes[num.Int64()]
		}
		sessionID = string(result)

		// check if this session id alredy exists
		rdb := db.RedisConnect()
		key := "session:" + sessionID
		res, err := rdb.Exists(ctx, key).Result()
		if err != nil {
			return "", fmt.Errorf("error creating a SessionID %v", err.Error())
		}

		if res > 0 {
			length = l
			continue
		} else {
			break
		}
	}
	return sessionID, nil
}
