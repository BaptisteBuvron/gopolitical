import React from "react";
import {data} from "../data";
import {getCountryColor} from "../utilities";

function TerritoryComponent() {
    const handleTerritoryClick = (territory: { country: any;  }) => {
        console.log(`Clicked on territory: ${territory.country}`);
    };

    return (
        <div className="Country-tab" style={{ marginTop: "10%" }}>
            {data["territories"].map((territory, index) => (
                <div
                    key={index}
                    className="territory"
                    style={{
                        backgroundColor: getCountryColor(territory.country),
                        left: `${territory.x * 30}px`,
                        top: `${territory.y * 30}px`,
                        width: `${30}px`,
                        height: `${30}px`,
                        position: "absolute", // Ensure absolute positioning
                    }}
                    onClick={() => handleTerritoryClick(territory)}
                >
                    1
                </div>
            ))}
        </div>
    );
}




export default TerritoryComponent;