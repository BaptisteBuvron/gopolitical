import React from "react";
import {data} from "../data";
import {getCountryColor} from "../utilities";
import LoggingButton from "./Details"; // Import the ButtonComponent




function TerritoryComponent() {
    const handleTerritoryClick = (territory: { country: any; }) => {
        console.log(`Clicked on territory: ${territory.country}`);
    };
    return (
        <div className="Country-tab" style={{marginTop: "10%"}}>
            {data["territories"].map((territory, index) => (
                <div key={index} className="territory"
                     style={{
                         backgroundColor: getCountryColor(territory.country),
                         left: `${territory.x * 20}px`,
                         top: `${territory.y * 20}px`,
                         width: `200px`,
                         height: `200px`,
                         justifyContent: "center",
                         border: "3px solid #fff",
                     }}
                     onClick={() => handleTerritoryClick(territory)}
                >
                    <p>{territory.country}</p>
                    <LoggingButton />
                </div>
            ))}
        </div>
    )
}

export default TerritoryComponent;