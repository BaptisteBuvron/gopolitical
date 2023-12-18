import React from "react";
import {data} from "../data";
import {getCountryColor} from "../utilities";


function TerritoryComponent() {
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
                     }}>
                    <p>{territory.country}</p>
                </div>
            ))}
        </div>
    )
}

export default TerritoryComponent;