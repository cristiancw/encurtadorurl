package url

type repositorioEmMemoria struct {
	urls   map[string]*URL
	clicks map[string]int
}

// NovoRepositorioMemoria função para criar um novo repositório em memória.
func NovoRepositorioMemoria() *repositorioEmMemoria {
	return &repositorioEmMemoria{
		make(map[string]*URL),
		make(map[string]int),
	}
}

func (r *repositorioEmMemoria) IDExiste(id string) bool {
	_, existe := r.urls[id]
	return existe
}

func (r *repositorioEmMemoria) BuscaPorID(id string) *URL {
	return r.urls[id]
}

func (r *repositorioEmMemoria) BuscaPorURL(url string) *URL {
	for _, u := range r.urls {
		if u.Destino == url {
			return u
		}
	}
	return nil
}

func (r *repositorioEmMemoria) Salvar(url URL) error {
	r.urls[url.ID] = &url
	return nil
}

func (r *repositorioEmMemoria) RegistrarClick(id string) {
	r.clicks[id]++
}

func (r *repositorioEmMemoria) BuscarClick(id string) int {
	return r.clicks[id]
}
