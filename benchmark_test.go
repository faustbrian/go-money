package money

import (
	"context"
	"math/big"
	"testing"

	rhymond "github.com/Rhymond/go-money"
	"github.com/faustbrian/go-international/currency"
	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/decimal"
	"github.com/faustbrian/go-math/integer"
	govalues "github.com/govalues/money"
)

func BenchmarkExactAdd(b *testing.B) {
	euro, _ := currency.Parse("EUR")
	monetaryContext, _ := DefaultContext(euro)
	left, _ := Parse("123.45", euro, monetaryContext)
	right, _ := Parse("67.89", euro, monetaryContext)
	assertLocalAddBenchmark(b, left, right)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = left.Add(right)
	}
}

func BenchmarkGovaluesAdd(b *testing.B) {
	left, _ := govalues.ParseAmount("EUR", "123.45")
	right, _ := govalues.ParseAmount("EUR", "67.89")
	result, err := left.Add(right)
	if err != nil || result.String() != "EUR 191.34" {
		b.Fatalf("correctness gate = %s, %v", result, err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = left.Add(right)
	}
}

func BenchmarkRhymondAdd(b *testing.B) {
	left := rhymond.New(12_345, rhymond.EUR)
	right := rhymond.New(6_789, rhymond.EUR)
	result, err := left.Add(right)
	if err != nil || result.Amount() != 19_134 {
		b.Fatalf("correctness gate = %d, %v", result.Amount(), err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = left.Add(right)
	}
}

func BenchmarkRejectOversizedTypedInputs(b *testing.B) {
	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 10_000
	limits.MaxOutputDigits = 10_000
	limits.MaxIntermediateBits = 20_000
	hugeCoefficient := new(big.Int).Lsh(big.NewInt(1), 10_000)
	hugeInteger, err := integer.FromBig(hugeCoefficient, limits)
	if err != nil {
		b.Fatal(err)
	}
	hugeDecimal, err := decimal.FromBig(hugeCoefficient, -2, limits)
	if err != nil {
		b.Fatal(err)
	}
	euro, _ := currency.Parse("EUR")
	monetaryContext, _ := DefaultContext(euro)
	value, _ := Parse("1.00", euro, monetaryContext)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = AmountFromDecimal(hugeDecimal)
		_, _ = FromMinorUnits(hugeInteger, euro, monetaryContext)
		_, _ = value.Allocate(context.Background(), []integer.Integer{hugeInteger})
	}
}

func assertLocalAddBenchmark(b *testing.B, left, right Money) {
	b.Helper()
	result, err := left.Add(right)
	if err != nil || result.String() != "191.34 EUR" {
		b.Fatalf("correctness gate = %s, %v", result, err)
	}
}
