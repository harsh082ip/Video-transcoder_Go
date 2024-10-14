package consts

import "time"

const (
	WEBPORT         = ":8006"
	LetterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	SessionTTL      = 600 * time.Second
	SessionIDlength = 16
)
