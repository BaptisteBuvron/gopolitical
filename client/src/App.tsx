import React from 'react';
import './App.css';
import CountryComponent from './CountryComponent';
import { countries } from './CountryList';
import TerritoryComponent from "./territories/TerritoryComponent";
import {data} from "./data";
import TerritoryDetailComponent from "./detailComponent/TerritoryDetailComponent";

function App() {
    return (
        <div id="app" className="col-12">
          {/*<header className="App-header">
              <h1>Liste des pays du monde</h1>
              <CountryComponent countries={countries} />
              <TerritoryComponent />

          </header>*/}
            <TerritoryDetailComponent data={data} x={0} y={0} />
        </div>
      );
}

export default App;
