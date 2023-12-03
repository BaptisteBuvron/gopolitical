import React from 'react';
import './App.css';

function App() {
    const countries = [
        { name: 'USA', color: '#B22284', x: 0, y: 0 },
        { name: 'Brésil', color: '#FF0000', x: 0, y: 150 },
        { name: 'Emirat Arabes unis', color: '#FF0000', x: 150, y: 0 },
        { name: 'France', color: '#0055A4', x: 150, y: 150 },
        { name: 'Japon', color: '#BC002D', x: 300, y: 0 },
        { name: 'Australie', color: '#0D5EAF', x: 300, y: 150 },
        { name: 'Canada', color: '#FF0000', x: 0, y: 300 },
        { name: 'Chine', color: '#DE2910', x: 150, y: 300 },
        { name: 'Royaume-Uni', color: '#00247D', x: 300, y: 300 },
        { name: 'Inde', color: '#FF9933', x: 450, y: 0 },
        { name: 'Allemagne', color: '#000000', x: 450, y: 150 },
        { name: 'Russie', color: '#0033A0', x: 450, y: 300 },
        { name: 'Afrique du Sud', color: '#007A4D', x: 600, y: 0 },
        { name: 'Mexique', color: '#CE1126', x: 600, y: 150 },
        { name: 'Corée du Sud', color: '#0033A0', x: 600, y: 300 },
    ];


    return (
        <div className="App">
          <header className="App-header">
              <h1>Liste des pays du monde</h1>
              <div className="Country-tab">
                  {countries.map((country, index) => (
                      <div key={index} className="Country"
                           style={{ backgroundColor: country.color,
                              left: `${country.x}px`,
                              top: `${country.y}px`, }}>
                          <p>{country.name}</p>
                      </div>
                  ))}
              </div>
          </header>
        </div>
      );
}

export default App;

