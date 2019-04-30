package util

import (
	"math/rand"
	"sync"
	"time"

	"github.com/speps/go-hashids"
)

var (
	hi       *hashids.HashID
	onceHash sync.Once
	baseTime int64 = 1554500000
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
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

//RandomID produces a new random ID
func RandomID() (string, error) {
	onceHash.Do(func() {
		config := GetConfig()

		//initialize hashids data
		data := hashids.NewData()
		data.Salt = config.Salt
		data.Alphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789" //remove ambiguous chars
		data.MinLength = 10
		hi, _ = hashids.NewWithData(data)

		//make rand nondeterministic
		rand.Seed(time.Now().UnixNano())
	})

	//our id is based on current epoch nanoseconds and a pseudo-random 63-bit int for uniqueness
	now := int(time.Now().Unix() - baseTime)
	rnd := rand.Intn(10000)

	id, err := hi.Encode([]int{now, rnd})
	if err != nil {
		return "", err
	}

	return id, nil
}
