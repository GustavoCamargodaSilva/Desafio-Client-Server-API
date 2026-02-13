package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResult struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar requisicao:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("TIMEOUT: tempo de resposta do servidor excedeu 300ms")
		}
		log.Fatal("Erro ao fazer requisicao:", err)
	}
	defer resp.Body.Close()

	var cotacao CotacaoResult
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		log.Fatal("Erro ao decodificar resposta:", err)
	}

	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", cotacao.Bid)), 0644)
	if err != nil {
		log.Fatal("Erro ao salvar arquivo:", err)
	}

	fmt.Printf("Cotacao salva com sucesso! Dólar: %s\n", cotacao.Bid)
}
