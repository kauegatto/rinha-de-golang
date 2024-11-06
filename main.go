package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"rinha_backend/models"
	"time"
)

func main() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("pgx", connStr)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	db.SetMaxOpenConns(200)

	http.Handle("POST /clientes/{id}/transacoes", depositHandler(db))
	http.Handle("GET /clients/{id}/transacoes", getExtrato(db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

type ExtractResponse struct {
	AccountBalance models.Client `json:"saldo"`
	Transactions   []Transaction `json:"ultimas_transacoes"`
}

func getExtrato(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		client, err := getClientById(id, db)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "Client not found")
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			fmt.Fprint(w, err.Error())
			return
		}

		transactions, err := getLastTenTransactionsByClientId(id, db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			fmt.Fprint(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := ExtractResponse{
			AccountBalance: client,
			Transactions:   transactions,
		}

		jsonResponse, err := json.Marshal(&response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			fmt.Fprint(w, err.Error())
			return
		}

		fmt.Fprint(w, string(jsonResponse))
	})
}

func depositHandler(db *sql.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id := r.PathValue("id")

			var payment models.Payment
			err := json.NewDecoder(r.Body).Decode(&payment)
			defer r.Body.Close()

			if err != nil {
				fmt.Printf("err: %s\n", err.Error())
				w.WriteHeader(400)
				return
			}

			client, err := getClientById(id, db)

			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			newBalance, err := payment.Validar(client.Balance)

			if err != nil {
				w.WriteHeader(422)
			}

			// todo update newbalance
			_, err = db.Exec("INSERT INTO transactions (client_id, amount, operation, description) VALUES ($1, $2, $3, $4);", id, payment.Valor, payment.Tipo, payment.Descricao)

			if err != nil {

			}
			_, err = db.Exec("UPDATE clients SET balance = $1 WHERE id = $2;", newBalance, id)

			if err != nil {
			}

			//

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"limite": ` + fmt.Sprintf("%d", client.AccountLimit) + `, "saldo": ` + fmt.Sprintf("%d", newBalance) + `}`))
		},
	)
}

func getClientById(id string, db *sql.DB) (models.Client, error) {
	rows := db.QueryRow("SELECT * FROM clients WHERE id = $1", id)

	var accountLimit int64
	var balance int64
	err := rows.Scan(&id, &accountLimit, &balance)
	if err != nil {
		return models.Client{}, err
	}

	return models.Client{
		Date:         time.Now().Format(time.RFC3339Nano),
		AccountLimit: accountLimit,
		Balance:      balance,
	}, nil
}

type Transaction struct {
	Amount      int       `json:"valor"`
	Operation   string    `json:"tipo"`
	Description string    `json:"descricao"`
	CreatedAt   time.Time `json:"realizada_em"`
}

func getLastTenTransactionsByClientId(id string, db *sql.DB) ([]Transaction, error) {
	rows, err := db.Query("SELECT amount, operation, description, created_at FROM transactions WHERE client_id = $1 ORDER BY created_at DESC LIMIT 10;", id)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction

	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(&transaction.Amount, &transaction.Operation, &transaction.Description, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
