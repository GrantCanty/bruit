package shared_types

import (
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type Node struct {
	Data Candle
	Next *Node
}

type List struct {
	Head   *Node
	Last   *Node
	Length uint
}

func NewList(head *Node, last *Node, length uint) *List {
	return &List{Head: head, Last: last, Length: length}
}

func (l *List) GetList() *List {
	return l
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

	if l.IsEmpty() {
		return
	}

	tmp := l.Head
	for tmp != nil {
		fmt.Println(string("\033[34m"), tmp.Data, string("\033[0m"))
		tmp = tmp.Next
	}
	locker.RUnlock()
	fmt.Println()
}

func (l List) IsEmpty() bool {
	return l.Length == 0
}

func (l *List) GetLast() *Node {
	return l.Last
}

func (l *List) EditCandle(oldCandle Candle, newCandle Candle) {
	oldCandle.SetHigh(newCandle.GetHigh())
	oldCandle.SetLow(newCandle.GetLow())
	oldCandle.SetClose(newCandle.GetClose(), newCandle.GetVolume())
	oldCandle.SetVWAP(newCandle.GetVWAP(), newCandle.GetVolume())
	oldCandle.SetCount(newCandle.GetCount(), newCandle.GetVolume())
	oldCandle.SetVolume(newCandle.GetVolume())
}

func (l *List) AddCandle(newCandle Candle, emptyCandles Candle, interval int) { // old candle should switch to list
	since := time.Since(l.GetLast().Data.GetEndTime().Time).Minutes()
	if since < time.Duration(interval).Minutes() { // if the time since the close of the last candle is less than the time of the connection's interval, the candle you received will just be added to the end
		newCandle.SetStartTime(l.GetLast().Data.GetEndTime().Time)
		node := Node{Data: newCandle, Next: nil}
		l.AddToEnd(&node)
	} else {
		newNodeCount := int(int(since) / interval)
		zero := decimal.New(0, 0)

		for i := int(0); i < newNodeCount; i++ {
			last := l.GetLast()
			close := last.Data.GetClose()

			emptyCandles.SetCandle(last.Data.GetStartTime(), last.Data.GetEndTime(), close, close, close, close, zero, zero, 0)
			emptyCandles.SetStartTime(emptyCandles.GetStartTime().Time.Add(time.Minute * time.Duration(interval)))
			emptyCandles.SetEndTime(emptyCandles.GetEndTime().Time.Add(time.Minute * time.Duration(interval)))

			node := Node{Data: emptyCandles, Next: nil}
			l.AddToEnd(&node)
		}
		node := Node{Data: newCandle, Next: nil}
		l.AddToEnd(&node)
	}
}
