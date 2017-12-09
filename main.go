package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/gorilla/handlers"
	"os"
	"google.golang.org/appengine"
	"fmt"
	"cloud.google.com/go/datastore"
	"encoding/json"
)

var datastoreClient *datastore.Client

// estrutura de dados do acesso
type Acesso struct {
	IdVisitante string `json:"idVisitante,omitempty"`
	Url string `json:"url,omitempty"`
	DataHora string `json:"dataHora,omitempty"`
}

type Usuario struct {
	IdVisitante string `json:"idVisitante,omitempty"`
	Email string `json:"email,omitempty"`
}

type AcessoUsuario struct {
	IdVisitante string `json:"idVisitante,omitempty"`
	Email string `json:"email,omitempty"`
	Url string `json:"url,omitempty"`
	DataHora string `json:"dataHora,omitempty"`
}


// função principal
func main() {
	registerHandlers()
	appengine.Main()
}

// função para direcionar chamadas rest
func registerHandlers() {

	router := mux.NewRouter()

	router.Headers("Access-Control-Allow-Origin", "*")

	router.Handle("/", http.RedirectHandler("/acesso", http.StatusFound)).Headers()

	router.HandleFunc("/lista", GetAcessos).Methods("GET")
	router.HandleFunc("/registra", RegistraAcesso).Methods("POST")
	router.HandleFunc("/registrausuario", RegistraUsuario).Methods("POST")


	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, router))

}

// Função busca acessos no banco NoSQL e retorna para o front-end
func GetAcessos(w http.ResponseWriter, req *http.Request) {

	ctx := appengine.NewContext(req)

	IDprojeto := "gleaming-orbit-188218"

	// Datastore
	var errdts error
	datastoreClient, errdts = datastore.NewClient(ctx, IDprojeto)
	if errdts != nil {

		http.Error(w, fmt.Sprintf(errdts.Error()), 500)
		return
	}

	var acessos []Acesso
	qa := datastore.NewQuery("RegistroAcesso")
	_, err := datastoreClient.GetAll(ctx, qa, &acessos)
	if err != nil {
		// Handle error
	}

	var usuarios []Usuario
	qu := datastore.NewQuery("RegistroUsuario")
	_, erro := datastoreClient.GetAll(ctx, qu, &usuarios)
	if erro != nil {
		// Handle error
	}

	var acessUsuarios []AcessoUsuario

	for j := 0 ; j < len(acessos); j++ {

		acessoUsuario := AcessoUsuario{}

		acessoUsuario.IdVisitante = acessos[j].IdVisitante
		acessoUsuario.DataHora = acessos[j].DataHora
		acessoUsuario.Url = acessos[j].Url
		acessoUsuario.Email = "Não informado"

		for i := 0; i < len(usuarios); i++ {
			if usuarios[i].IdVisitante == acessos[j].IdVisitante {
				acessoUsuario.Email = usuarios[i].Email
			}
		}

		acessUsuarios = append(acessUsuarios, acessoUsuario)

	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(acessUsuarios)

}

// Função registra o acesso do Visitante no banco NoSQL
func RegistraAcesso(w http.ResponseWriter, req *http.Request){

	ctx := appengine.NewContext(req)

	var acesso Acesso

	_ = json.NewDecoder(req.Body).Decode(&acesso)

	IDprojeto := "gleaming-orbit-188218"


	// Datastore
	var errdts error
	datastoreClient, errdts = datastore.NewClient(ctx, IDprojeto)
	if errdts != nil {

		http.Error(w, fmt.Sprintf(errdts.Error()), 500)
		return
	}

	tipo := "RegistroAcesso"

	instaKey := datastore.IncompleteKey(tipo, nil)


	dataAcesso := &Acesso{
		IdVisitante: acesso.IdVisitante,
		Url: acesso.Url,
		DataHora: acesso.DataHora,
	}

	// Salva a nova entidade no banco NoSQL.
	_, erro := datastoreClient.Put(ctx, instaKey,dataAcesso)
	if erro != nil {
		http.Error(w, fmt.Sprintf("Failed to save task: %v", erro), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(acesso)

}

// Função registra o e-mail e o ID do Visitante no banco NoSQL
func RegistraUsuario(w http.ResponseWriter, req *http.Request){

	ctx := appengine.NewContext(req)

	var usuario Usuario

	_ = json.NewDecoder(req.Body).Decode(&usuario)

	IDprojeto := "gleaming-orbit-188218"


	// Datastore
	var errdts error
	datastoreClient, errdts = datastore.NewClient(ctx, IDprojeto)
	if errdts != nil {

		http.Error(w, fmt.Sprintf(errdts.Error()), 500)
		return
	}

	tipo := "RegistroUsuario"

	instaKey := datastore.IncompleteKey(tipo, nil)


	dataUsuario := &Usuario{
		IdVisitante: usuario.IdVisitante,
		Email: usuario.Email,
	}

	// Salva a nova entidade no banco NoSQL.
	_, erro := datastoreClient.Put(ctx, instaKey,dataUsuario)
	if erro != nil {
		http.Error(w, fmt.Sprintf("Failed to save task: %v", erro), 500)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(usuario)
}

