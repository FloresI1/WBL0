package main

import (
	"log"
	"net/http"
	"text/template"
)

func getDeliveryAndPaymentByOrderUIDHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем orderUID из параметров запроса
	orderUID := r.URL.Query().Get("order_uid")
	if orderUID == "" {
		http.Error(w, "Invalid orderUID parameter", http.StatusBadRequest)
		log.Printf("Invalid orderUID parameter. Request URL: %s", r.URL)
		return
	}

	log.Printf("Processing request for orderUID: %s", orderUID)

	if cachedData, ok := getDataFromCache(orderUID); ok {
		log.Printf("Retrieved data from cache: %+v", cachedData)
		renderDeliveryAndPaymentPage(w, cachedData)
		return
	}

	log.Printf("Data not found in cache. Fetching from other sources...")

	var delivery Delivery
	var payment Payment
	var items Items
	var info Info

	cacheMutex.Lock()
	cache.Store(orderUID, PageData{
		OrderUID: orderUID,
		Delivery: delivery,
		Payment:  payment,
		Items:    items,
		Info:     info,
	})
	cacheMutex.Unlock()

	log.Printf("Data saved to cache: %+v", PageData{
		OrderUID: orderUID,
		Delivery: delivery,
		Payment:  payment,
		Items:    items,
		Info:     info,
	})

	renderDeliveryAndPaymentPage(w, &PageData{
		OrderUID: orderUID,
		Delivery: delivery,
		Payment:  payment,
		Items:    items,
		Info:     info,
	})
}

func getDataFromCache(orderUID string) (*PageData, bool) {
    cachedData, ok := cache.Load(orderUID)
    if !ok || cachedData == nil {
        log.Printf("The data for orderUID: %s was not found in the cache.", orderUID)
        
        // Если данных нет в кеше, получаем их из базы данных
        dbData, err := fetchDataFromDB(orderUID)
        if err != nil {
            log.Printf("Error retrieving data from the database: %v", err)
            return nil, false
        }

        // Сохраняем полученные данные в кеше
        pageData := &PageData{
			OrderUID: orderUID,
			Delivery: dbData[0].Delivery,
			Payment:  dbData[0].Payment,
			Items:    dbData[0].Items,
			Info:     dbData[0].Info,
		}
		cacheMutex.Lock()
		cache.Store(orderUID, dbData)
		cacheMutex.Unlock()
		
		log.Printf("Data has been retrieved from the database and cached for orderUID: %s.", orderUID)
		return pageData, true
    }

    comboData, ok := cachedData.([]*Combo)
    if !ok || len(comboData) == 0 {
        log.Printf("Invalid data format in the cache for orderUID: %s.", orderUID)
        return nil, false
    }

    pageData := &PageData{
        OrderUID: orderUID,
        Delivery: comboData[0].Delivery,
        Payment:  comboData[0].Payment,
        Items:    comboData[0].Items,
        Info:     comboData[0].Info,
    }

    log.Printf("Data has been retrieved from the cache for orderUID: %s.", orderUID)
    return pageData, true
}

func renderDeliveryAndPaymentPage(w http.ResponseWriter, data *PageData) {
	tmpl, err := template.ParseFiles("delivery.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing HTML template: %v", err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing HTML template: %v", err)
		return
	}
}

func debugCacheHandler(w http.ResponseWriter, r *http.Request) {
	cache.Range(func(key, value interface{}) bool {
		orderUID := key.(string)
		pageData, ok := value.([]*Combo)
		if !ok || pageData == nil {
			log.Printf("Invalid data format for OrderUID: %s\n", orderUID)
			return true
		}

		log.Printf("OrderUID: %s, Data: %+v\n", orderUID, pageData[0])

		return true
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Debug cache information printed to the console"))
}