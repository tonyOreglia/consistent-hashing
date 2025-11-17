package hash

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

func HashKey(key string) uint64 {
	sum := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint64(sum[:8]) // first 8 bytes → uint64
}

func HashId(v string) string {
	sum := sha256.Sum256([]byte(v))
	// Use first 8 bytes of hash, hex-encoded, for a short readable string
	return hex.EncodeToString(sum[:8])
}
