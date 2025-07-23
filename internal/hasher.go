package internal

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GetFileHash(data []byte, hashtype string) (string, error) {

	switch strings.ToLower(hashtype) {
	case "":
		h := md5.New()
		h.Write(data)
		hashBytes := h.Sum(nil)
		hashstring := hex.EncodeToString(hashBytes)
		return hashstring, nil

	default:
		h := sha256.New()
		h.Write(data)
		hashBytes := h.Sum(nil)
		hashstring := hex.EncodeToString(hashBytes)

		return hashstring, nil
	}
}
