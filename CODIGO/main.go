package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_NAME"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func cadastroHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "./static/html/cadastro.html")
	} else if r.Method == "POST" {
		db, err := connectDB()
		if err != nil {
			http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		email := r.FormValue("email")
		senha := r.FormValue("senha")

		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM clientes WHERE email=?)", email).Scan(&exists)
		if err != nil {
			http.Error(w, "Erro ao verificar email", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Este e-mail já está cadastrado!", http.StatusBadRequest)
			return
		}

		hashSenha, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erro ao criptografar a senha", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO clientes (email, senha) VALUES (?, ?)", email, hashSenha)
		if err != nil {
			http.Error(w, "Erro ao inserir os dados no banco", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "./static/html/login.html")
	} else if r.Method == "POST" {
		db, err := connectDB()
		if err != nil {
			http.Error(w, "Erro ao conectar ao banco de dados", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		email := r.FormValue("email")
		senha := r.FormValue("senha")

		var hashSenha string
		err = db.QueryRow("SELECT senha FROM clientes WHERE email=?", email).Scan(&hashSenha)
		if err == sql.ErrNoRows {
			http.Error(w, "E-mail não encontrado", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Erro ao buscar a senha no banco de dados", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashSenha), []byte(senha))
		if err != nil {
			http.Error(w, "Senha incorreta", http.StatusUnauthorized)
			return
		}

		fmt.Fprintf(w, "Login realizado com sucesso!")
	}
}

func main() {
	loadEnv()
	http.HandleFunc("/cadastro", cadastroHandler)
	http.HandleFunc("/login", loginHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	fmt.Println("Servidor rodando na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
