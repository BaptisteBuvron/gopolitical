import React from "react";
import {Simulation} from "../Entity";
import SimulationErrorComponent from "./SimulationErrorComponent";

interface CountryComponentProps {
    simulation: Simulation | undefined;
}

function CountryComponent({simulation}: CountryComponentProps) {
    if(simulation === undefined) {
        return (
            <SimulationErrorComponent />
        )
    }
    return (
        <div className="Country-tab">
            {Array.from(simulation.countries.values()).map((country, index) => (
                <div key={index} className="country"
                     style={{
                         backgroundColor: "#" + country.color,
                         width: "100px",
                     }}>
                    <p>{country.agent.name}</p>
                </div>
            ))}
        </div>
    )
}

export default CountryComponent;