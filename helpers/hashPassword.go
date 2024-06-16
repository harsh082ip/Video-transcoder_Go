package helpers

import "golang.org/x/crypto/bcrypt"

// The cost factor (in this case, 14) is a measure of how
// computationally expensive the hashing should be. Higher cost
// factors result in slower hash generation but also make it more
// difficult for attackers to perform brute-force or rainbow table attacks.

func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
