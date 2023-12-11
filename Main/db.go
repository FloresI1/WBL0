package main

import (
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/lib/pq"
)

func init() {
	var err error
	// Подключение к базе данных PostgreSQL
	db, err = sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable password=1")
	if err != nil {
		log.Fatal(err)
	}

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database")

	// Восстановление кэша из БД
	err = restoreCacheFromDB()
	if err != nil {
		log.Printf("Error restoring cache from DB: %v", err)
	}
}
	// Сохраняем данные в БД
func saveToDB(combo Combo) {
	_, err := db.Exec(
		"INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		combo.Delivery.OrderUID, combo.Delivery.Name, combo.Delivery.Phone, combo.Delivery.Zip, combo.Delivery.City, combo.Delivery.Address, combo.Delivery.Region, combo.Delivery.Email,
	)
	if err != nil {
		log.Printf("Error saving delivery to DB: %v", err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		combo.Payment.OrderUID, combo.Payment.Transaction, combo.Payment.RequestID, combo.Payment.Currency, combo.Payment.Provider, combo.Payment.Amount, combo.Payment.PaymentDt, combo.Payment.Bank, combo.Payment.DeliveryCost, combo.Payment.GoodsTotal, combo.Payment.CustomFee,
	)
	if err != nil {
		log.Printf("Error saving payment to DB: %v", err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		combo.Items.OrderUID, combo.Items.ChrtID, combo.Items.TrackNumber, combo.Items.Price, combo.Items.RID, combo.Items.Name, combo.Items.Sale, combo.Items.Size, combo.Items.TotalPrice, combo.Items.NmID, combo.Items.Brand, combo.Items.Status,
	)
	if err != nil {
		log.Printf("Error saving items to DB: %v", err)
		return
	}

	_, err = db.Exec(
		"INSERT INTO info (order_uid, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		combo.Info.OrderUID, combo.Info.Locale, combo.Info.InternalSignature, combo.Info.CustomerID, combo.Info.DeliveryService, combo.Info.ShardKey, combo.Info.SMID, combo.Info.DateCreated, combo.Info.OOFShard,
	)
	if err != nil {
		log.Printf("Error saving info to DB: %v", err)
		return
	}

	log.Println("Saved combo to DB")
}

func fetchDataFromDB(orderUID string) ([]*Combo, error) {
    // Инициализируем данные для хранения полученных данных
    orderData := &Combo{
        Delivery: Delivery{},
        Info:     Info{},
        Items:    Items{},
        Payment:  Payment{},
    }

    err := db.QueryRow("SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1", orderUID).
        Scan(&orderData.Delivery.Name, &orderData.Delivery.Phone, &orderData.Delivery.Zip, &orderData.Delivery.City,
            &orderData.Delivery.Address, &orderData.Delivery.Region, &orderData.Delivery.Email)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении данных из таблицы delivery: %v", err)
    }

    err = db.QueryRow("SELECT locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM info WHERE order_uid = $1", orderUID).
        Scan(&orderData.Info.Locale, &orderData.Info.InternalSignature, &orderData.Info.CustomerID, &orderData.Info.DeliveryService,
            &orderData.Info.ShardKey, &orderData.Info.SMID, &orderData.Info.DateCreated, &orderData.Info.OOFShard)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении данных из таблицы info: %v", err)
    }

    err = db.QueryRow("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE order_uid = $1", orderUID).
        Scan(&orderData.Items.ChrtID, &orderData.Items.TrackNumber, &orderData.Items.Price, &orderData.Items.RID, &orderData.Items.Name,
            &orderData.Items.Sale, &orderData.Items.Size, &orderData.Items.TotalPrice, &orderData.Items.NmID, &orderData.Items.Brand, &orderData.Items.Status)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении данных из таблицы items: %v", err)
    }

    err = db.QueryRow("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1", orderUID).
        Scan(&orderData.Payment.Transaction, &orderData.Payment.RequestID, &orderData.Payment.Currency,
            &orderData.Payment.Provider, &orderData.Payment.Amount, &orderData.Payment.PaymentDt, &orderData.Payment.Bank,
            &orderData.Payment.DeliveryCost, &orderData.Payment.GoodsTotal, &orderData.Payment.CustomFee)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении данных из таблицы payment: %v", err)
    }

    // Формируем срез данных и возвращаем его
    return []*Combo{orderData}, nil
}