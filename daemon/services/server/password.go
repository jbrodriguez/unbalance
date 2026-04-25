package server

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

const (
	passwordMinChars    = 8
	passwordMaxChars    = 256
	argon2Time          = 1
	argon2Memory        = 64 * 1024
	argon2Parallelism   = 4
	argon2SaltLength    = 16
	argon2KeyLength     = 32
	argon2VariantPrefix = "$argon2id$"
)

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < passwordMinChars {
		return fmt.Errorf("password must be at least %d characters", passwordMinChars)
	}

	if length > passwordMaxChars {
		return fmt.Errorf("password must be %d characters or fewer", passwordMaxChars)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	if err := validatePassword(password); err != nil {
		return "", err
	}

	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Parallelism, argon2KeyLength)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2Memory, argon2Time, argon2Parallelism, encodedSalt, encodedHash), nil
}

func verifyPassword(storedHash, password string) (bool, bool, error) {
	switch {
	case isBcryptHash(storedHash):
		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return false, false, nil
			}
			return false, false, err
		}
		return true, true, nil
	case strings.HasPrefix(storedHash, argon2VariantPrefix):
		ok, err := verifyArgon2idPassword(storedHash, password)
		return ok, false, err
	case storedHash == "":
		return false, false, nil
	default:
		return false, false, fmt.Errorf("unsupported password hash format")
	}
}

func isBcryptHash(hash string) bool {
	return strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$")
}

type argon2Params struct {
	memory      uint32
	time        uint32
	parallelism uint8
	salt        []byte
	hash        []byte
}

func verifyArgon2idPassword(encodedHash, password string) (bool, error) {
	params, err := parseArgon2idHash(encodedHash)
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey([]byte(password), params.salt, params.time, params.memory, params.parallelism, uint32(len(params.hash)))
	return subtle.ConstantTimeCompare(computedHash, params.hash) == 1, nil
}

func parseArgon2idHash(encodedHash string) (*argon2Params, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return nil, fmt.Errorf("invalid argon2id hash format")
	}

	versionParts := strings.SplitN(parts[2], "=", 2)
	if len(versionParts) != 2 || versionParts[0] != "v" {
		return nil, fmt.Errorf("invalid argon2id version")
	}
	version, err := strconv.Atoi(versionParts[1])
	if err != nil {
		return nil, err
	}
	if version != argon2.Version {
		return nil, fmt.Errorf("unsupported argon2 version %d", version)
	}

	params := &argon2Params{}
	for _, item := range strings.Split(parts[3], ",") {
		kv := strings.SplitN(item, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid argon2 parameter")
		}

		switch kv[0] {
		case "m":
			value, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return nil, err
			}
			params.memory = uint32(value)
		case "t":
			value, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return nil, err
			}
			params.time = uint32(value)
		case "p":
			value, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil {
				return nil, err
			}
			params.parallelism = uint8(value)
		default:
			return nil, fmt.Errorf("unsupported argon2 parameter %q", kv[0])
		}
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, err
	}

	params.salt = salt
	params.hash = hash
	return params, nil
}
