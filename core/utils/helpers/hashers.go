package helpers

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	b64 "encoding/base64"
	"encoding/json"

	"golang.org/x/crypto/pbkdf2"
)

const (
	UNUSABLE_PASSWORD_PREFIX        = "!"
	UNUSABLE_PASSWORD_SUFFIX_LENGTH = 40
	RANDOM_STRING_CHARS             = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	algorithm                       = "pbkdf2_sha256"
	iterations                      = 260000
	digest                          = "sha256"
	saltEntropy                     = 128
)

var letterRunes = []rune(RANDOM_STRING_CHARS)

// GetRandomString Return a securely generated random string.
// The bit length of the returned value can be calculated with the formula:
// log_2(len(allowed_chars)^length)
// For example, with default `allowed_chars` (26+26+10), this gives:
//   - length: 12, bit length =~ 71 bits
//   - length: 22, bit length =~ 131 bits
func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// IsPasswordUsable return True if this password wasn't generated by
// MakePassword(None).
func IsPasswordUsable(encoded string) bool {
	if encoded == "" {
		return false
	}
	return strings.HasPrefix(encoded, UNUSABLE_PASSWORD_PREFIX)
}

// MakePassword Turn a plain-text password into a hash for database storage

// Same as encode() but generate a new random salt. If password is None then
// return a concatenation of UNUSABLE_PASSWORD_PREFIX and a random string,
// which disallows logins. Additional random string reduces chances of gaining
// access to staff or superuser accounts. See ticket #20079 for more info.
func MakePassword(password string, salt string) (string, error) {
	if password == "" {
		return "", errors.New("password is not a valid password")
	}
	hasher := PBKDF2PasswordHasher{}
	if salt == "" {
		salt = hasher.Salt()
	}

	return hasher.Encode(password, salt, iterations)

}

type PBKDF2PasswordHasher struct {
}

type DecodePasswordHasher struct {
	Algorithm  string `json:"algorithm"`
	Iterations int64  `json:"iterations"`
	Salt       string `json:"salt"`
	Hash       string `json:"hash"`
}

// Salt generate a cryptographically secure nonce salt in ASCII with an entropy
// of at least `salt_entropy` bits.
func (h *PBKDF2PasswordHasher) Salt() string {
	l := len(RANDOM_STRING_CHARS)
	charCount := math.Ceil(saltEntropy / math.Log2(float64(l)))
	return GetRandomString(int(charCount))
}

func (h *PBKDF2PasswordHasher) Encode(password string, salt string, iterations int) (string, error) {
	if password == "" {
		return "", errors.New("password must not be empty")
	}
	if strings.Contains(salt, "$") {
		return "", errors.New("salt shold not have character $")
	}
	hash := pbkdf2.Key([]byte(password), []byte(salt), iterations, sha256.Size, sha256.New)
	return fmt.Sprintf("%s$%d$%s$%s", algorithm, iterations, salt, b64.StdEncoding.EncodeToString(hash)), nil
}

func (h *PBKDF2PasswordHasher) Decode(encoded string) (DecodePasswordHasher, error) {
	decodeArray := strings.Split(encoded, "$")
	algorithm, itr, salt, hash := decodeArray[0], decodeArray[1], decodeArray[2], decodeArray[3]
	it, err := strconv.ParseInt(itr, 10, 64)
	return DecodePasswordHasher{Algorithm: algorithm, Iterations: it, Salt: salt, Hash: hash}, err
}

func (h *PBKDF2PasswordHasher) Verify(password, encoded string) bool {
	decoded, err := h.Decode(encoded)
	if err != nil {
		return false
	}
	encoded_2, err := h.Encode(password, decoded.Salt, int(decoded.Iterations))

	if err != nil {
		return false
	}
	return strings.Compare(encoded, encoded_2) == 0
}

func (h *PBKDF2PasswordHasher) SafeSummary(encoded string) ([]byte, error) {
	d, _ := h.Decode(encoded)
	return json.Marshal(d)
}

func (h *PBKDF2PasswordHasher) MustUpdate(encoded string) bool {

	return false
}

func (h *PBKDF2PasswordHasher) HardenRuntime(password, encoded string) (string, error) {
	decoded, _ := h.Decode(encoded)
	extra_iterations := iterations - decoded.Iterations
	if extra_iterations > 0 {
		return h.Encode(password, decoded.Salt, int(extra_iterations))
	}

	return password, nil
}