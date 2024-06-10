package main

import (
	"database/sql"
	"encoding/json"

	//"fmt"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"

	"net/http"

	memorycache "github.com/maxchagin/go-memorycache-example"
)

type Order struct {
	Order_uid    string   `json:"order_uid"`
	Track_number string   `json:"track_number"`
	Delivery     Delivery `json:"delivery"`
}

type Delivery struct {
	Name    string `json:"name"`
	City    string `json:"city"`
	Address string `json:"address"`
}

var id string
var order Order

// var orderInfo string
var orderID string
var cacheValue []byte

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {

	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Order ID is missing", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=own_user password=0000 dbname=products sslmode=disable")
	if err != nil {
		http.Error(w, "Error opening database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM goods WHERE order_uid = $1", orderID)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var column1Value, column2Value, column3Value string
		err := rows.Scan(&column1Value, &column2Value, &column3Value)
		if err != nil {
			http.Error(w, "Error scanning database rows", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Column 1: %s, Column 2: %s, Column 3: %s\n", column1Value, column2Value, column3Value)
	}
}

func getCache() {

	cache := memorycache.New(5*time.Minute, 10*time.Minute)
	cache.Set(order.Order_uid, order, 5*time.Minute)
	//cacheValue, _ := cache.Get(order.Order_uid)
}

func main() {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("nats", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("context jsm", err)
	}

	_, err = js.Subscribe("foo", func(msg *nats.Msg) {

		err := json.Unmarshal([]byte(msg.Data), &order)
		if err != nil {

			log.Fatal("json", err)
		}

		db, err := sql.Open("postgres", "host=localhost port=5432 user=own_user password=0000 dbname=products sslmode=disable")
		if err != nil {
			log.Fatal(" open db", err)
		}
		defer db.Close()

		query := "INSERT INTO goods (order_uid, track_number, name, city, address) VALUES ($1, $2, $3, $4, $5 )"
		_, err = db.Exec(query, order.Order_uid, order.Track_number, order.Delivery.Name, order.Delivery.City, order.Delivery.Address)
		if err != nil {

			log.Fatal("db", err)

		}

		// Подтверждение получения сообщения
		msg.Ack()
	}, nats.Durable("main-consumer2"), nats.ManualAck())
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/order", getOrderHandler)

	log.Fatal(http.ListenAndServe(":9999", nil))

	// Ожидание сигнала для завершения программы
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Программа завершена.")
}
