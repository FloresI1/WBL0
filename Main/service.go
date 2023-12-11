package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func consumeDeliveries(ctx context.Context) {
	// Подключение к серверу NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Подписка на тему "your-nats-channel"
	subscription, err := nc.Subscribe("your-nats-channel", func(msg *nats.Msg) {
		var pars Combo
		err := json.Unmarshal(msg.Data, &pars)
		if err != nil {
			log.Printf("Error decoding delivery: %v", err)
			return
		}
		saveDeliveryToCacheAndDB(pars)
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Subscribed to channel: %s\n", subscription.Subject)

	<-ctx.Done()
	log.Println("Shutting down delivery processing...")
}

func saveDeliveryToCacheAndDB(combo Combo) {
	log.Printf("Saving delivery to cache: %+v\n", combo)

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Получаем существующий срез для orderUID или создаем новый
	combos, ok := cache.Load(combo.OrderUID)
	if !ok {
		combos = []*Combo{}
	}

	existingCombos, ok := combos.([]*Combo)
	if !ok {
		log.Printf("Invalid data format in cache for OrderUID: %s\n", combo.OrderUID)
		return
	}

	// Добавляем новое значение в срез
	combos = append(existingCombos, &combo)

	combosSlice, ok := combos.([]*Combo)
	if !ok {
		log.Printf("Invalid data format in cache for OrderUID: %s\n", combo.OrderUID)
		return
	}

	// Обновляем значение в кеше
	cache.Store(combo.OrderUID, combosSlice)

	log.Printf("Saved combo to cache: %+v\n", combosSlice)

	 saveToDB(combo)
}
	
func restoreCacheFromDB() error {
	tables := []string{"delivery", "payment", "items", "info"}

	for _, tableName := range tables {
		var query string
		var columns string
		var dest []interface{}

		switch tableName {
		case "delivery":
			columns = "order_uid, name, phone, zip, city, address, region, email"
			dest = make([]interface{}, 8)
			dest[0], dest[1], dest[2], dest[3], dest[4], dest[5], dest[6], dest[7] = new(string), new(string), new(string), new(string), new(string), new(string), new(string), new(string)
		case "payment":
			columns = "order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee"
			dest = make([]interface{}, 11)
			dest[0], dest[1], dest[2], dest[3], dest[4], dest[5], dest[6], dest[7], dest[8], dest[9], dest[10] = new(string), new(string), new(string), new(string), new(string), new(int), new(time.Time), new(string), new(int), new(int), new(int)
		case "items":
			columns = "order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status"
			dest = make([]interface{}, 12)
			dest[0], dest[1], dest[2], dest[3], dest[4], dest[5], dest[6], dest[7], dest[8], dest[9], dest[10], dest[11] = new(string), new(int), new(string), new(int), new(string), new(string), new(int), new(string), new(int), new(int), new(string), new(int)
		case "info":
			columns = "order_uid, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard"
			dest = make([]interface{}, 9)
			dest[0], dest[1], dest[2], dest[3], dest[4], dest[5], dest[6], dest[7], dest[8] = new(string), new(string), new(string), new(string), new(string), new(string), new(string), new(time.Time), new(string)
		}

		query = "SELECT " + columns + " FROM " + tableName
		rows, err := db.Query(query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var orderUID, name, phone, zip, city, address, region, email string
			var transaction, requestID, currency, provider, bank string
			var amount, deliveryCost, goodsTotal, customFee, chrtID, sale, totalPrice, nmID, status, price, smID int
			var paymentDt, dateCreated time.Time
			var trackNumber, rid, itemName, size, brand string
			var locale, internalSignature, customerID, deliveryService, shardKey,oofShard string

			// Сканирование строки в переменные
			switch tableName {
			case "delivery":
				if err := rows.Scan(&orderUID, &name, &phone, &zip, &city, &address, &region, &email); err != nil {
					log.Printf("Error scanning row from %s table: %v", tableName, err)
					continue
				}
			case "payment":
				if err := rows.Scan(&orderUID, &transaction, &requestID, &currency, &provider, &amount, &paymentDt, &bank, &deliveryCost, &goodsTotal, &customFee); err != nil {
					log.Printf("Error scanning row from %s table: %v", tableName, err)
					continue
				}
			case "items":
				if err := rows.Scan(&orderUID, &chrtID, &trackNumber, &price, &rid, &itemName, &sale, &size, &totalPrice, &nmID, &brand, &status); err != nil {
					log.Printf("Error scanning row from %s table: %v", tableName, err)
					continue
				}
			case "info":
				if err := rows.Scan(&orderUID, &locale, &internalSignature, &customerID, &deliveryService, &shardKey, &smID, &dateCreated, &oofShard); err != nil {
					log.Printf("Error scanning row from %s table: %v", tableName, err)
					continue
				}
			}
		}
	}

	return nil
}
