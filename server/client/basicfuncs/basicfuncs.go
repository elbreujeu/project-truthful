package basicfuncs

import (
	"math/rand"
	"strconv"
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
