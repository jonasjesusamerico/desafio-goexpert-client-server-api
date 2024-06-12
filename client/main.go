package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const URL_SERVER = "http://localhost:8080"

type Response struct {
	Bid string `json:"bid"`
}

type ErrorRequest struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	rate, err := buscaCotacao(ctx)
	if err != nil {
		fmt.Println("Erro ao buscar dados ao servidor: " + err.Error())
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Println("Falha ao criar arquivo txt:", err)
		return
	}
	defer file.Close()

	fmt.Println("Cotação do dólar: " + rate)
	_, err = file.WriteString("Dólar: " + rate)
	if err != nil {
		fmt.Println("Falha ao escrever a cotação no arquivo:", err)
	}
}

func buscaCotacao(ctx context.Context) (string, error) {

	resp, err := doRequestWithContext(ctx, "GET", "/cotacao")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var response ErrorRequest
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return "", err
		}

		return "", errors.New(response.Error + " - " + response.Details)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Bid, nil
}

func doRequestWithContext(ctx context.Context, method string, uri string) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, method, URL_SERVER+uri, nil)

	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return
}
