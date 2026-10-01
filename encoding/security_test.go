package encoding_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/faustbrian/go-international/currency"
	"github.com/faustbrian/go-money/v2"
	moneyencoding "github.com/faustbrian/go-money/v2/encoding"
)

func TestSQLScanRejectsOversizedInputBeforeCopying(t *testing.T) {
	euro, err := currency.Parse("EUR")
	if err != nil {
		t.Fatal(err)
	}
	context, err := money.DefaultContext(euro)
	if err != nil {
		t.Fatal(err)
	}
	original, err := money.Parse("1.00", euro, context)
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := moneyencoding.MarshalText(original)
	if err != nil {
		t.Fatal(err)
	}
	exactLimit := string(canonical) + strings.Repeat(" ", moneyencoding.MaxEncodedBytes-len(canonical))
	for _, test := range []struct {
		name   string
		source any
	}{
		{name: "bytes", source: []byte(exactLimit)},
		{name: "string", source: exactLimit},
	} {
		t.Run("inclusive "+test.name, func(t *testing.T) {
			value := moneyencoding.SQLMoney{}
			if err := value.Scan(test.source); err != nil {
				t.Fatalf("Scan(exact limit %s) error = %v", test.name, err)
			}
			if equal, err := value.Money.Equal(original); err != nil || !equal {
				t.Fatalf("Scan(exact limit %s) = %s, %v", test.name, value.Money, err)
			}
		})
	}
	for _, source := range []any{[]byte(exactLimit + " "), exactLimit + " ", nil} {
		value := moneyencoding.SQLMoney{Money: original}
		if err := value.Scan(source); !errors.Is(err, moneyencoding.ErrInvalidEncoding) {
			t.Fatalf("Scan(%T over limit or nil) error = %v", source, err)
		}
		if equal, err := value.Money.Equal(original); err != nil || !equal {
			t.Fatalf("Scan(%T over limit or nil) changed receiver", source)
		}
	}

	oversized := strings.Repeat("x", 64*1024)
	for _, test := range []struct {
		name   string
		source any
	}{
		{name: "bytes", source: []byte(oversized)},
		{name: "string", source: oversized},
	} {
		t.Run(test.name, func(t *testing.T) {
			value := moneyencoding.SQLMoney{Money: original}
			if err := value.Scan(test.source); !errors.Is(err, moneyencoding.ErrInvalidEncoding) {
				t.Fatalf("Scan(oversized %s) error = %v", test.name, err)
			}
			if equal, err := value.Money.Equal(original); err != nil || !equal {
				t.Fatalf("Scan(oversized %s) changed receiver", test.name)
			}

			result := testing.Benchmark(func(b *testing.B) {
				for range b.N {
					if err := value.Scan(test.source); !errors.Is(err, moneyencoding.ErrInvalidEncoding) {
						b.Fatalf("Scan(oversized %s) error = %v", test.name, err)
					}
				}
			})
			if result.AllocedBytesPerOp() > 4*moneyencoding.MaxEncodedBytes {
				t.Errorf("Scan(oversized %s) allocated %d bytes per call before rejection", test.name, result.AllocedBytesPerOp())
			}
		})
	}
}
