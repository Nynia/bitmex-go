package main

import (
	"fmt"
	"github.com/Nynia/bitmex-go"
	"github.com/sumorf/bitmex-api/swagger"
	"log"
	"strings"
)

const (
	INIT_POSITION int     = 61
	PRICE_DIST    float64 = 24
	PROFIT_DIST   float64 = 12
	UNIT_AMOUNT   int     = 375
	SYMBOL1       string  = "XBTUSD"
	SYMBOL2       string  = "XBTZ19"
)

func main() {
	b := bitmex.New(bitmex.HostTestnet, "JbGVj9i92-ChCyLPTPqcAQMW", "MEcZ1CKCouHPKj9MvRjuNyjSYGD16eGbhctntULdzrnG9x9k")
	subscribeInfos := []bitmex.SubscribeInfo{
		{Op: bitmex.BitmexWSOrder, Param: SYMBOL1},
		{Op: bitmex.BitmexWSOrder, Param: SYMBOL2},
	}

	err := b.Subscribe(subscribeInfos)
	if err != nil {
		log.Fatal(err)
	}
	buy_amount := 0
	sell_amount := 0
	total := 0
	b.On(bitmex.BitmexWSOrder, func(m []*swagger.Order, action string) {
		ord := m[0]
		if ord.OrdStatus == "Filled" {
			log.Printf("side: %s, symbol: %s, cum_qty: %.2f, order_px: %.2f, orderID: %s",
				ord.Side, ord.Symbol, ord.OrderQty, ord.Price, ord.OrderID)
			if int(ord.OrderQty) != UNIT_AMOUNT {
				log.Printf("开仓：%s", ord.Symbol)
				new_orders := make([]string, INIT_POSITION)
				for i := 0; i < INIT_POSITION; i++ {
					price := ord.Price + float64(i)*PRICE_DIST + PROFIT_DIST
					new_orders[i] = fmt.Sprintf(`{"symbol":"%s","side":"Sell","orderQty":%d,"ordType":"Limit","price":%f}`, ord.Symbol, UNIT_AMOUNT, price)
				}
				ordstr := "[" + strings.Join(new_orders, ",") + "]"
				_, _ = b.NewBuckOrder(ordstr)
				for i := 0; i < INIT_POSITION; i++ {
					price := ord.Price - float64(i+1)*PRICE_DIST
					new_orders[i] = fmt.Sprintf(`{"symbol":"%s","side":"Buy","orderQty":%d,"ordType":"Limit","price":%f}`, ord.Symbol, UNIT_AMOUNT, price)
				}
				ordstr = "[" + strings.Join(new_orders, ",") + "]"
				_, _ = b.NewBuckOrder(ordstr)
			} else {
				if ord.Side == bitmex.SIDE_BUY {
					buy_amount++
					go b.SendOrder(ord.Symbol, "Sell", ord.OrderQty, ord.Price+PROFIT_DIST, total+1)
				} else {
					sell_amount++
					go b.SendOrder(ord.Symbol, "Buy", ord.OrderQty, ord.Price-PROFIT_DIST, total+1)
				}
				total += 1
				log.Printf("TOTAL: %d\tBUY: %d\tSELL: %d", total, buy_amount, sell_amount)
			}
		}
	})
	b.StartWS()
	_, _ = b.SyncOrders()
	go b.MonitorTmpOrder()

	forever := make(chan bool)
	<-forever
}
