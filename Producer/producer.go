package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

type Delivery struct {
	OrderUID string `json:"order_uid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
	OrderUID     string    `json:"order_uid"`
	Transaction  string    `json:"transaction"`
	RequestID    string    `json:"request_id"`
	Currency     string    `json:"currency"`
	Provider     string    `json:"provider"`
	Amount       int       `json:"amount"`
	PaymentDt    time.Time `json:"payment_dt"`
	Bank         string    `json:"bank"`
	DeliveryCost int       `json:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total"`
	CustomFee    int       `json:"custom_fee"`
}

type Items struct {
	OrderUID    string `json:"order_uid"`
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
type Info struct {
	OrderUID    string `json:"order_uid"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SMID              int    `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OOFShard          string    `json:"oof_shard"`
}
type DeliveryAndPaymentandItems struct {
	OrderUID string   `json:"order_uid"`
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    Items    `json:"items"`
	Info     Info     `json:"info"`
}

func produceDeliveries() {
	// Подключение к серверу NATS Streaming
	nc, err := nats.Connect("nats://0.0.0.0:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Обработка сигнала завершения
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Инициализация генератора случайных чисел
	randGenerator := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Генерация случайного количества сообщений от 1 до 10
	numMessages := randGenerator.Intn(100) + 1

	for i := 0; i < numMessages; i++ {
		orderUID := generateRandomOrderUID(randGenerator)

		delivery := generateRandomDelivery(randGenerator, orderUID)

		payment := generateRandomPayment(randGenerator, orderUID)

		items := generateRandomItems(randGenerator, orderUID)

		info := generateRandomInfo(randGenerator, orderUID)

		deliveryAndPayment := DeliveryAndPaymentandItems{
			OrderUID: orderUID,
			Delivery: delivery,
			Payment:  payment,
			Items:    items,
			Info: 	  info,
		}

		// Преобразование в JSON
		msgBytes, err := json.Marshal(deliveryAndPayment)
		if err != nil {
			log.Printf("Error encoding delivery and payment: %v", err)
			continue
		}

		// Отправка в канал NATS Streaming
		err = nc.Publish("your-nats-channel", msgBytes)
		if err != nil {
			log.Printf("Error publishing delivery and payment: %v", err)
			continue
		}

		// Подтверждение отправки
		log.Printf("Sent delivery and payment: %s\n", string(msgBytes))
	}

	// Ожидание сигнала завершения
	<-signalCh
}

func generateRandomOrderUID(randGenerator *rand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	orderUID := make([]byte, 16)
	for i := range orderUID {
		orderUID[i] = charset[randGenerator.Intn(len(charset))]
	}
	return string(orderUID)
}

func generateRandomDelivery(randGenerator *rand.Rand, orderUID string) Delivery {
	return Delivery{
		OrderUID: orderUID,
		Name:     generateRandomString(randGenerator, 10),
		Phone:    generateRandomString(randGenerator, 10),
		Zip:      generateRandomString(randGenerator, 5),
		City:     generateRandomString(randGenerator, 8),
		Address:  generateRandomString(randGenerator, 15),
		Region:   generateRandomString(randGenerator, 8),
		Email:    generateRandomEmail(randGenerator),
	}
}

func generateRandomPayment(randGenerator *rand.Rand, orderUID string) Payment {
	return Payment{
		OrderUID:     orderUID,
		Transaction:  generateRandomString(randGenerator, 10),
		RequestID:    generateRandomString(randGenerator, 8),
		Currency:     "USD",
		Provider:     "Stripe",
		Amount:       randGenerator.Intn(100) + 1,
		PaymentDt:    time.Now(),
		Bank:         "Bank of Example",
		DeliveryCost: randGenerator.Intn(20) + 1,
		GoodsTotal:   randGenerator.Intn(500) + 1,
		CustomFee:    randGenerator.Intn(10) + 1,
	}
}

func generateRandomItems(randGenerator *rand.Rand, orderUID string) Items {
	return Items{
		OrderUID:    orderUID,
		ChrtID:      randGenerator.Intn(100) + 1,
		TrackNumber: generateRandomString(randGenerator, 8),
		Price:       randGenerator.Intn(50) + 1,
		RID:         generateRandomString(randGenerator, 6),
		Name:        generateRandomString(randGenerator, 10),
		Sale:        randGenerator.Intn(20) + 1,
		Size:        generateRandomString(randGenerator, 5),
		TotalPrice:  randGenerator.Intn(100) + 1,
		NmID:        randGenerator.Intn(50) + 1,
		Brand:       generateRandomString(randGenerator, 8),
		Status:      randGenerator.Intn(5) + 1,
	}
}
func generateRandomInfo(randGenerator *rand.Rand, orderUID string) Info {
	return Info{
		OrderUID:    orderUID,
		Locale:            generateRandomString(randGenerator, 5),
		InternalSignature: generateRandomString(randGenerator, 10),
		CustomerID:        generateRandomString(randGenerator, 8),
		DeliveryService:   generateRandomString(randGenerator, 6),
		ShardKey:          generateRandomString(randGenerator, 8),
		SMID:              randGenerator.Intn(6),
		DateCreated:       time.Now(),
		OOFShard:          generateRandomString(randGenerator, 8),
	}
}
func generateRandomString(randGenerator *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[randGenerator.Intn(len(charset))]
	}
	return string(result)
}

func generateRandomEmail(randGenerator *rand.Rand) string {
	return generateRandomString(randGenerator, 8) + "@example.com"
}

func main() {
	// Запуск функции отправки доставок
	produceDeliveries()
}