package exchange

import (
	"container/heap"
	"fmt"
	"time"
)

type Order struct {
	CreatedAt int64
	Creator   string
	Amt       int
	Price     float64
	Buy       bool
}

func (o *Order) ReduceAmt(amt int) int {
	o.Amt -= amt
	return o.Amt
}

type Transaction struct {
	From     string
	To       string
	Security string
	Amt      int
	Price    float64
	Time     int64
}

func (t *Transaction) String() string {
	//amt := strconv.Itoa(t.Amt)
	//price := strconv.FormatFloat(t.Price, 'E', -1, 64)
	//time := strconv.FormatFloat(t.Time, 'E', -1, 64)
	//return t.From + " -> " + t.To + ": " + amt + " " + t.Security + " @ $" + price + " at time " + time
	return fmt.Sprint(t.From, " -> ", t.To, ": ", t.Amt, " ", t.Security, " @ $", t.Price, " at time ", t.Time)
}

type OrderPQ []*Order

func (pq OrderPQ) Len() int      { return len(pq) }
func (pq OrderPQ) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq OrderPQ) Less(i, j int) bool {
	if pq[i].Price != pq[j].Price {
		if pq[i].Buy {
			if pq[i].Price < pq[j].Price {
				return false
			}
		} else if pq[i].Price > pq[j].Price {
			return false
		}
		return true
	}

	return pq[i].CreatedAt < pq[j].CreatedAt
}

func (pq *OrderPQ) Push(o interface{}) {
	*pq = append(*pq, o.(*Order))
}

func (pq *OrderPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return x
}

type OrderBook struct {
	security     string
	buys         OrderPQ
	sells        OrderPQ
	transactions []*Transaction
}

func NewOrderBook(security string) OrderBook {
	buy := make(OrderPQ, 0)
	sell := make(OrderPQ, 0)
	heap.Init(&buy)
	heap.Init(&sell)
	return OrderBook{
		security:     security,
		buys:         buy,
		sells:        sell,
		transactions: make([]*Transaction, 0),
	}
}

func MinInt(a int, b int) int {
	if b < a {
		return b
	}
	return a
}

func (b *OrderBook) PlaceBuy(amt int, price float64, creator string, createdAt int64) {
	for b.sells.Len() > 0 && b.sells[0].Price <= price && amt > 0 {
		top := b.sells[0]
		amtMatched := MinInt(top.Amt, amt)
		amt -= amtMatched

		if top.ReduceAmt(amtMatched) == 0 {
			heap.Pop(&b.sells)
		}

		b.transactions = append(b.transactions, &Transaction{
			From:     top.Creator,
			To:       creator,
			Security: b.security,
			Amt:      amtMatched,
			Price:    price,
			Time:     time.Now().UnixMilli(),
		})
	}

	if amt != 0 {
		heap.Push(&b.buys, &Order{
			CreatedAt: createdAt,
			Creator:   creator,
			Amt:       amt,
			Price:     price,
			Buy:       true,
		})
	}
}

func (b *OrderBook) PlaceSell(amt int, price float64, creator string, createdAt int64) {
	for b.buys.Len() > 0 && b.buys[0].Price >= price && amt > 0 {
		top := b.buys[0]
		amtMatched := MinInt(top.Amt, amt)
		amt -= amtMatched

		if top.ReduceAmt(amtMatched) == 0 {
			heap.Pop(&b.buys)
		}

		b.transactions = append(b.transactions, &Transaction{
			From:     creator,
			To:       top.Creator,
			Security: b.security,
			Amt:      amtMatched,
			Price:    price,
			Time:     time.Now().UnixMilli(),
		})
	}

	if amt != 0 {
		heap.Push(&b.sells, &Order{
			CreatedAt: createdAt,
			Creator:   creator,
			Amt:       amt,
			Price:     price,
			Buy:       false,
		})
	}
}

func (b *OrderBook) MarketPrice() float64 {
	if len(b.transactions) == 0 {
		return -1
	}
	return b.transactions[len(b.transactions)-1].Price
}
