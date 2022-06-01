package moneydb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseIBAN(t *testing.T) {
	requireValid := func(expected IBAN, str string) {
		iban, err := ParseIBAN(str)
		require.NoError(t, err)
		require.Equal(t, expected, iban)
	}
	requireInvalid := func(str string) {
		_, err := ParseIBAN(str)
		require.Error(t, err)
	}

	requireValid(IBAN("DE11520513735120710131"), "DE11520513735120710131")
	requireValid(IBAN("DE11520513735120710131"), "de11520513735120710131")
	requireValid(IBAN("GB33BUKB20201555555555"), "GB33BUKB20201555555555")
	requireValid(IBAN("DE02120300000000202051"), "DE02 1203 0000 0000 2020 51")
	requireInvalid("")
	requireInvalid("§$/$&(&")
	requireInvalid("DE")
	requireInvalid("DE1152051373512071013")
	requireInvalid("DE115205137351207101319")
	requireInvalid("GB01BARC20714583608387")
	requireInvalid("DE11520513735120710132")
}

func TestIBANString(t *testing.T) {
	require.Equal(t, "DE11520513735120710131", IBAN("DE11520513735120710131").String())
}

func TestParseEuro(t *testing.T) {
	requireValid := func(expected Money, str string) {
		m, err := ParseEuro(str)
		require.NoError(t, err)
		require.Equal(t, expected, m)
	}
	requireInvalid := func(str string) {
		_, err := ParseEuro(str)
		require.Error(t, err)
	}

	requireValid(Money(0), "0")
	requireValid(Money(0), "00")
	requireValid(Money(0), "0,00")
	requireValid(Money(1), "0,01")
	requireValid(Money(8), "0,08")
	requireValid(Money(10), "0,1")
	requireValid(Money(90), "0,9")
	requireValid(Money(10), "0,10")
	requireValid(Money(80), "0,80")
	requireValid(Money(100), "1")
	requireValid(Money(900), "9")
	requireValid(Money(100), "1,0")
	requireValid(Money(110), "1,1")
	requireValid(Money(190), "1,9")
	requireValid(Money(100), "1,00")
	requireValid(Money(1000), "10")
	requireValid(Money(10000), "100")
	requireValid(Money(-100), "-1")
	requireValid(Money(-1), "-0,01")
	requireValid(Money(307448), "3074,48")
	requireValid(Money(-307448), "-3074,48")
	requireValid(Money(-307448), "-3074.48")
	requireInvalid("")
	requireInvalid("sdf")
	requireInvalid("+1")
	requireInvalid("0,")
	requireInvalid("0,000")
	requireInvalid(",00")
	requireInvalid("- 1")
}

func TestMoneyString(t *testing.T) {
	require.Equal(t, "0,00 €", Money(0).String())
	require.Equal(t, "0,01 €", Money(1).String())
	require.Equal(t, "0,09 €", Money(9).String())
	require.Equal(t, "0,10 €", Money(10).String())
	require.Equal(t, "0,90 €", Money(90).String())
	require.Equal(t, "1,00 €", Money(100).String())
	require.Equal(t, "9,00 €", Money(900).String())
	require.Equal(t, "-0,01 €", Money(-1).String())
	require.Equal(t, "-0,09 €", Money(-9).String())
	require.Equal(t, "-0,10 €", Money(-10).String())
	require.Equal(t, "-0,90 €", Money(-90).String())
	require.Equal(t, "-1,00 €", Money(-100).String())
	require.Equal(t, "-9,00 €", Money(-900).String())
	require.Equal(t, "3074,48 €", Money(307448).String())
	require.Equal(t, "-3074,48 €", Money(-307448).String())
}

func TestMoneyPlus(t *testing.T) {
	m1 := Money(934)
	m2 := Money(-19)
	m3 := m1.Plus(m2)
	require.Equal(t, Money(915), m3)
	require.Equal(t, Money(934), m1)
	require.Equal(t, Money(-19), m2)
}

func TestMoneyMinus(t *testing.T) {
	m1 := Money(934)
	m2 := Money(19)
	m3 := m1.Minus(m2)
	require.Equal(t, Money(915), m3)
	require.Equal(t, Money(934), m1)
	require.Equal(t, Money(19), m2)
}
