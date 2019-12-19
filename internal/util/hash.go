package util

import (
	"runtime"

	"github.com/Techassi/gomark/internal/constants"
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

func HashPassword() (string, error) {
	p := &Argon2Params{}
	p.init()

	return "", nil
}
