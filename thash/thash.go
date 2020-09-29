package thash

import "crypto/sha256"

// hash generates a 256-bit hash
func Hash(buf []byte) []byte {
	sha256 := func(x []byte) []byte {
		h := sha256.Sum256(x)
		return h[:]
	}
	return sha256(sha256(buf))
}

// TODO: replace with something like scrypt
func Mine(buf []byte) []byte {
	return Hash(buf)
}
