package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GetFileHash(data []byte, hashtype string) (string, error) {

	switch strings.ToLower(hashtype) {
	case "md5":
		h := md5.New()
		_, err := h.Write(data)
		if err != nil {
			return "", err
		}
		hashBytes := h.Sum(nil)
		hashstring := hex.EncodeToString(hashBytes)
		return hashstring, nil

	default:
		h := sha256.New()
		_, err := h.Write(data)
		if err != nil {
			return "", err
		}
		hashBytes := h.Sum(nil)
		hashstring := hex.EncodeToString(hashBytes[:])

		return hashstring, nil
	}
}
