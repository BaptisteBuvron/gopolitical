import React from 'react';
import './App.css';
import TerritoryComponent from "./territories/TerritoryComponent";

function App() {



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
