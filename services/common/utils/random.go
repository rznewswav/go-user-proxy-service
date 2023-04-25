package utils

import (
	"math/rand"
	"time"
)

var s1 = rand.NewSource(time.Now().UnixNano())
var Random = rand.New(s1)
