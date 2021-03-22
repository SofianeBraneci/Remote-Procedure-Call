package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"

	st "../client/structures"
	tr "./travaux"
)

var ADRESSE = "localhost"

var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"}
var map_persones = make(map[int]*personne_serv)

type Request struct {
	personeID    int
	method       string
	responseChan chan string
}

// type d'un paquet de personne stocke sur le serveur, n'implemente pas forcement personne_int (qui n'existe pas ici)
type personne_serv struct {
	// A FAIRE

	id       int
	personne st.Personne
	afaire   []func(st.Personne) st.Personne
	statut   string
}

// cree une nouvelle personne_serv, est appelé depuis le client, par le proxy, au moment ou un producteur distant
// produit une personne_dist
func creer(id int) *personne_serv {
	// A FAIRE
	persone := personne_serv{id: id, personne: pers_vide, afaire: nil, statut: "V"}
	return &persone
}

// Méthodes sur les personne_serv, on peut recopier des méthodes des personne_emp du client
// l'initialisation peut être fait de maniere plus simple que sur le client
// (par exemple en initialisant toujours à la meme personne plutôt qu'en lisant un fichier)
func (p *personne_serv) initialise() {
	// A FAIRE

	p.personne = st.Personne{Nom: "BRANECI", Prenom: "Sofiane", Age: 22, Sexe: "M"}
	p.statut = "R"
	funcs := rand.Intn(4) + 1
	for i := 0; i < funcs; i++ {
		p.afaire = append(p.afaire, tr.UnTravail())
	}

}

func (p *personne_serv) travaille() {
	// A FAIRE

	if len(p.afaire) == 0 {
		p.statut = "C"

	} else {
		fun := p.afaire[0]
		println(fun)
		p.afaire = p.afaire[1:]
		p.personne = fun(p.personne)
	}
}

func (p *personne_serv) vers_string() string {
	// A FAIRE
	age := strconv.Itoa(p.personne.Age)
	return p.personne.Nom + " " + p.personne.Prenom + " " + age + " " + p.personne.Sexe

}

func (p *personne_serv) donne_statut() string {
	// A FAIRE
	return p.statut
}

// Goroutine qui maintient une table d'association entre identifiant et personne_serv
// il est contacté par les goroutine de gestion avec un nom de methode et un identifiant
// et il appelle la méthode correspondante de la personne_serv correspondante
func mainteneur(requests chan Request) {
	// A FAIRE

	for {
		request := <-requests
		println(request.method)
		switch request.method {
		case "creer":
			p := creer(request.personeID)
			map_persones[request.personeID] = p
			request.responseChan <- "PERSONE IS CREATED SUCCESSFULLY\n"
		case "initialise":
			persone, ok := map_persones[request.personeID]
			if ok {
				persone.initialise()
				request.responseChan <- "PERSONE IS INITIALISED SUCCESSFULLY\n"
			}
		case "donne_statut":
			persone, ok := map_persones[request.personeID]
			if ok {
				statu := persone.donne_statut()
				request.responseChan <- statu + "\n"
			}
		case "vers_string":
			persone, ok := map_persones[request.personeID]
			if ok {
				str := persone.vers_string()
				request.responseChan <- str + "\n"

			}
		case "travaille":
			persone, ok := map_persones[request.personeID]
			if ok {
				persone.travaille()
				request.responseChan <- "PERSONE HASE DONE SOME WORK\n"
			}
		}

	}
}

// Goroutine de gestion des connections
// elle attend sur la socketi un message content un nom de methode et un identifiant et appelle le mainteneur avec ces arguments
// elle recupere le resultat du mainteneur et l'envoie sur la socket, puis ferme la socket
func gere_connection(conn net.Conn, requestChan chan Request) {
	// A FAIRE

	responseCh := make(chan string)
	command, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	println("HANDLER RECEIVED THE FOLLOWING COMMAND " + command)
	splited := strings.Split(command, ",")
	id, _ := strconv.Atoi(splited[0])
	cmdLen := len(splited[1])
	request := Request{personeID: id, method: splited[1][:cmdLen-1], responseChan: responseCh}
	requestChan <- request
	response := <-request.responseChan
	fmt.Fprintf(conn, response)
	println("A RESPONSE WAS SENT TO THE PERSONE DIST WITH THE ID OF: " + splited[0])

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Format: client <port>")
		return
	}
	port, _ := strconv.Atoi(os.Args[1]) // doit être le meme port que le client
	addr := ADRESSE + ":" + fmt.Sprint(port)
	// A FAIRE: creer les canaux necessaires, lancer un mainteneur
	requestChan := make(chan Request)
	go mainteneur(requestChan)
	println("MAINTAINER IS LAUNCHED")
	ln, _ := net.Listen("tcp", addr) // ecoute sur l'internet electronique
	fmt.Println("Ecoute sur", addr)
	for {
		conn, _ := ln.Accept() // recoit une connection, cree une socket
		fmt.Println("Accepte une connection.")
		go gere_connection(conn, requestChan) // passe la connection a une routine de gestion des connections
	}
}
