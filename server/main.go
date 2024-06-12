package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func databaseConnection() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "./cambio.db")
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {

	db, err := databaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable := `
		CREATE TABLE IF NOT EXISTS cambio (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			cotacao TEXT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", handlerCotacao)
	fmt.Println("Listening on :8080")
	log.Println(http.ListenAndServe(":8080", nil))
}

type USDBRL struct {
	Bid string `json:"bid"`
}

type Cambio struct {
	USDBRL USDBRL `json:"USDBRL"`
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Cria um contexto com um timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	// Faz pesquisa numa api requisitando a cotação
	rate, err := buscaCotacaoExterna(ctx)
	if err != nil {
		http.Error(w, "{\"message\": \"Servidor demorou muito a responder a requisição\"}", http.StatusRequestTimeout)
		return
	}

	db, err := databaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := salvaCambio(ctx, db, rate); err != nil {
		log.Println("Falha ao salvar a cotação do cambio:", err)
	}

	response := map[string]string{"bid": rate}
	json.NewEncoder(w).Encode(response)
}

func buscaCotacaoExterna(ctx context.Context) (string, error) {
	// Faz requisição aplicando o contexto com o timeout
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var cambio Cambio
	if err := json.NewDecoder(resp.Body).Decode(&cambio); err != nil {
		return "", err
	}

	return cambio.USDBRL.Bid, nil
}

func salvaCambio(ctx context.Context, db *sql.DB, rate string) error {
	query := "INSERT INTO cambio (cotacao, timestamp) VALUES (?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		_, err = stmt.Exec(rate, time.Now())
		return err
	}
}
