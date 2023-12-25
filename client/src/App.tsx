import React, {useEffect, useState} from 'react';
import './App.css';
import TerritoryComponent from "./territories/TerritoryComponent";
import {data} from "./data";
import {Simulation} from "./Entity";
import {json} from "stream/consumers";

function App() {

    const [simulation, setSimulation] = useState<Simulation>();

    useEffect(() => {
        let socket = new WebSocket('ws://localhost:8080/ws');
        // Réessayer la connexion WebSocket lorsque le serveur n'est pas disponible
        const interval = setInterval(() => {
            if (socket.readyState === WebSocket.CLOSED) {
                socket = new WebSocket('ws://localhost:8080/ws');
            }
        }, 1000);

        socket.onmessage = function(event) {
            const data = JSON.parse(event.data);

            //parse data to Simulation
            //json beautifier print
            console.log(data);

            let simulation = new Simulation(data);
            setSimulation(simulation);
        }

        socket.onclose = function(event) {
            clearInterval(interval);
        }

        // Fermer la connexion WebSocket lors du démontage du composant
        return () => socket.close();
    }, []); // Effect sera exécuté une seule fois après le rendu initial



    return (
        <div className="App">
          <header className="App-header">
              <h1>Carte du monde</h1>
              <TerritoryComponent />
          </header>
        </div>
      );
}

export default App;
