package main

import (
	"github.com/sumorf/bitmex-api"
	"github.com/sumorf/bitmex-api/swagger"
	"log"
)

func main() {
	b := bitmex.New(bitmex.HostTestnet, "DEt3D-w0hMPCGwqx3MW0jlU8", "glojC1IZXq94N3MlB8CLM3qeW7cjSUxzqbYRdkz9jJMV0p8q")
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSOrderBookL2, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSOrder, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSPosition, Param: "XBTUSD"},
		{Op: bitmex.BitmexWSMargin, Param: "XBTUSD"},
	}
	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}
	b.On(bitmex.BitmexWSOrderBookL2, func(m bitmex.OrderBookDataL2, symbol string) {
		//ob := m.OrderBook()
		//fmt.Printf("\rOrderbook Asks: %#v Bids: %#v                            ", ob.Asks[0], ob.Bids[0])
	}).On(bitmex.BitmexWSOrder, func(m []*swagger.Order, action string) {
		ord := m[0]
		if ord.OrdStatus == "Filled" {
			if ord.Side == bitmex.SIDE_BUY {
				go b.SendOrder(ord.Symbol, "Sell", ord.OrderQty, ord.Price+10)
			} else {
				go b.SendOrder(ord.Symbol, "Buy", ord.OrderQty, ord.Price-10)
			}
		}
		//fmt.Printf("Order action=%v orders=%#v\n", action, m[0].OrderID)
	}).On(bitmex.BitmexWSPosition, func(m []*swagger.Position, action string) {
		//fmt.Printf("Position action=%v positions=%#v\n", action, m)
	}).On(bitmex.BitmexWSMargin, func(m []*swagger.Margin, action string) {
		//fmt.Printf("Wallet action=%v margins=%#v\n", action, m)
	})
	b.StartWS()
	//fmt.Printf(b.GetTmpOrder("XBTUSD", "Buy"))
	//b.SendOrder("XBTUSD", "Buy", 1000, 12500)
	// Get orderbook by rest api
	//b.GetOrderBook(10, "XBTUSD")
	//// Place a limit buy order
	//b.PlaceOrder(bitmex.SIDE_BUY, bitmex.ORD_TYPE_LIMIT, 0, 6000.0, 1000, "", "", "XBTUSD")
	//b.GetOrders("XBTUSD")
	//b.GetOrder("{OrderID}", "XBTUSD")
	//b.AmendOrder("{OrderID}", 6000.5)
	//b.CancelOrder("{OrderID}")
	//b.CancelAllOrders("XBTUSD")
	//b.GetPosition("XBTUSD")
	//b.GetMargin()

	b.SyncOrders()

	forever := make(chan bool)
	<-forever
}
