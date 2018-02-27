package logclient

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// ComputeHmac256 computes HMAC with given secret
func computeHmac256(message string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
