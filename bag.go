package money

import (
	"cmp"
	"sort"
	"strings"

	"github.com/faustbrian/go-international/currency"
)

// MaxMoneyBagEntries bounds input, heterogeneous output, and lookup work.
const MaxMoneyBagEntries = 1_000

// MoneyBag is an immutable deterministic collection keyed by exact currency
// and context identity. Its zero value is an empty bag.
type MoneyBag struct{ values []Money }

// NewMoneyBag validates and combines values with identical identities. The
// input count must not exceed MaxMoneyBagEntries even when all inputs combine;
// excess returns ErrMoneyBagLimit. Empty input produces an empty bag. Inputs
// are not retained, and invalid values or sum-limit failures return no partial bag.
func NewMoneyBag(values ...Money) (MoneyBag, error) {
	if len(values) > MaxMoneyBagEntries {
		return MoneyBag{}, ErrMoneyBagLimit
	}
	bag := MoneyBag{}
	var err error
	for _, value := range values {
		bag, err = bag.Add(value)
		if err != nil {
			return MoneyBag{}, err
		}
	}

	return bag, nil
}

// Add returns a new bag, combining only an identical currency/context entry.
func (bag MoneyBag) Add(value Money) (MoneyBag, error) {
	if !value.Valid() {
		return MoneyBag{}, ErrInvalidMoney
	}
	values := append([]Money(nil), bag.values...)
	for index, existing := range values {
		if existing.currency == value.currency && existing.context == value.context {
			combined, err := existing.Add(value)
			if err != nil {
				return MoneyBag{}, err
			}
			values[index] = combined
			return MoneyBag{values: values}, nil
		}
	}
	if len(values) >= MaxMoneyBagEntries {
		return MoneyBag{}, ErrMoneyBagLimit
	}
	values = append(values, value)
	sort.Slice(values, func(left, right int) bool {
		if values[left].currency != values[right].currency {
			return strings.Compare(values[left].currency.String(), values[right].currency.String()) == -1
		}
		return contextLess(values[left].context, values[right].context)
	})

	return MoneyBag{values: values}, nil
}

// Values returns an independent slice ordered by currency code, then context
// kind, scale, cash step and currency binding. An empty bag returns nil.
func (bag MoneyBag) Values() []Money { return append([]Money(nil), bag.values...) }

// Get returns the value with exact currency and context identity. An absent
// entry returns invalid Money and false; no conversion or context inference occurs.
func (bag MoneyBag) Get(code currency.Code, context Context) (Money, bool) {
	for _, value := range bag.values {
		if value.currency == code && value.context == context {
			return value, true
		}
	}

	return Money{}, false
}

func contextLess(left, right Context) bool {
	if comparison := cmp.Compare(left.kind, right.kind); comparison != 0 {
		return comparison == -1
	}
	if comparison := cmp.Compare(left.scale, right.scale); comparison != 0 {
		return comparison == -1
	}
	if comparison := cmp.Compare(left.cashStep, right.cashStep); comparison != 0 {
		return comparison == -1
	}

	return strings.Compare(left.currency.String(), right.currency.String()) == -1
}
