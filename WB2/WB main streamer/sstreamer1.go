package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	// Открытие файла JSON
	jsonFile, err := os.Open("order.json")
	if err != nil {
		log.Fatalf("Ошибка при открытии файла: %v", err)
	}
	defer jsonFile.Close()

	// Чтение содержимого файла JSON
	jsonData, err := io.ReadAll(jsonFile)
	fmt.Println(&jsonData)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	// Подключение к серверу NATS
	nc, err := nats.Connect("htpp://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Создание контекста JetStream
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		log.Fatal(err)
	}

	// Добавление стрима
	//_, err = js.AddStream(&nats.StreamConfig{
	//	Name:     "foo",
	//	Subjects: []string{"foo"},
	//})
	//if err != nil {
	//	log.Fatalf("Ошибка при добавлении стрима: %v", err)
	//}

	// Публикация содержимого файла JSON
	_, err = js.Publish("foo", []byte(jsonData))
	if err != nil {
		log.Fatalf("Ошибка при публикации: %v", err)
	}

}
