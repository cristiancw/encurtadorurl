package url

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

const (
	tamanho  = 5
	simbolos = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-+"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Tamanho do identificador: %d\n", tamanho)
	fmt.Printf("Simbolos usados: %s\n", simbolos)
	ConfigurarRepo(NovoRepositorioMemoria())
}

// URL usada na aplicação
type URL struct {
	ID      string    `json: "id"`
	Criacao time.Time `json: "criacao"`
	Destino string    `json: "destino"`
}

// Stats busca o status da url.
func (url *URL) Stats() *Stats {
	clicks := repo.BuscarClick(url.ID)
	return &Stats{URL: url, Clicks: clicks}
}

type Stats struct {
	URL    *URL `json:"url"`
	Clicks int  `json:"clicks"`
}

// Repositorio define as funções que devem manipular os dados armazenados.
type Repositorio interface {
	IDExiste(id string) bool
	BuscaPorID(id string) *URL
	BuscaPorURL(url string) *URL
	Salvar(url URL) error
	RegistrarClick(id string)
	BuscarClick(id string) int
}

var repo Repositorio

// ConfigurarRepo define o tipo de repositório.
func ConfigurarRepo(r Repositorio) {
	repo = r
}

// BuscarOuCriarNovaURL usa o parametro passado para identificar se
// já existe essa url cadastrada e devolve sua versão curta,
// senão cria um novo registro também retornando sua versão curta.
// Caso a url seja inválida vai retornar um erro.
func BuscarOuCriarNovaURL(destino string) (*URL, bool, error) {

	// Busca uma já criada e devolve ela
	if u := repo.BuscaPorURL(destino); u != nil {
		return u, false, nil
	}

	// Não tem deve ser criada
	// Primeiro valida o que foi passado
	if _, err := url.ParseRequestURI(destino); err != nil {
		return nil, false, err // Não sendo válido devolve o erro encontrado
	}

	// Se não existe e o parametro esta okay, vamos criar uma nova...
	url := URL{ID: gerarID(), Criacao: time.Now(), Destino: destino}
	// ...agora que foi criada, salva ela...
	repo.Salvar(url)
	// ... e devolve ela
	return &url, true, nil
}

// Buscar busca uma URL cadastrada baseada no id.
func Buscar(id string) *URL {
	return repo.BuscaPorID(id)
}

// RegistrarClick o id passado para calcular quantas vezes foi acessado.
func RegistrarClick(id string) {
	repo.RegistrarClick(id)
}

// BuscarClick quantas vezes foi acessado o id passado.
func BuscarClick(id string) int {
	return repo.BuscarClick(id)
}

func gerarID() string {
	novoID := func() string {
		id := make([]byte, tamanho, tamanho)
		for i := range id {
			id[i] = simbolos[rand.Intn(len(simbolos))]
		}
		return string(id)
	}

	for {
		if id := novoID(); !repo.IDExiste(id) {
			return id
		}
	}
}
