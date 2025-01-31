package encryption

import (
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"sync"

	"github.com/minio/sha256-simd"
	"golang.org/x/crypto/sha3"
)

const HASH_LENGTH = 32

var sha3Pool = sync.Pool{
	New: func() interface{} {
		return sha3.New256()
	},
}

type HashBytes [HASH_LENGTH]byte

/*Hash - hash the given data and return the hash as hex string */
func Hash(data interface{}) string {
	return hex.EncodeToString(RawHash(data))
}

/*RawHash - Logic to hash the text and return the hash bytes */
func RawHash(data interface{}) []byte {
	var databuf []byte
	switch dataImpl := data.(type) {
	case []byte:
		databuf = dataImpl
	case HashBytes:
		databuf = dataImpl[:]
	case string:
		databuf = []byte(dataImpl)
	default:
		panic("unknown type")
	}
	hash := sha3Pool.Get().(hash.Hash)
	hash.Write(databuf)
	res := hash.Sum(nil)
	hash.Reset()
	sha3Pool.Put(hash)
	return res
}

func ShaHash(data interface{}) []byte {
	var databuf []byte
	switch dataImpl := data.(type) {
	case []byte:
		databuf = dataImpl
	case HashBytes:
		databuf = dataImpl[:]
	case string:
		databuf = []byte(dataImpl)
	default:
		panic("unknown type")
	}
	hash := sha256.New()
	_, _ = hash.Write(databuf)
	return hash.Sum(nil)
}

/*FastHash - sha1 hash the given data and return the hash as hex string */
func FastHash(data interface{}) string {
	return hex.EncodeToString(RawFastHash(data))
}

/*RawFastHash - Logic to sha1 hash the text and return the hash bytes */
func RawFastHash(data interface{}) []byte {
	var databuf []byte
	switch dataImpl := data.(type) {
	case []byte:
		databuf = dataImpl
	case HashBytes:
		databuf = dataImpl[:]
	case string:
		databuf = []byte(dataImpl)
	default:
		panic("unknown type")
	}
	hash := sha1.New()
	hash.Write(databuf)
	var buf []byte
	return hash.Sum(buf)
}
