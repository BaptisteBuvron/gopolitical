package gopolitical

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	Simulation *Simulation
	Clients    map[*websocket.Conn]struct{}
	mu         *sync.Mutex
}

func NewWebSocket(simulation *Simulation) *WebSocket {
	return &WebSocket{
		Simulation: simulation,
		Clients:    make(map[*websocket.Conn]struct{}),
		mu:         &sync.Mutex{},
	}
}

func (ws *WebSocket) Start() {
	http.HandleFunc("/ws", ws.handleWebSocket)

	port := 8080
	fmt.Printf("Serveur WebSocket écoutant sur le port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"json"}, // Add the desired subprotocol(s)
}

func (ws *WebSocket) handleWebSocket(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		Info("Websocket", "La connection n'a pas pu être établie : %v", err)
		return
	}
	defer conn.Close()

	// Ajouter la nouvelle connexion à la liste des clients
	ws.mu.Lock()
	ws.Clients[conn] = struct{}{}
	ws.mu.Unlock()

	Debug("Websocket", "Connection au client")

	// Attendre des mises à jour depuis la connexion
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			Info("Websocket", "La connexion a été interrompu : %v", err)
			break
		}
	}

	// Retirer la connexion fermée de la liste des clients
	ws.mu.Lock()
	delete(ws.Clients, conn)
	ws.mu.Unlock()
}

func (ws *WebSocket) SendUpdate() {
	// Convertir la simulation en JSON
	ws.mu.Lock()
	defer ws.mu.Unlock()
	simulationJSON, err := json.Marshal(ws.Simulation)
	if err != nil {
		Info("Websocket", "Erreur lors de la sérialisation JSON de la simulation: %v", err)
		return
	}

	// Parcourir toutes les connexions actives et envoyer la mise à jour
	for conn := range ws.Clients {
		err := conn.WriteMessage(websocket.TextMessage, simulationJSON)
		if err != nil {
			Info("Websocket", "Erreur lors de l'envoi de la mise à jour via WebSocket: %v", err)
		}
	}
}
