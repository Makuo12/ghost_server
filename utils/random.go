package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)
const alphabet = "abcdefghijklmnopqrstuvwxyz"
//const alphnumeric = "0a1b2c3d4e5f6g7h8i9jklmnopqrstuvwxyz"
func init() {
	//the code below would make sure each time we generate random number is different
	rand.Seed(time.Now().UnixNano())
}
var MyMessage = "I like to dance and eat cheese with ice cream"
//now we write a function that returns a random integer 
func RandomInt(min int64, max int64) int64{
	return min + rand.Int63n(max - min + 1)
}

func RandomNumber(n int) string {
	var sb strings.Builder
	for i:=0;i<=n;i++{
		c := rand.Intn(9)
		sb.WriteString(strconv.Itoa(c))

	}
	return sb.String()
}
func RandomString(n int) string{
	var sb strings.Builder
	k := len(alphabet)

	for i:=0; i<n; i++{
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomEmail() string{
	return fmt.Sprintf("%v@gmail.com",RandomString(6))
}
func RandomName() string{
	return fmt.Sprintf("%v",RandomString(7))
}

