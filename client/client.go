package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	st "./structures" // contient la structure Personne
	tr "./travaux"    // contient les fonctions de travail sur les Personnes
)

var ADRESSE string = "localhost"                           // adresse de base pour la Partie 2
var FICHIER_SOURCE string = "./conseillers-municipaux.txt" // fichier dans lequel piocher des personnes
var TAILLE_SOURCE int = 450000                             // inferieure au nombre de lignes du fichier, pour prendre une ligne au hasard
const TAILLE_G = 5                                         // taille du tampon des gestionnaires
var NB_G int = 2                                           // nombre de gestionnaires
var NB_P int = 2                                           // nombre de producteurs
var NB_O int = 4                                           // nombre d'ouvriers
var NB_PD int = 2                                          // nombre de producteurs distants pour la Partie 2
var PORT int
var PROXY_CHANNEL = make(chan Request)
var pers_vide = st.Personne{Nom: "", Prenom: "", Age: 0, Sexe: "M"} // une personne vide
var COUNTER = 1

// paquet de personne, sur lequel on peut travailler, implemente l'interface personne_int
type personne_emp struct {
	// A FAIRE
	personne st.Personne
	ligne    int
	afaire   []func(st.Personne) st.Personne
	statut   string
}

// paquet de personne distante, pour la Partie 2, implemente l'interface personne_int
type personne_dist struct {
	// A FAIRE
	id int
}

// interface des personnes manipulees par les ouvriers, les
type personne_int interface {
	initialise()          // appelle sur une personne vide de statut V, remplit les champs de la personne et passe son statut à R
	travaille()           // appelle sur une personne de statut R, travaille une fois sur la personne et passe son statut à C s'il n'y a plus de travail a faire
	vers_string() string  // convertit la personne en string
	donne_statut() string // renvoie V, R ou C
}

type Request struct {
	command      string
	responseChan chan string
}

// fabrique une personne à partir d'une ligne du fichier des conseillers municipaux
// à changer si un autre fichier est utilisé
func personne_de_ligne(l string) st.Personne {
	separateur := regexp.MustCompile("\u0009") // oui, les donnees sont separees par des tabulations ... merci la Republique Francaise
	separation := separateur.Split(l, -1)
	naiss, _ := time.Parse("2/1/2006", separation[7])
	a1, _, _ := time.Now().Date()
	a2, _, _ := naiss.Date()
	agec := a1 - a2
	return st.Personne{Nom: separation[4], Prenom: separation[5], Sexe: separation[6], Age: agec}
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES ***

func lectrice(line int, returnChannel chan string) {
	f, err := os.Open(FICHIER_SOURCE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Text()
	scanner.Text()
	counter := 2
	for scanner.Scan() {
		if counter == line {
			returnChannel <- scanner.Text()
		}
		counter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (p *personne_emp) initialise() {
	// A FAIRE

	returnChannel := make(chan string)
	go lectrice(p.ligne, returnChannel)

	personneLine := <-returnChannel
	p.personne = personne_de_ligne(personneLine)
	funcs := rand.Intn(4) + 1
	for i := 0; i < funcs; i++ {
		p.afaire = append(p.afaire, tr.UnTravail())
	}
	p.statut = "R"
}

func (p *personne_emp) travaille() {
	// A FAIRE

	if len(p.afaire) == 0 {
		p.statut = "C"

	} else {
		fun := p.afaire[0]
		p.afaire = p.afaire[1:]
		fun(p.personne)
	}

}

func (p *personne_emp) vers_string() string {
	// A FAIRE
	age := strconv.Itoa(p.personne.Age)
	return p.personne.Nom + " " + p.personne.Prenom + " " + age + " " + p.personne.Sexe

}

func (p *personne_emp) donne_statut() string {
	// A FAIRE
	return p.statut
}

// *** METHODES DE L'INTERFACE personne_int POUR LES PAQUETS DE PERSONNES DISTANTES (PARTIE 2) ***
// ces méthodes doivent appeler le proxy (aucun calcul direct)

func (p personne_dist) initialise() {
	// A
	id := strconv.Itoa(p.id)
	command := id + ",initialise\n"
	res := make(chan string)
	req := Request{command: command, responseChan: res}
	PROXY_CHANNEL <- req
	response := <-req.responseChan
	println("INIT FROM DIST PERSONE RECEIVED THE FOLLOWING RESPONSE: " + response + " " + p.donne_statut())

}

func (p personne_dist) travaille() {
	// A FAIRE
	id := strconv.Itoa(p.id)
	command := id + ",travaille\n"
	res := make(chan string)
	req := Request{command: command, responseChan: res}
	PROXY_CHANNEL <- req
	response := <-req.responseChan
	println("TRAVAILLE FROM DIST PERSONE RECEIVED THE FOLLOWING RESPONSE: " + response)
}

func (p personne_dist) vers_string() string {
	// A FAIRE

	id := strconv.Itoa(p.id)
	command := id + ",vers_string\n"
	res := make(chan string)
	req := Request{command: command, responseChan: res}
	PROXY_CHANNEL <- req
	response := <-req.responseChan
	println("VERS STRING FROM DIST PERSONE RECEIVED THE FOLLOWING RESPONSE: " + response)
	return response

}

func (p personne_dist) donne_statut() string {
	// A FAIRE
	id := strconv.Itoa(p.id)
	command := id + ",donne_statut\n"
	res := make(chan string)
	req := Request{command: command, responseChan: res}
	PROXY_CHANNEL <- req
	response := <-req.responseChan
	println("DONNE STATUT FROM DIST PERSONE RECEIVED THE FOLLOWING RESPONSE: " + id + " " + response)
	return response
}

// *** CODE DES GOROUTINES DU SYSTEME ***

// Partie 2: contacté par les méthodes de personne_dist, le proxy appelle la méthode à travers le réseau et récupère le résultat
// il doit utiliser une connection TCP sur le port donné en ligne de commande
func proxy(requestChan chan Request, port int) {
	// A FAIRE
	serverAddress := ADRESSE + ":" + strconv.Itoa(port)
	println("Proxy is launched")
	for {
		req := <-requestChan
		// Init connection
		conn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			log.Fatal(err)
		}
		command := req.command
		// sending
		println("PROXY HAS RECEIVED A REQUEST WITH THE FOLLOWING COMMAND : " + command)
		fmt.Fprintf(conn, command)
		response, _ := bufio.NewReader(conn).ReadString('\n')
		conn.Close()
		req.responseChan <- response
	}
}

// Partie 1 : contacté par la méthode initialise() de personne_emp, récupère une ligne donnée dans le fichier source
func lecteur() {
	// A FAIRE
}

// Partie 1: récupèrent des personne_int depuis les gestionnaires, font une opération dépendant de donne_statut()
// Si le statut est V, ils initialise le paquet de personne puis le repasse aux gestionnaires
// Si le statut est R, ils travaille une fois sur le paquet puis le repasse aux gestionnaires
// Si le statut est C, ils passent le paquet au collecteur
func ouvrier(producteursChan chan personne_int, consumerChannel chan personne_int, collectorChan chan personne_int) {
	// A FAIRE
	println("WORKER")

	for {
		persone := <-consumerChannel
		statut := persone.donne_statut()
		if strings.Contains(statut, "\n") {
			statut = strings.Replace(statut, "\n", "", 1)
		}
		println("WORKER STATUT IS " + statut)
		if statut == "V" {
			persone.initialise()
			println("WORKER GOT AN EMPTY PERSONE")
			producteursChan <- persone

		} else if statut == "R" {
			println("WORKER GOT A FULL PACKET")
			persone.travaille()
			producteursChan <- persone
		} else {
			println("WORKER SENDING TO THE COLLLECTOR")
			collectorChan <- persone
		}

	}

}

// Partie 1: les producteurs cree des personne_int implementees par des personne_emp initialement vides,
// de statut V mais contenant un numéro de ligne (pour etre initialisee depuis le fichier texte)
// la personne est passée aux gestionnaires
func producteur(productionChannel chan personne_int) {
	// A FAIRE

	for {
		emptyPersone := personne_emp{personne: pers_vide, ligne: rand.Intn(2000) + 2, statut: "V"}
		//println(emptyPersone.ligne)
		productionChannel <- &emptyPersone
		println("PRODUCTEUR SENDING")
		time.Sleep(1 * time.Second)
	}

}

// Partie 2: les producteurs distants cree des personne_int implementees par des personne_dist qui contiennent un identifiant unique
// utilisé pour retrouver l'object sur le serveur
// la creation sur le client d'une personne_dist doit declencher la creation sur le serveur d'une "vraie" personne, initialement vide, de statut V
func producteur_distant(producteurChan chan personne_int, proxyChannel chan Request) {

	println("Producteur distant")
	// A FAIRE
	for {

		p := personne_dist{id: COUNTER}
		id := strconv.Itoa(COUNTER)
		COUNTER++
		command := id + ",creer\n"
		res := make(chan string)
		req := Request{command: command, responseChan: res}
		proxyChannel <- req
		response := <-req.responseChan
		println("PRODUCTEUR FROM DIST PERSONE RECEIVED THE FOLLOWING RESPONSE: " + response)
		producteurChan <- p
		println("PRODUCTEUR HAS SENT A PERSONE TO THE PRODUCERS CHANNEL, INTO THE NEXT ONE")
		time.Sleep(1 * time.Second)
	}

}

// Partie 1: les gestionnaires recoivent des personne_int des producteurs et des ouvriers et maintiennent chacun une file de personne_int
// ils les passent aux ouvriers quand ils sont disponibles
// ATTENTION: la famine des ouvriers doit être évitée: si les producteurs inondent les gestionnaires de paquets, les ouvrier ne pourront
// plus rendre les paquets surlesquels ils travaillent pour en prendre des autres
func gestionnaire(producteurChan chan personne_int, consumerChannel chan personne_int) {
	// A FAIRE
	current_size := 0
	var queue []personne_int
	println("GESTIONNAIRE ON")
	for {

		select {
		case persone := <-producteurChan:

			if current_size >= TAILLE_G {

				for current_size > 0 {
					println("queue is beeing emptied")
					current_size--
					current_persone := queue[0]
					queue = queue[1:]
					consumerChannel <- current_persone
				}

			}
			println("queue is not full")
			queue = append(queue, persone)
			current_size++

		default:
			//println("current size  = ", current_size)
			if current_size > 0 {
				persone := queue[0]
				queue = queue[1:]
				consumerChannel <- persone
				current_size--
				println("Personne is sent to the workers")
			}

		}

	}
}

// Partie 1: le collecteur recoit des personne_int dont le statut est c, il les collecte dans un journal
// quand il recoit un signal de fin du temps, il imprime son journal.
func collecteur(stop chan int, collecteurChan chan personne_int) {

	// A FAIRE
	journal := ""
	for {

		select {
		case <-stop:
			println("---------------------------COLLECTOR ENDING-------------------------")
			println(journal)
			stop <- 0
			return

		case persone := <-collecteurChan:
			println("COLLECTOR RECEIVED")
			journal = journal + persone.vers_string() + "\n"

		}
	}

}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // graine pour l'aleatoire
	if len(os.Args) < 3 {
		fmt.Println("Format: client <port> <millisecondes d'attente>")
		return
	}
	PORT, _ := strconv.Atoi(os.Args[1])   // utile pour la partie 2
	millis, _ := strconv.Atoi(os.Args[2]) // duree du timeout
	fintemps := make(chan int)
	// A FAIRE
	// creer les canaux
	// lancer les goroutines (parties 1 et 2): 1 lecteur, 1 collecteur, des producteurs, des gestionnaires, des ouvriers
	producersChan := make(chan personne_int)
	consumersChan := make(chan personne_int)
	collecteurChan := make(chan personne_int)
	//stop := make(chan int)
	go gestionnaire(producersChan, consumersChan)
	go proxy(PROXY_CHANNEL, PORT)
	go producteur_distant(producersChan, PROXY_CHANNEL)
	//go producteur(producersChan)
	go ouvrier(producersChan, consumersChan, collecteurChan)
	//go ouvrier(producersChan, consumersChan, collecteurChan)
	// go ouvrier(producersChan, consumersChan, collecteurChan)
	go collecteur(fintemps, collecteurChan)

	// lancer les goroutines (partie 2): des producteurs distants, un proxy

	time.Sleep(time.Duration(millis) * time.Second)
	fintemps <- 0
	<-fintemps
}
