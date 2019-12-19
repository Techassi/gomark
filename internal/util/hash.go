package util

import (
	"encoding/base64"
	"fmt"
	"runtime"

	"github.com/Techassi/gomark/internal/constants"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

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
