import React from 'react';
import './App.css';
import CountryComponent from './components/CountryComponent';
import { countries } from './CountryList';
import TerritoryComponent from "./components/territories/TerritoriesComponent";
import {data} from "./data";
import {Container} from "react-bootstrap";
import TerritoryDetailComponent from "./components/territoryDetail/TerritoryDetailComponent";
import {Territory} from "./models/types";

function App() {
    return (
        <Container id="app" fluid>
            <header className="header p-3">
                <h1>Liste des territoires du monde</h1>
                <TerritoryComponent />
            </header>
          {/*

              <CountryComponent countries={countries} />
              <TerritoriesComponent />

          */}
            <div className="row justify-content-evenly g-4 col-12 pb-5">
                {
                    data["territories"].map((territory: Territory, index) => (
                        <TerritoryDetailComponent key={index} data={data} x={territory.x} y={territory.y} />
                    ))
                }
            </div>
        </Container>
      );
}

export default App;
