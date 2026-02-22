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
	http.ListenAndServe(":8000", nil)
}

func metodoPegar(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var names []Name

	file, err := os.OpenFile("Nome.json", os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "Erro ao abrir arquivo", http.StatusInternalServerError)
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
