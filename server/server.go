package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/evandrorm89/desafio-client-server-api/common"
	_ "github.com/mattn/go-sqlite3"
)

type Cambio struct {
	Usdbr USDBRL `json:"USDBRL"`
}

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	Varbid     string `json:"varBid"`
	Pctchange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	Createdate string `json:"create_date"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Request iniciada")
	defer log.Println("Request finalizada")

	if ctx.Err() != nil {
		log.Println("Request cancelada pelo cliente")
		http.Error(w, "Request cancelada pelo cliente", http.StatusRequestTimeout)
		return
	}

	cambio, err := BuscaCambio()
	if err != nil {
		log.Printf("Erro ao buscar o Cambio: %v", err)
		http.Error(w, "Falha ao obter o cambio", http.StatusInternalServerError)
		return
	}

	response := common.CambioResponse{
		Cambio: cambio.Usdbr.Bid,
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Erro ao serializar o json: %v", err)
		http.Error(w, "Falha ao serializar o json", http.StatusInternalServerError)
		return
	}

	err = saveDb(cambio)

	if err != nil {
		log.Printf("Erro ao salvar os dados no banco: %v", err)
		http.Error(w, "Erro ao salvar os dados no banco", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func BuscaCambio() (*Cambio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var c Cambio
	err = json.Unmarshal(body, &c)
	if err != nil {
		return nil, err
	}
	log.Println(c.Usdbr)
	return &c, nil

}

func saveDb(cambio *Cambio) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	db, err := sql.Open("sqlite3", "./cotacao.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into cotacao(value) values(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, cambio.Usdbr.Bid)
	if err != nil {
		return err
	}
	return nil
}
