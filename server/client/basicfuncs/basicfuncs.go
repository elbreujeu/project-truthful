package basicfuncs

import (
	"math/rand"
	"strconv"
	"strings"
)

func GenerateRandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func ConvertQueryParameterToInt(rawParameterString string, defaultValue int) (int, error) {
	var toReturn int
	var err error
	if rawParameterString == "" {
		toReturn = defaultValue
	} else {
		toReturn, err = strconv.Atoi(rawParameterString)
		if err != nil {
			return 0, err
		}
	}
	return toReturn, nil
}

func DeleteNonAlphanumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}
