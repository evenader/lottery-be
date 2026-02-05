// Copyright (c) 2026 evenader. All rights reserved.

package utils

import "golang.org/x/crypto/bcrypt"

func BcryptWithString(str string) (string, error) {
	// DefaultCost 默认是 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(bytes), err
}
