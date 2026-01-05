package utils

import (
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

const (
	SchemeID   = "fiber-aead-xchacha20p1305-v1"
	HKDFInfo   = "hkdf:" + SchemeID
	DefaultAAD = "ctx:global"

	Argon2Time    = 1
	Argon2Memory  = 64 * 1024 // 64MB
	Argon2Threads = 4
	Argon2KeyLen  = 32
	Argon2SaltLen = 16
)

var hkdfSalt = []byte{0x72, 0x8b, 0x55, 0x1e, 0x9a, 0x11, 0x3f, 0xcc, 0x20, 0x01, 0x9e, 0x02, 0xa7, 0x44, 0xd1, 0x5a}

type CryptoService struct{ aead cipher.AEAD }

func NewCryptoService(appKey []byte) (*CryptoService, error) {
	if len(appKey) == 0 {
		return nil, errors.New("APP_KEY is empty")
	}

	// Turunkan key 32B dengan HKDF(SHA-256)
	r := hkdf.New(sha256.New, appKey, hkdfSalt, []byte(HKDFInfo))
	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, fmt.Errorf("hkdf: %w", err)
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("new aead: %w", err)
	}

	for i := range key {
		key[i] = 0
	}

	return &CryptoService{aead: aead}, nil
}

func (c *CryptoService) Encrypt(plain []byte) (string, error) {
	nonce := make([]byte, chacha20poly1305.NonceSizeX) // 24B, WAJIB acak
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("nonce rng: %w", err)
	}

	ct := c.aead.Seal(nil, nonce, plain, []byte(DefaultAAD))
	out := append(append(make([]byte, 0, len(nonce)+len(ct)), nonce...), ct...)

	return base64.StdEncoding.EncodeToString(out), nil
}

func (c *CryptoService) Decrypt(b64 string) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("base64: %w", err)
	}

	if len(raw) < chacha20poly1305.NonceSizeX {
		return nil, errors.New("ciphertext too short")
	}

	nonce := raw[:chacha20poly1305.NonceSizeX]
	ct := raw[chacha20poly1305.NonceSizeX:]

	pt, err := c.aead.Open(nil, nonce, ct, []byte(DefaultAAD))
	if err != nil {
		fmt.Printf("❌❌❌ DECRYPTION FAILED ❌❌❌\n")
		fmt.Printf("Error: %v\n", err)

		// Debug: Try without AAD
		fmt.Printf("\n--- DEBUG: Trying without AAD ---\n")
		ptNoAAD, err2 := c.aead.Open(nil, nonce, ct, nil)
		if err2 != nil {
			fmt.Printf("Also failed without AAD: %v\n", err2)
		} else {
			fmt.Printf("⚠️  SUCCESS without AAD! Plaintext: %s\n", string(ptNoAAD))
			fmt.Printf("⚠️  This indicates AAD mismatch!\n")
		}

		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return pt, nil
}

func (c *CryptoService) HashString(password string) (string, error) {
	salt := make([]byte, Argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt rng: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)

	combined := append(salt, hash...)

	return base64.StdEncoding.EncodeToString(combined), nil
}

func (c *CryptoService) VerifyString(password, hashedPassword string) (bool, error) {
	decoded, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, fmt.Errorf("base64: %w", err)
	}

	if len(decoded) < Argon2SaltLen {
		return false, errors.New("hashed string too short")
	}

	salt := decoded[:Argon2SaltLen]
	storedHash := decoded[Argon2SaltLen:]

	computedHash := argon2.IDKey([]byte(password), salt, Argon2Time, Argon2Memory, Argon2Threads, Argon2KeyLen)

	if len(computedHash) != len(storedHash) {
		return false, errors.New("not identic")
	}

	for i := 0; i < len(computedHash); i++ {
		if computedHash[i] != storedHash[i] {
			return false, errors.New("not identic")
		}
	}

	return true, nil
}
