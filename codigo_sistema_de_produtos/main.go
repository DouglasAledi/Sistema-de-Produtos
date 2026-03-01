package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type Name struct {
	Nome string `json:"nome"`
}

func main() {

	http.HandleFunc("/nome", metodoPegar)
	http.HandleFunc("/nome", metodoAdicionar)
	http.ListenAndServe(":8000", nil)
}

func metodoPegar(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var names []Name

	file, err := os.OpenFile("Nome.json", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&names)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(names)
}

func metodoAdicionar(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	file, err := os.OpenFile("Nome.json", os.O_RDWR|os.O_CREATE, 0644)

	var names []Name

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler o arquivo", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&names)

	if err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	var novo Name

	err = json.NewDecoder(r.Body).Decode(&novo)

	if err != nil {
		http.Error(w, "Erro ao ler corpo", http.StatusBadRequest)
		return
	}

	names = append(names, novo)

	file.Truncate(0)
	file.Seek(0, 0)

	err = json.NewEncoder(file).Encode(names)

	if err != nil {
		http.Error(w, "Erro ao salvar JSON", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
