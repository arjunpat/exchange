package main

import (
	"testing"
)

func TestPlaceBuy(t *testing.T) {
	var ts []*Transaction

	b := NewOrderBook("GOOG")
	b.PlaceBuy(1, 100, "David", 1)
	b.PlaceBuy(5, 100, "Brian", 2)
	b.PlaceBuy(20, 100, "Arjun", 3)
	b.PlaceBuy(4, 101, "Kevin", 4)
	if b.buys.Len() != 4 {
		t.Fatal("Length incorrect")
	}

	b.PlaceSell(2, 101, "Andrew", 5)
	ts = append(ts, &Transaction{"Andrew", "Kevin", "GOOG", 2, 101, -1})

	b.PlaceSell(4, 100, "Bob", 6)
	ts = append(ts,
		&Transaction{"Bob", "Kevin", "GOOG", 2, 100, -1},
		&Transaction{"Bob", "David", "GOOG", 1, 100, -1},
		&Transaction{"Bob", "Brian", "GOOG", 1, 100, -1},
	)
	if b.MarketPrice() != 100 {
		t.Fatal("Market price failed. Expected: 100 Got:", b.MarketPrice())
	}

	b.PlaceSell(5, 98, "Jake", 7)
	ts = append(ts,
		&Transaction{"Jake", "Brian", "GOOG", 4, 98, -1},
		&Transaction{"Jake", "Arjun", "GOOG", 1, 98, -1},
	)
	if b.MarketPrice() != 98 {
		t.Fatal("Market price failed. Expected: 98 Got:", b.MarketPrice())
	}

	for i := range ts {
		if !transEqual(b.transactions[i], ts[i]) {
			t.Fatal("Failed on ", i, "Expected:", ts[i], "Got:", b.transactions[i])
		}
	}
}

func TestPlaceSell(t *testing.T) {
	var ts []*Transaction

	b := NewOrderBook("GOOG")
	b.PlaceSell(1, 100, "David", 1)
	b.PlaceSell(5, 100, "Brian", 2)
	b.PlaceSell(20, 100, "Arjun", 3)
	b.PlaceSell(4, 99, "Kevin", 4)
	if b.sells.Len() != 4 {
		t.Fatal("Length incorrect")
	}

	b.PlaceBuy(2, 99, "Andrew", 5)
	ts = append(ts, &Transaction{"Kevin", "Andrew", "GOOG", 2, 99, -1})

	b.PlaceBuy(4, 100, "Bob", 6)
	ts = append(ts,
		&Transaction{"Kevin", "Bob", "GOOG", 2, 100, -1},
		&Transaction{"David", "Bob", "GOOG", 1, 100, -1},
		&Transaction{"Brian", "Bob", "GOOG", 1, 100, -1},
	)
	if b.MarketPrice() != 100 {
		t.Fatal("Market price failed. Expected: 100 Got:", b.MarketPrice())
	}

	b.PlaceBuy(5, 103, "Jake", 7)
	ts = append(ts,
		&Transaction{"Brian", "Jake", "GOOG", 4, 103, -1},
		&Transaction{"Arjun", "Jake", "GOOG", 1, 103, -1},
	)
	if b.MarketPrice() != 103 {
		t.Fatal("Market price failed. Expected: 98 Got:", b.MarketPrice())
	}

	for i := range ts {
		if !transEqual(b.transactions[i], ts[i]) {
			t.Fatal("Failed on ", i, "Expected:", ts[i], "Got:", b.transactions[i])
		}
	}
}

func transEqual(a *Transaction, b *Transaction) bool {
	return a.Security == b.Security && a.Amt == b.Amt && a.Price == b.Price && a.From == b.From && a.To == a.To
}
