package util

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyz"
const SPECIALCHARS = "!@#$%^&*()_-+=<>?"

// Initialize the random number generator with a new seed based on the current time
func init() {
	rand.NewSource(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max (inclusive) and returns it.
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString get a length size as n, generates a random string of length n using characters from the alphabet and returns it.
func RandomString(n int) (string, error) {
	var sb strings.Builder
	k := len(ALPHABET)
	for i := 0; i < n; i++ {
		c := ALPHABET[rand.Intn(k)]
		err := sb.WriteByte(c)
		if err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

// GenerateRandomEmail generates a random email address using a random name, random number, and a domain from the predefined list. It returns the generated email address.
func GenerateRandomEmail() (string, error) {
	randomName, err := RandomString(6)
	if err != nil {
		return "", errors.New("RandomString got error: " + err.Error())
	}
	randomNumber := strconv.FormatInt(RandomInt(1000, 9999), 10)
	domain := []string{"gmail.com", "yahoo.com", "hotmail.com"} // Add more domains if needed
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	randomDomain := domain[random.Intn(len(domain))]

	email := randomName + randomNumber + "@" + randomDomain
	return email, nil
}

// RandomPassword generates a random password of the specified length using characters from the alphabet, uppercase alphabet, special characters, and digits.
func RandomPassword(length int) (string, error) {
	var sb strings.Builder
	charset := ALPHABET + strings.ToUpper(ALPHABET) + SPECIALCHARS + "0123456789"
	charsetLength := len(charset)

	for i := 0; i < length; i++ {
		c := charset[rand.Intn(charsetLength)]
		err := sb.WriteByte(c)
		if err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

// RandomDateBeforeNow generates a random date before the current date in the format "2006-01-02".
func RandomDateBeforeNow() string {
	now := time.Now()
	randomTime := Randate(now)
	return randomTime.Format("2006-01-02")
}

// Randate generates a random date up to the specified max date.
func Randate(max time.Time) time.Time {
	year := rand.Intn(max.Year() + 1)
	month := time.Month(rand.Intn(int(time.December)) + 1)
	day := rand.Intn(28) + 1 // Assuming max day as 28 for simplicity

	randomTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return randomTime
}

// RandomDateTimeAfterNow generates a random date and time after the current date and time in the format "2006-01-02 15:04:05".
func RandomDateTimeAfterNow() string {
	now := time.Now()
	randomTime := randatetime(now)
	return randomTime.Format("2006-01-02 15:04:05")
}

// randatetime generates a random date and time after the specified min date and time.
func randatetime(min time.Time) time.Time {
	n := time.Since(min).Nanoseconds()
	if n > 0 {
		randomNanoseconds := rand.Int63n(n)
		randomDuration := time.Duration(randomNanoseconds)
		randomTime := min.Add(randomDuration)
		return randomTime
	} else {
		return min
	}
}
