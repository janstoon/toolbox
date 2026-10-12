package bricks

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/janstoon/toolbox/tricks"
)

var (
	// ErrIbanIncorrectLength is returned when an IBAN is shorter than 5 or longer than 34 characters.
	ErrIbanIncorrectLength = errors.New("iban length incorrect")
	// ErrIbanUnknownCountry is returned when no BBAN validator is registered for the IBAN country code.
	ErrIbanUnknownCountry = errors.New("iban country unknown")
	// ErrIbanCheckFailure is returned when the IBAN mod-97 check digits do not match.
	ErrIbanCheckFailure = errors.New("iban check failure")
	// ErrIbanInvalidBban is returned when the registered country validator rejects the BBAN.
	ErrIbanInvalidBban = errors.New("iban invalid bban")

	// ErrPanIncorrectLength is returned when a primary account number is shorter than 8 or longer than 19 characters.
	ErrPanIncorrectLength = errors.New("pan length incorrect")
)

// InternationalBankAccountNumber is an ISO 13616 IBAN in uppercase, without spaces.
type InternationalBankAccountNumber struct {
	full                   string
	countryCode            string
	checkDigits            string
	basicBankAccountNumber string

	// Country is the country registered for the IBAN's ISO alpha-2 code, or nil when that country is unknown.
	Country *Country
}

// String returns the compact uppercase IBAN.
func (iban InternationalBankAccountNumber) String() string {
	return iban.full
}

type (
	// BbanValidator checks the basic bank account number of one country.
	// Register it with RegisterBbanValidator.
	BbanValidator          func(bban string) error
	bbanValidatorCatalogue struct {
		countryCode string
		validator   BbanValidator
	}
)

type bbanValidatorBank struct {
	sync.RWMutex

	byCountryCode   map[string]BbanValidator
	normCountryCode func(string) string
}

func (bank *bbanValidatorBank) push(c bbanValidatorCatalogue) {
	bank.Lock()
	defer bank.Unlock()

	tricks.InsertIfNotExist(bank.byCountryCode, bank.normCountryCode(c.countryCode), c.validator)
}

func (bank *bbanValidatorBank) lookupByCountryCode(code string) BbanValidator {
	bank.RLock()
	defer bank.RUnlock()

	return bank.byCountryCode[bank.normCountryCode(code)]
}

var bbanValidators = bbanValidatorBank{
	byCountryCode:   make(map[string]BbanValidator),
	normCountryCode: strings.ToUpper,
}

// ParseInternationalBankAccountNumber parses a compact IBAN.
// The country code must have a validator registered with RegisterBbanValidator.
// On failure the error wraps ErrInvalidArgument and one of the ErrIban errors.
func ParseInternationalBankAccountNumber(iban string) (*InternationalBankAccountNumber, error) {
	const (
		ibanMinLength                   = 5
		ibanMaxLength                   = 34
		ibanCountryCodeStart            = 0
		ibanCountryCodeEnd              = 2
		ibanCheckDigitsStart            = 2
		ibanCheckDigitsEnd              = 4
		ibanBasicBankAccountNumberStart = 4
	)

	if l := len(iban); l < ibanMinLength || l > ibanMaxLength {
		return nil, errors.Join(ErrInvalidArgument, ErrIbanIncorrectLength)
	}

	iban = strings.ToUpper(iban)
	cc := iban[ibanCountryCodeStart:ibanCountryCodeEnd]
	cd := iban[ibanCheckDigitsStart:ibanCheckDigitsEnd]
	bban := iban[ibanBasicBankAccountNumberStart:]

	if err := validateIban(cc, cd, bban); err != nil {
		return nil, errors.Join(ErrInvalidArgument, err)
	}

	validator := bbanValidators.lookupByCountryCode(cc)
	if validator == nil {
		return nil, errors.Join(ErrInvalidArgument, ErrIbanUnknownCountry)
	}

	if err := validator(bban); err != nil {
		return nil, errors.Join(ErrInvalidArgument, err)
	}

	return &InternationalBankAccountNumber{
		full:                   iban,
		countryCode:            cc,
		checkDigits:            cd,
		basicBankAccountNumber: bban,

		Country: LookupCountryByIsoAlphaTwoCode(cc),
	}, nil
}

var base36DigitsTable = map[rune]string{
	'0': "0", '1': "1", '2': "2", '3': "3", '4': "4",
	'5': "5", '6': "6", '7': "7", '8': "8", '9': "9",
	'A': "10", 'B': "11", 'C': "12", 'D': "13", 'E': "14",
	'F': "15", 'G': "16", 'H': "17", 'I': "18", 'J': "19",
	'K': "20", 'L': "21", 'M': "22", 'N': "23", 'O': "24",
	'P': "25", 'Q': "26", 'R': "27", 'S': "28", 'T': "29",
	'U': "30", 'V': "31", 'W': "32", 'X': "33", 'Y': "34",
	'Z': "35",
}

func validateIban(cc, cd, bban string) error {
	rearranged := fmt.Sprintf("%s%s%s", bban, cc, cd)
	var numeric strings.Builder
	for _, c := range rearranged {
		numeric.WriteString(base36DigitsTable[c])
	}

	value := new(big.Int)
	value.SetString(numeric.String(), 10)

	if value.Mod(value, big.NewInt(97)).Int64() == 1 {
		return nil
	}

	return ErrIbanCheckFailure
}

// RegisterBbanValidator registers validator for an ISO 3166-1 alpha-2 country code.
// The match is case-insensitive. Registering the same country twice keeps the first validator.
func RegisterBbanValidator(countryCode string, validator BbanValidator) {
	bbanValidators.push(bbanValidatorCatalogue{
		countryCode: countryCode,
		validator:   validator,
	})
}

// PrimaryAccountNumber is a payment-card number of 8 to 19 characters.
// Check-digit validation is not applied yet.
type PrimaryAccountNumber struct {
	full                       string
	issuerIdentificationNumber string
	accountNumber              string
	checkDigit                 string
}

// String returns the digits passed to ParsePrimaryAccountNumber.
func (pan PrimaryAccountNumber) String() string {
	return pan.full
}

// ParsePrimaryAccountNumber accepts a primary account number of 8 to 19 characters.
// It returns an error that wraps ErrInvalidArgument and ErrPanIncorrectLength when the length is outside that range.
func ParsePrimaryAccountNumber(pan string) (*PrimaryAccountNumber, error) {
	const (
		panMinLength = 8
		panMaxLength = 19
	)

	if l := len(pan); l < panMinLength || l > panMaxLength {
		return nil, errors.Join(ErrInvalidArgument, ErrPanIncorrectLength)
	}

	// todo: validate pan

	return &PrimaryAccountNumber{
		full:                       pan,
		issuerIdentificationNumber: "",
		accountNumber:              "",
		checkDigit:                 "",
	}, nil
}

// BankAccount names an account at a bank. NewBankAccount records the account name.
// BankName returns the bank identifier, which stays empty until a later setter stores it.
type BankAccount struct {
	bank string
	name string

	iban InternationalBankAccountNumber
	pan  PrimaryAccountNumber
}

// NewBankAccount returns an account labeled with name.
func NewBankAccount(name string) BankAccount {
	return BankAccount{
		name: name,
	}
}

// BankName returns the bank identifier stored on the account.
func (ba BankAccount) BankName() string {
	return ba.bank
}
