package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// Cria uma pasta chamada uploads se não existir
	os.MkdirAll("./uploads", os.ModePerm)

	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)

	fmt.Println("🚀 Mini Cloud rodando em: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
		<h1>☁️ Mini Cloud</h1>
		<form enctype="multipart/form-data" action="/upload" method="post">
			<input type="file" name="file">
			<input type="submit" value="Enviar">
		</form>
	`)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Pega o arquivo
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao pegar o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Cria o arquivo na pasta uploads
	dst, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, "Erro ao salvar o arquivo", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copia o conteúdo
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Erro ao copiar o arquivo", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "✔️ Arquivo %s enviado com sucesso!", handler.Filename)
}
