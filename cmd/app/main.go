package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bonfimjustino7/kafka-go-example/internal/infra/akafka"
	"github.com/bonfimjustino7/kafka-go-example/internal/infra/repository"
	"github.com/bonfimjustino7/kafka-go-example/internal/infra/web"
	"github.com/bonfimjustino7/kafka-go-example/internal/usecase"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	repository := repository.NewProductRepositoryMysql(db)
	createProductUseCase := usecase.NewCreateProductUseCase(repository)
	listProductUseCase := usecase.NewListProductUseCase(repository)

	productsHandlers := web.NewProductHandlers(createProductUseCase, listProductUseCase)

	r := chi.NewRouter()
	r.Post("/products", productsHandlers.CreateProductHandler)
	r.Get("/products", productsHandlers.ListProductHandler)

	go http.ListenAndServe(":8000", r)

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		dto := usecase.CreateProductInputDto{}
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			fmt.Println("Erro ao receber mensagem", err)
			continue
		}
		_, err = createProductUseCase.Execute(dto)

		if err != nil {
			fmt.Println("Erro ao processar mensagem", err)
			continue
		}
	}
}
