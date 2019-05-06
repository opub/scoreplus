package util

import (
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/speps/go-hashids"
)

var (
	hi       *hashids.HashID
	onceHash sync.Once
	baseTime int64 = 1554500000
)

const (
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

	log.Debug().Ints64("parts", parts).Msg("encoding parts")

	x, err := hi.EncodeInt64(parts)
	if err != nil {
		log.Error().Ints64("parts", parts).Msg("failed to encode id")
	}
	return x
}

//DecodeCookie decodes an ID and date from a string
func DecodeCookie(hash string) (int64, time.Time) {
	onceHash.Do(initHashIDs)

	x, err := hi.DecodeInt64WithError(hash)
	if err != nil || len(x) != 3 {
		log.Error().Str("hash", hash).Msg("failed to decode id")
		return 0, time.Now().AddDate(-1, 0, 0)
	}

	log.Debug().Ints64("parts", x).Msg("decoded parts")

	id := x[1]
	date := time.Unix(x[2], 0)
	return id, date
}

func initHashIDs() {
	config := GetConfig()

	//initialize hashids data
	data := hashids.NewData()
	data.Salt = config.Salt
	data.Alphabet = letterBytes
	data.MinLength = 24
	hi, _ = hashids.NewWithData(data)

	//make rand nondeterministic
	rand.Seed(time.Now().UnixNano())
}
