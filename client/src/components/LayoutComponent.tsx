import CountryComponent from "./CountryComponent";
import {countries} from "../CountryList";
import {useEffect, useState} from "react";
import {Simulation} from "../Entity";
import TerritoryComponent from "./territories/TerritoriesComponent";
import HeaderComponent from "./HeaderComponent";

interface LayoutComponentProps {
    simulation: Simulation | undefined;
}

function LayoutComponent(props: LayoutComponentProps) {
    console.log(props.simulation);
    return (
        <div>
            <HeaderComponent />
            {/*<header className="App-header">
                <h1>Liste des pays du monde</h1>
                <CountryComponent countries={countries} />

            </header>*/}
            <body className="text-bg-dark">
                {props.simulation && <TerritoryComponent simulation={props.simulation} />}
            </body>
        </div>
    );
}

export default LayoutComponent;


