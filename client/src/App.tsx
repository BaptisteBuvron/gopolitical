import React from 'react';
import './App.css';
import CountryComponent from './CountryComponent';
import { countries } from './CountryList';

function App() {



    return (
        <div className="App">
          <header className="App-header">
              <h1>Liste des pays du monde</h1>
              <CountryComponent countries={countries} />
          </header>
        </div>
      );
}

export default App;
