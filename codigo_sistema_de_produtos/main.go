package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Name struct {
	Nome string `json:"nome"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/nome", metodoVerNomes).Methods("GET")
	r.HandleFunc("/nome/{nome}", metodoVerNome).Methods("GET")
	r.HandleFunc("/nome", metodoAdicionar).Methods("POST")
	r.HandleFunc("/nome/{nome}", metodoUpdate).Methods("PUT")
	r.HandleFunc("/nome/{nome}", metodoDelete).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}

func metodoVerNomes(w http.ResponseWriter, r *http.Request) {

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

func metodoVerNome(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	nome := vars["nome"]

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

	for _, v := range names {
		if v.Nome == nome {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			json.NewEncoder(w).Encode(v)
			return
		}
	}
	http.Error(w, "Nome não encontrado", http.StatusNotFound)
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

func metodoUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	file, err := os.OpenFile("Nome.json", os.O_RDWR|os.O_CREATE, 0644)

	var names []Name

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

	var procurado Name
	var novoNome string

	err = json.NewDecoder(r.Body).Decode(&procurado)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler corpo", http.StatusBadRequest)
		return
	}

	for i, v := range names {
		if v.Nome == procurado.Nome {
			names[i].Nome = novoNome
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	file.Truncate(0)
	file.Seek(0, 0)
	json.NewEncoder(file).Encode(names)
}

func metodoDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	file, err := os.OpenFile("Nome.json", os.O_RDWR|os.O_CREATE, 0644)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	var names []Name
	var procurado Name

	err = json.NewDecoder(file).Decode(&names)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler corpo", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&procurado)

	if err != nil {
		http.Error(w, "Erro ao ler corpo", http.StatusBadRequest)
		return
	}

	for i, v := range names {
		if v.Nome == procurado.Nome {
			names = append(names[:i], names[i+1:]...)
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	file.Truncate(0)
	file.Seek(0, 0)
	json.NewEncoder(file).Encode(names)
}
