package splicerouter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func RandomInt() string {
	return fmt.Sprint(time.Now().Nanosecond())
}

func IsLast(index int, len int) bool {
	return index+1 == len
}

func RandomString() string {
	return uuid.NewString()
}
