package business

import (
	"regexp"
	"strings"
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func cleanPhoneNumber(phoneNumber string) string {
	re := regexp.MustCompile(`[^\d]`)
	return re.ReplaceAllString(phoneNumber, "")
}

func hasFiveIdenticalDigits(s string) bool {
	if len(s) < 5 {
		return false
	}

	for i := 0; i <= len(s)-5; i++ {
		if s[i] == s[i+1] && s[i] == s[i+2] && s[i] == s[i+3] && s[i] == s[i+4] {
			return true
		}
	}
	return false
}

func hasFiveConsecutiveDigits(s string) bool {
	consecutivePatterns := []string{
		"01234", "12345", "23456", "34567", "45678", "56789", "67890",
		"98765", "87654", "76543", "65432", "54321", "43210",
	}
	for _, pattern := range consecutivePatterns {
		if strings.Contains(s, pattern) {
			return true
		}
	}
	return false
}

func (s *ValidationService) IsValidPhoneNumber(phoneNumber string) bool {

	cleanedNumber := cleanPhoneNumber(phoneNumber)

	initialPattern := `^(998|998)(50|55|65|77|88|90|91|93|94|97|99)\d{3}\d{2}\d{2}$`
	reInitial := regexp.MustCompile(initialPattern)

	if !reInitial.MatchString(cleanedNumber) {
		return false
	}

	if len(cleanedNumber) < 7 {
		return false
	}
	last7Digits := cleanedNumber[len(cleanedNumber)-7:]

	if hasFiveIdenticalDigits(last7Digits) {
		return false
	}

	if hasFiveConsecutiveDigits(last7Digits) {
		return false
	}

	return true
}
