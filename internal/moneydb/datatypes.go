package moneydb

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	jbub_iban "github.com/jbub/banking/iban"
)

var (
	patternMoney = regexp.MustCompile(`^\s*(-?)(\d+)([.,](\d{1,2}))?\s*$`)
)

type IBAN string

func ParseIBAN(str string) (IBAN, error) {
	ibanObj, err := jbub_iban.Parse(strings.ToUpper(strings.ReplaceAll(str, " ", "")))
	if err != nil {
		return "", err
	}
	return IBAN(ibanObj.String()), nil
}

func (iban IBAN) String() string {
	return string(iban)
}

type Money int64

func ParseEuro(str string) (Money, error) {
	m := patternMoney.FindStringSubmatch(str)
	if len(m) != 5 {
		return 0, fmt.Errorf("invalid format")
	}

	euros, _ := strconv.ParseInt(m[2], 10, 64)
	cents, _ := strconv.ParseInt(m[4], 10, 64)
	if len(m[4]) == 1 {
		cents *= 10
	}
	totalCents := 100*euros + cents
	if m[1] == "-" {
		totalCents = -totalCents
	}

	return Money(totalCents), nil
}

func (m Money) String() string {
	totalCents := m

	suffix := " €"

	var prefix string
	if totalCents < 0 {
		prefix = "-"
		totalCents = -totalCents
	}

	str := fmt.Sprintf("%d", totalCents)
	if len(str) == 1 {
		return prefix + "0,0" + str + suffix
	} else if len(str) == 2 {
		return prefix + "0," + str + suffix
	} else {
		return prefix + str[:len(str)-2] + "," + str[len(str)-2:] + suffix
	}
}

func (m Money) Plus(other Money) Money {
	return m + other
}

func (m Money) Minus(other Money) Money {
	return m - other
}
