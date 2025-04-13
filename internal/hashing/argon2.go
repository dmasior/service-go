package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

type (
	Argon2       struct{}
	Argon2Params struct {
		memory      uint32
		iterations  uint32
		parallelism uint8
		saltLength  uint32
		keyLength   uint32
	}
)

var defaultParams = &Argon2Params{
	memory:      64 * 1024, // 64MB
	iterations:  3,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

func NewArgon2() *Argon2 {
	return &Argon2{}
}

func (a *Argon2) HashPassword(password string) (string, error) {
	salt, err := generateRandomBytes(defaultParams.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		defaultParams.iterations,
		defaultParams.memory,
		defaultParams.parallelism,
		defaultParams.keyLength,
	)

	// Encode salt and hash to base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id$v=19$m=memory,t=iterations,p=parallelism$salt$hash
	encodedHash := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		defaultParams.memory,
		defaultParams.iterations,
		defaultParams.parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func (a *Argon2) VerifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	if vals[1] != "argon2id" {
		return false, fmt.Errorf("unsupported algorithm: %s", vals[1])
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != 19 {
		return false, fmt.Errorf("unsupported version: %d", version)
	}

	params := &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.memory, &params.iterations, &params.parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false, err
	}
	saltLength, err := safeIntToUint32(len(salt))
	if err != nil {
		return false, err
	}
	params.saltLength = saltLength

	hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return false, err
	}
	keyLength, err := safeIntToUint32(len(hash))
	if err != nil {
		return false, err
	}
	params.keyLength = keyLength

	// Compute hash from the provided password
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.iterations,
		params.memory,
		params.parallelism,
		params.keyLength,
	)

	// Compare the computed hash with the stored hash
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// safeIntToUint32 safely converts an int to uint32, returning an error if the value is negative or exceeds uint32 max value
func safeIntToUint32(val int) (uint32, error) {
	if val < 0 {
		return 0, fmt.Errorf("cannot convert negative value %d to uint32", val)
	}
	if val > math.MaxUint32 {
		return 0, fmt.Errorf("value %d exceeds maximum uint32 value", val)
	}
	return uint32(val), nil
}
