import React from 'react';
import './App.css';
import CountryComponent from './components/CountryComponent';
import { countries } from './CountryList';
import TerritoriesComponent from "./components/territories/TerritoriesComponent";
import {data} from "./data";
import TerritoryDetailComponent from "./components/territoryDetail/TerritoryDetailComponent";
import {Territory} from "./models/types";
import Container from 'react-bootstrap/Container';

function App() {
    return (
        <Container id="app" fluid>
            <header className="header p-3">
                <h1>Liste des territoires du monde</h1>
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
