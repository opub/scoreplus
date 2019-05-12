package util

import (
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/speps/go-hashids"
)

var (
	hic      *hashids.HashID
	hil      *hashids.HashID
	onceHash sync.Once
	baseTime int64 = 1554500000
)

const (
	linkChars     = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789" //removed ambiguous chars
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

//RandomString generates random string using masking with source
func RandomString(length int) string {
	b := make([]byte, length)
	l := len(letterBytes)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

//EncodeCookie encodes an ID and date as a string
func EncodeCookie(id int64, date time.Time) string {
	onceHash.Do(initHashIDs)
	parts := []int64{rand.Int63n(100000), id, date.Unix()}
	x, err := hic.EncodeInt64(parts)
	if err != nil {
		log.Error().Ints64("parts", parts).Msg("failed to encode cookie")
	}
	return x
}

//DecodeCookie decodes an ID and date from a string
func DecodeCookie(hash string) (int64, time.Time) {
	onceHash.Do(initHashIDs)
	x, err := hic.DecodeInt64WithError(hash)
	if err != nil || len(x) != 3 {
		log.Error().Str("hash", hash).Msg("failed to decode cookie")
		return 0, time.Now().AddDate(-1, 0, 0)
	}

	id := x[1]
	date := time.Unix(x[2], 0)
	return id, date
}

//EncodeLink encodes an ID as a string
func EncodeLink(id int64, t int) string {
	onceHash.Do(initHashIDs)
	x, err := hil.EncodeInt64([]int64{id, int64(t)})
	if err != nil {
		log.Error().Int64("id", id).Int("type", t).Msg("failed to encode link")
	}
	return x
}

//DecodeLink decodes an ID from a string
func DecodeLink(hash string) int64 {
	onceHash.Do(initHashIDs)
	x, err := hil.DecodeInt64WithError(hash)
	if err != nil || len(x) != 2 {
		log.Error().Str("hash", hash).Msg("failed to decode link")
		return 0
	}
	return x[0]
}

func initHashIDs() {
	config := GetConfig()

	//initialize hashids data for cookies
	data := hashids.NewData()
	data.Salt = config.Salt
	data.Alphabet = letterBytes
	data.MinLength = 24
	hic, _ = hashids.NewWithData(data)

	//initialize hashids data for links
	data = hashids.NewData()
	data.Salt = config.Salt
	data.Alphabet = linkChars
	data.MinLength = 5
	hil, _ = hashids.NewWithData(data)

	//make rand nondeterministic
	rand.Seed(time.Now().UnixNano())
}
