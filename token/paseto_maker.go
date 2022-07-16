package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey string
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size : should be at least %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}, nil
}

func (pm *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return pm.paseto.Encrypt([]byte(pm.symmetricKey), payload, nil)
}

func (pm *PasetoMaker) VerifyToken(token string) (*Payload, error) {

	payload := &Payload{}

	err := pm.paseto.Decrypt(token, []byte(pm.symmetricKey), payload, nil)
	if err != nil {
		return nil, Err_Invalid_Token
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
