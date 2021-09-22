package fixture

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func randNumberRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numberRunes[rand.Intn(len(numberRunes))]
	}

	return string(b)
}

func randStringLowerRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes)/2)]
	}

	return string(b)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandID() string {
	return randNumberRunes(15)
}

func Username() string {
	return RandStringRunes(RandInt(4, 15))
}

func Email() string {
	email := fmt.Sprintf("%s@example.com", randStringLowerRunes(RandInt(5, 10)))
	return strings.ToLower(email)
}

func RandStr(n int) string {
	return RandStringRunes(n)
}

// generateAvatar returns an gravatar using the md5 hash of the email
func generateAvatar(email string) string {
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://gravatar.com/avatar/%s?d=identicon", hex.EncodeToString(hash[:]))
}
