package shared_types

import (
	"fmt"
	"sync"

	"github.com/shopspring/decimal"
)

type Candle interface {
	GetCandle() Candle
	GetStartTime() UnixTime
	GetEndTime() UnixTime
	GetHigh() decimal.Decimal
	SetHigh(num decimal.Decimal)
	GetLow() decimal.Decimal
	SetLow(num decimal.Decimal)
	GetClose() decimal.Decimal
	SetClose(num decimal.Decimal)
	GetVWAP() decimal.Decimal
	SetVWAP(num decimal.Decimal)
	GetVolume() decimal.Decimal
	SetVolume(num decimal.Decimal)
	GetCount() int
	SetCount(num int)
}

type CandleData struct {
	StartTime UnixTime
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	Trades    int
}

type Node struct {
	Data Candle
	Next *Node
}

type List struct {
	Head   *Node
	Last   *Node
	Length uint
}

func (l *List) AddToEnd(n *Node) {
	if l.Head == nil {
		l.Head = n
		l.Last = n
		l.Length++
		return
	}
	tmp := l.Last
	tmp.Next = n
	l.Last = n
	l.Length++
}

func (l List) Print(locker *sync.RWMutex) {
	locker.RLock()
	tmp := l.Head
	for tmp != nil {
		fmt.Println(string("\033[34m"), tmp.Data, string("\033[0m"))
		tmp = tmp.Next
	}
	locker.RUnlock()
	fmt.Println()
}

func (l List) Empty() bool {
	return l.Length == 0
}

/*func (l List) PrintRW(locker *sync.RWMutex) {
	defer locker.RUnlock()
	locker.RLock()
	tmp := l.Head
	for tmp != nil {
		fmt.Println(string("\033[34m"), &tmp, string("\033[0m"))
		tmp = tmp.Next
	}
	fmt.Println()
}*/
