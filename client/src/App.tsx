import React from 'react';
import './App.css';
import CountryComponent from './CountryComponent';
import { countries } from './CountryList';
import TerritoryComponent from "./territories/TerritoryComponent";
import {data} from "./data";

function App() {



    return (
        <div className="App">
          <header className="App-header">
              <h1>Liste des pays du monde</h1>
              <TerritoryComponent />
          </header>
        </div>
      );
}

export default App;
