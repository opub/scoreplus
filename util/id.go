package util

import (
	"math/rand"
	"time"

	"github.com/speps/go-hashids"
)

var hi *hashids.HashID
var baseTime int64 = 1554500000

func init() {
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	//initialize hashids data
	data := hashids.NewData()
	data.Salt = config.Salt
	data.Alphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789" //remove ambiguous chars
	data.MinLength = 6
	hi, _ = hashids.NewWithData(data)

	//make rand nondeterministic
	rand.Seed(time.Now().UnixNano())
}

//RandomID produces a new random ID
func RandomID() (string, error) {

	//our id is based on current epoch nanoseconds and a pseudo-random 63-bit int for uniqueness
	now := int(time.Now().Unix() - baseTime)
	rnd := rand.Intn(100000)

	id, err := hi.Encode([]int{now, rnd})
	if err != nil {
		return "", err
	}

	return id, nil
}
