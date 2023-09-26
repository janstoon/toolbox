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
	ErrIbanIncorrectLength = errors.New("iban length incorrect")
	ErrIbanUnknownCountry  = errors.New("iban country unknown")
	ErrIbanCheckFailure    = errors.New("iban check failure")
	ErrIbanInvalidBban     = errors.New("iban invalid bban")
)

type InternationalBankAccountNumber struct {
	countryCode            string
	checkDigits            string
	basicBankAccountNumber string
}

func (iban InternationalBankAccountNumber) String() string {
	var sb strings.Builder
	sb.WriteString(iban.countryCode)
	sb.WriteString(iban.checkDigits)
	sb.WriteString(iban.basicBankAccountNumber)

	return sb.String()
}

type (
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

func ParseInternationalBankAccountNumber(ibanStr string) (*InternationalBankAccountNumber, error) {
	const (
		ibanMinLength                   = 5
		ibanMaxLength                   = 34
		ibanCountryCodeStart            = 0
		ibanCountryCodeEnd              = 2
		ibanCheckDigitsStart            = 2
		ibanCheckDigitsEnd              = 4
		ibanBasicBankAccountNumberStart = 4
	)

	if l := len(ibanStr); l < ibanMinLength || l > ibanMaxLength {
		return nil, errors.Join(ErrInvalidInput, ErrIbanIncorrectLength)
	}

	ibanStr = strings.ToUpper(ibanStr)
	cc := ibanStr[ibanCountryCodeStart:ibanCountryCodeEnd]
	cd := ibanStr[ibanCheckDigitsStart:ibanCheckDigitsEnd]
	bban := ibanStr[ibanBasicBankAccountNumberStart:]

	if err := validateIban(cc, cd, bban); err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	validator := bbanValidators.lookupByCountryCode(cc)
	if validator == nil {
		return nil, errors.Join(ErrInvalidInput, ErrIbanUnknownCountry)
	}

	if err := validator(bban); err != nil {
		return nil, errors.Join(ErrInvalidInput, err)
	}

	return &InternationalBankAccountNumber{
		countryCode:            cc,
		checkDigits:            cd,
		basicBankAccountNumber: bban,
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

	if 1 == value.Mod(value, big.NewInt(97)).Int64() {
		return nil
	}

	return ErrIbanCheckFailure
}

func RegisterBbanValidator(countryCode string, validator BbanValidator) {
	bbanValidators.push(bbanValidatorCatalogue{
		countryCode: countryCode,
		validator:   validator,
	})
}

type PrimaryAccountNumber struct {
	issuerIdentificationNumber string
	accountNumber              string
}

func ParsePrimaryAccountNumber(str string) (*PrimaryAccountNumber, error) {
	return nil, ErrNotImplemented
}

type IranPrimaryAccountNumber struct {
	Institute   [6]byte
	Product     [2]byte
	Serial      [7]byte
	CheckDigits [1]byte
}

type BankAccount struct {
	iban InternationalBankAccountNumber
	pan  PrimaryAccountNumber

	bank string
}

func bankAccountFromIban(iban InternationalBankAccountNumber) BankAccount {
	return BankAccount{}
}

func bankAccountFromPan(pan PrimaryAccountNumber) BankAccount {
	return BankAccount{}
}

func (ba BankAccount) BankName() string {
	return ba.bank
}
