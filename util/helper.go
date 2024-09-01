package util

import "net/mail"

func CheckNumberPowerOfTwo(n int) int {
	return n & (n - 1)
}

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
