package handlers

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// PasswordValidationRule encompasses a password validation rule.
type PasswordValidationRule interface {
	Validate(password string) error
}

// PasswordRuleFunction is a validation rule function to be plugged in to a chain of Validate calls.
type PasswordRuleFunction func(string) error

// Validate applies a given PasswordRuleFunction to a password.
func (f PasswordRuleFunction) Validate(password string) error {
	return f(password)
}

// PasswordValidator validates passwords by applying a sequence of rules.
type PasswordValidator struct {
	rules []PasswordValidationRule
}

// Validate applies validation rules in sequence on the password.
func (v PasswordValidator) Validate(password string) []error {
	var errors []error

	for _, rule := range v.rules {
		if err := rule.Validate(password); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

var (
	minimumLengthRule PasswordValidationRule = PasswordRuleFunction(func(password string) error {
		if utf8.RuneCountInString(password) < 8 {
			return fmt.Errorf("must be at least 8 characters")
		}

		return nil
	})

	twoSpecialSymbolsRule PasswordValidationRule = PasswordRuleFunction(func(password string) error {
		count := 0

		for _, r := range password {
			if unicode.IsPunct(r) || unicode.IsSymbol(r) {
				count++
			}
		}

		if count < 2 {
			return fmt.Errorf("must contain at least 2 special symbols")
		}

		return nil
	})

	oneNumberRule PasswordValidationRule = PasswordRuleFunction(func(password string) error {
		for _, r := range password {
			if unicode.IsDigit(r) {
				return nil
			}
		}
		return fmt.Errorf("must contain at least 1 number")
	})

	oneUppercaseLetterRule PasswordValidationRule = PasswordRuleFunction(func(password string) error {
		for _, r := range password {
			if unicode.IsLetter(r) && unicode.IsUpper(r) {
				return nil
			}
		}

		return fmt.Errorf("must contain at least 1 uppercase letter")
	})
)

// BasePasswordValidator exposes a password validator which requires len > 8,
// two special symbols, at least one capital letter and at least one number.
var BasePasswordValidator = PasswordValidator{
	rules: []PasswordValidationRule{
		minimumLengthRule, twoSpecialSymbolsRule,
		oneNumberRule, oneUppercaseLetterRule,
	},
}
