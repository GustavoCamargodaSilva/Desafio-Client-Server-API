package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type CotacaoResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

type CotacaoResult struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		log.Fatal("Erro ao abrir banco de dados:", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal("Erro ao criar tabela:", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		bid, err := buscarCotacao()
		if err != nil {
			log.Println("Erro ao buscar cotacao:", err)
			http.Error(w, "Erro ao buscar cotacao", http.StatusInternalServerError)
			return
		}

		err = salvarCotacao(db, bid)
		if err != nil {
			log.Println("Erro ao salvar cotacao no banco:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CotacaoResult{Bid: bid})
	})

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buscarCotacao() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("TIMEOUT: tempo de busca da cotacao na API excedeu 200ms")
		}
		return "", err
	}
	defer resp.Body.Close()

	var cotacao CotacaoResponse
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return "", err
	}

	return cotacao.USDBRL.Bid, nil
}

func salvarCotacao(db *sql.DB, bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := db.ExecContext(ctx, "INSERT INTO cotacoes (bid) VALUES (?)", bid)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("TIMEOUT: tempo de persistencia no banco excedeu 10ms")
		}
		return err
	}

	return nil
}
