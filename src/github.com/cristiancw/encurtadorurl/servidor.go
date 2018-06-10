package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cristiancw/encurtadorurl/url"
)

var (
	porta   int
	urlBase string
	stats   chan string
)

// Headers para definir os parametros do cabeçalho.
type Headers map[string]string

// RedirecionadorStruct outra maneira de utilizar handlers.
type RedirecionadorStruct struct {
	stats2 chan string
}

func main() {
	stats = make(chan string)
	defer close(stats)
	go registrarEstatisticas(stats)

	http.HandleFunc("/api/encurtar", Encurtador)
	http.HandleFunc("/api/stats/", Visualizar)
	http.HandleFunc("/r/", Redirecionador)
	http.Handle("/r2/", &RedirecionadorStruct{stats})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", porta), nil))
}

func init() {
	porta = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", porta)
}

// Encurtador função que vai receber a requisição com a Url atual e criar uma nova reduzida.
func Encurtador(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST,"})
		return
	}

	url, nova, erro := url.BuscarOuCriarNovaURL(extrairURL(r))

	if erro != nil {
		responderCom(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if nova {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase, url.ID)
	responderCom(w, status, Headers{
		"Location": urlCurta,
		"Link":     fmt.Sprintf("<%s/api/stats/%s>; rel=\"stats\"", urlBase, url.ID),
	})
}

// Visualizar função que vai receber a requisição para retornar o status da url encurtada.
func Visualizar(w http.ResponseWriter, r *http.Request) {
	buscarURLeExecutarFuncao(w, r, func(urlO *url.URL) {
		json, err := json.Marshal(urlO.Stats())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		responderComJSON(w, string(json))
	})
}

// Redirecionador função que vai receber a requisição a Url reduzida e redirecionar para a Url original.
func Redirecionador(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirecionador 1")
	buscarURLeExecutarFuncao(w, r, func(urlO *url.URL) {
		http.Redirect(w, r, urlO.Destino, http.StatusMovedPermanently)
		// Poderia ser assim, porém, com o passar do tempo mais metricas vão aparecer
		// como não queremos tornar lento um processo por coletar suas metricas
		// trabalhamos com uma goroutine
		// url.RegistrarClick(id)

		stats <- urlO.ID
	})
}

// ServeHTTP segunda maneira de fazer um handler para o redirecionador.
func (red *RedirecionadorStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Redirecionador 2")
	buscarURLeExecutarFuncao(w, r, func(urlO *url.URL) {
		http.Redirect(w, r, urlO.Destino, http.StatusMovedPermanently)
		red.stats2 <- urlO.ID
	})
}

func responderCom(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}

func responderComJSON(w http.ResponseWriter, json string) {
	responderCom(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, json)
}

func extrairURL(r *http.Request) string {
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	return string(url)
}

func registrarEstatisticas(ids <-chan string) {
	for id := range ids {
		url.RegistrarClick(id)
		fmt.Printf("Click registrado com sucesso para o ID: %v\n", id)
	}
}

func buscarURLeExecutarFuncao(w http.ResponseWriter, r *http.Request, executar func(*url.URL)) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]
	if urlO := url.Buscar(id); urlO != nil {
		executar(urlO)
	} else {
		http.NotFound(w, r)
	}
}
