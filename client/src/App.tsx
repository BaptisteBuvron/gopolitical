import React, {useEffect, useState} from 'react';
import './App.css';
import TerritoryComponent from "./components/territories/TerritoriesComponent";
import {Simulation} from "./Entity";
import {Route, Routes} from "react-router-dom";
import LayoutComponent from "./components/LayoutComponent";

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
        <Routes>
            <Route path="/" element={<LayoutComponent simulation={simulation}/>}>
                {/*<Route index element={<Home />} />
                <Route path="blogs" element={<Blogs />} />
                <Route path="contact" element={<Contact />} />
                <Route path="*" element={<NoPage />} />*/}
            </Route>
        </Routes>
    /*<div className="App">
        <header className="App-header">
            <h1>Liste des pays du monde</h1>
            <CountryComponent countries={countries} />
            <h1>Carte du monde</h1>
            {simulation && <TerritoryComponent simulation={simulation} />}
        </header>
    </div>*/
      );
}

export default App;
