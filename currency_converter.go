package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

func main() {
	a, b, c, d := Price{textformat: "251", numberformat: decimal.NewFromInt(250)}.RoundToNine()

	println("txt1: " + a)
	println("txt2: " + b)
	println("dec1: " + c.String())
	println("dec2: " + d.String())

	a, b, c, d = Price{textformat: "250.749322", numberformat: decimal.NewFromFloat(250.749322)}.RoundToNine()

	println("2txt1: " + a)
	println("2txt2: " + b)
	println("2dec1: " + c.String())
	println("2dec2: " + d.String())
}

func FormatPrice(ctx context.Context, price decimal.Decimal, currencyISO string) (string, error) {

	currencyISO = func() string {
		if currencyISO == "" {
			if sessionCurrency != defaultCurrency {
				price = ChangeCurrency(price, defaultCurrency, sessionCurrency)
			}
			return sessionCurrency
		}
		return currencyISO
	}()

	output := price.StringFixed(int32(currencies[currencyISO].Precision))

	if output == "" || currencyISO == "" || currencies[currencyISO].IsoCode == "" || currencies[currencyISO].Symbol == "" {
		return "", fmt.Errorf("cannot get currency")
	}

	return output + currencies[currencyISO].Symbol, nil
}

func ChangeCurrency(price decimal.Decimal, sourceCurrencyISO string, targetCurrencyISO string) decimal.Decimal {

	conversionRate :=
		currencies[sourceCurrencyISO].ConversionRate.
			Mul(decimal.NewFromInt(1).
				Div(currencies[targetCurrencyISO].ConversionRate))

	return price.Mul(conversionRate)
}

type Price struct {
	textformat   string
	numberformat decimal.Decimal
}

func (s Price) RoundToNine() (string, string, decimal.Decimal, decimal.Decimal) {

	txt1 := ""
	txt2 := ""
	dec1 := decimal.Zero
	dec2 := decimal.Zero

	parts := strings.Split(s.textformat, ".")
	if len(parts) == 2 {
		txt1 = parts[0] + ".9" + strings.Repeat("9", len(parts[1])-1)
		txt2 = strings.ReplaceAll(parts[0], string(parts[0][len(parts[0])-1]), "9") + ".9" + strings.Repeat("9", len(parts[1])-1)
		dec1 = func() decimal.Decimal {
			r, e := decimal.NewFromString(txt1)
			if e != nil {
				return decimal.Zero
			}
			return r
		}()
		dec2 = func() decimal.Decimal {
			r, e := decimal.NewFromString(txt2)
			if e != nil {
				return decimal.Zero
			}
			return r
		}()
	} else if s.textformat != "" {
		txt1 = strings.Replace(s.textformat, s.textformat[:len(s.textformat)], "9", 1)
	}

	return txt1, txt2, dec1, dec2
}

type Currency struct {
	ConversionRate decimal.Decimal
	Precision      int
	IsoCode        string
	Symbol         string
}

func NewCurrency(
	pn func() (decimal.Decimal, error),
	precision int,
	isoCode string,
	symbol string,
) Currency {

	return Currency{
		ConversionRate: func() decimal.Decimal {
			CR, err := pn()
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
			return CR
		}(),
		Precision: precision,
		IsoCode:   isoCode,
		Symbol:    symbol,
	}

}

var currencies = CreateCurrencies()

const defaultCurrency = "PLN"
const sessionCurrency = "EUR"

func CreateCurrencies() map[string]Currency {

	CzeskaKuruna := NewCurrency(
		func() (decimal.Decimal, error) {
			return decimal.NewFromString("5.887000")
		},
		0,
		"CZK",
		"Kč",
	)

	FuntSzterling := NewCurrency(
		func() (decimal.Decimal, error) {
			return decimal.NewFromString("0.198264")
		},
		2,
		"GBP",
		"£",
	)

	Euro := NewCurrency(
		func() (decimal.Decimal, error) {
			return decimal.NewFromString("0.198264")
		},
		2,
		"EUR",
		"€",
	)

	PolskiZloty := NewCurrency(
		func() (decimal.Decimal, error) {
			return decimal.NewFromInt(1), nil
		},
		2,
		"PLN",
		"zł",
	)

	currencies := map[string]Currency{}
	currencies["CZK"] = CzeskaKuruna
	currencies["GBP"] = FuntSzterling
	currencies["EUR"] = Euro
	currencies["PLN"] = PolskiZloty

	return currencies
}
