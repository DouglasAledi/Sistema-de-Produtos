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
	r.HandleFunc("/nome", metodo).Methods("GET")
	r.HandleFunc("/nome/{nome}", metodo).Methods("GET")
	r.HandleFunc("/nome", metodo).Methods("POST")
	r.HandleFunc("/nome/{nome}", metodo).Methods("PUT")
	r.HandleFunc("/nome/{nome}", metodo).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}

func metodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nome := vars["nome"]

	switch r.Method {
	case http.MethodGet:
		if nome != "" {
			metodoVerNome(w, r)
		} else {
			metodoVerNomes(w, r)
		}
	case http.MethodPost:
		metodoAdicionar(w, r)
	case http.MethodDelete:
		metodoDelete(w, r)
	case http.MethodPut:
		metodoUpdate(w, r)
	}
}

func metodoVerNomes(w http.ResponseWriter, r *http.Request) {

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

	file, err := os.OpenFile("Nome.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo", http.StatusInternalServerError)
		return
	}

	var names []Name

	defer file.Close()

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	if names == nil {
		names = []Name{}
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

	defer r.Body.Close()
	defer file.Close()

	w.WriteHeader(http.StatusCreated)
}

func metodoUpdate(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	vars := mux.Vars(r)
	nomeAntigo := vars["nome"]

	file, err := os.OpenFile("Nome.json", os.O_RDWR|os.O_CREATE, 0644)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	var names []Name

	err = json.NewDecoder(file).Decode(&names)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	var novo Name

	err = json.NewDecoder(r.Body).Decode(&novo)

	if err != nil {
		http.Error(w, "Erro ao ler corpo", http.StatusBadRequest)
		return
	}

	encontrado := false

	for i, v := range names {
		if v.Nome == nomeAntigo {
			names[i].Nome = novo.Nome
			encontrado = true
			break
		}
	}

	if !encontrado {
		http.Error(w, "Nome não encontrado", http.StatusNotFound)
		return
	}

	file.Truncate(0)
	file.Seek(0, 0)
	json.NewEncoder(file).Encode(names)

	w.WriteHeader(http.StatusOK)
}

func metodoDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	nome := vars["nome"]

	file, err := os.OpenFile("Nome.json", os.O_RDWR|os.O_CREATE, 0644)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var names []Name

	err = json.NewDecoder(file).Decode(&names)

	if err != nil && err != io.EOF {
		http.Error(w, "Erro ao ler JSON", http.StatusInternalServerError)
		return
	}

	encontrado := false

	for i, v := range names {
		if v.Nome == nome {
			names = append(names[:i], names[i+1:]...)
			encontrado = true
			break
		}
	}

	if !encontrado {
		http.Error(w, "Nome não encontrado", http.StatusNotFound)
		return
	}

	file.Truncate(0)
	file.Seek(0, 0)
	json.NewEncoder(file).Encode(names)

	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Deletado com sucesso"))

}
