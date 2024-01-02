import React, {useEffect, useState} from 'react';
import './App.css';
import TerritoriesComponent from "./components/territories/TerritoriesComponent";
import {Simulation} from "./Entity";
import {Navigate, Route, Routes} from "react-router-dom";
import LayoutComponent from "./components/LayoutComponent";
import CountriesComponent from "./components/CountriesComponent";
import MarketComponent from "./components/market/MarketComponent";

function App() {
    const [simulation, setSimulation] = useState<Simulation>();

    useEffect(() => {
        let socket = new WebSocket('ws://localhost:8080/ws');
        // Réessayer la connexion WebSocket lorsque le serveur n'est pas disponible
        const interval = setInterval(() => {
            if (socket.readyState === WebSocket.CLOSED ) {
                socket = new WebSocket('ws://localhost:8080/ws');
            }
        }, 1000);

        socket.onmessage = function(event) {
            const data = JSON.parse(event.data);

            //console.log(data);

            let simulation = new Simulation(data);
            //console.log(simulation)
            setSimulation(simulation);
        }

        socket.onclose = function(event) {
            //setSimulation(undefined);
            clearInterval(interval);
        }
;
    }, []); // Effect sera exécuté une seule fois après le rendu initial
    return (
        <Routes>
            <Route path="/" element={<LayoutComponent simulation={simulation} />}>
                <Route index element={<TerritoriesComponent simulation={simulation} />}/>
                <Route path="/countries" element={<CountriesComponent simulation={simulation}/>} />
                <Route path="/market" element={<MarketComponent simulation={simulation} />} />
            </Route>
            <Route
                path="*"
                element={<Navigate to="/" />}
            />
        </Routes>
      );
}

export default App;
