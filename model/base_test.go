package model

import (
	"github.com/opub/scoreplus/util"
)

func random() string {
	value, _ := util.RandomID()
	return value
}
