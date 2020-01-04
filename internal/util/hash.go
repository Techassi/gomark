package util

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/adler32"
	"runtime"
	"strings"

	"github.com/Techassi/gomark/internal/constants"

	"golang.org/x/crypto/argon2"
)

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// ARGON2 FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func (p *Argon2Params) init() {
	p.Memory = 64 * constants.ARGON2_KIBIBYTE
	p.Iterations = 3
	p.Parallelism = uint8(runtime.NumCPU() / 4)
	p.SaltLength = 16
	p.KeyLength = 32
}

func HashPassword(pass string) (string, error) {
	p := &Argon2Params{}
	p.init()

	// Generate cryptographically secure random salt
	salt, err := RandomByteSlice(p.SaltLength)
	if err != nil {
		return "", err
	}

	// Generate the salted hash of the password
	hash := argon2.IDKey([]byte(pass), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Generate the encoded representation of the password
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	passB64 := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, saltB64, passB64)

	return encoded, nil
}

func ComparePassword(inputPass, hashedPass string) (bool, error) {
	p, salt, hash, err := decodeHash(hashedPass)
	if err != nil {
		return false, err
	}

	compareHash := argon2.IDKey([]byte(inputPass), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	if subtle.ConstantTimeCompare(hash, compareHash) == 1 {
		return true, nil
	}

	return false, nil
}

func decodeHash(encodedHash string) (p *Argon2Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// ENTITY FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func EntityHash(name, url string) string {
	h := adler32.New()
	s := fmt.Sprintf("%s%s", name, url)

	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func EntityHashPlusString(name string) string {
	h := adler32.New()
	s := fmt.Sprintf("%s%s", name, RandomString(10))

	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
