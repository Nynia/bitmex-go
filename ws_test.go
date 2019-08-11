package bitmex

import (
	"fmt"
	"log"
	"testing"
)

func TestBitMEXConnect(t *testing.T) {
	b := New(HostTestnet, "", "")
	//b.SetHttpProxy("http://127.0.0.1:1080")
	subscribeInfos := []SubscribeInfo{
		{Op: BitmexWSOrderBookL2, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(BitmexWSOrderBookL2, func(ob OrderBookDataL2, symbol string) {
		m := ob.OrderBook()
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	})

	b.StartWS()

	select {}
}

func TestBitMEXWS(t *testing.T) {
	b := New(HostTestnet, "", "")
	b.SetHttpProxy("http://127.0.0.1:1080")
	subscribeInfos := []SubscribeInfo{
		{Op: BitmexWSOrderBookL2, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}

	b.On(BitmexWSOrderBookL2, func(ob OrderBookDataL2, symbol string) {
		m := ob.OrderBook()
		fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", m.Asks[0], m.Bids[0])
	})

	b.StartWS()

	select {}
}
