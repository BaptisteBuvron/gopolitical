import React, { useState } from "react";
import TerritoryDetailComponent from "../territoryDetail/TerritoryDetailComponent";
import {Territory} from "../../Entity";
import {Simulation} from "../../Entity";

interface TerritoryComponentProps {
    simulation: Simulation;
}

function TerritoryComponent({ simulation }: TerritoryComponentProps) {
    const [showModal, setShowModal] = useState(false);
    const [selectedTerritory, setSelectedTerritory] = useState<Territory | null>(null);

    const handleTerritoryClick = (territory: Territory, index: number) => {
        //Si on reclique sur le même territoire = fermeture modal
        //Sinon, ouverture du modal
        console.log(territory)
        if (selectedTerritory && selectedTerritory === territory) {
            setShowModal(false);
            setSelectedTerritory(null);
            document.getElementById("territory-" + index)?.classList.remove("selected-territory");
        } else {
            let territories: HTMLCollectionOf<Element> = document.getElementsByClassName("territory");

            // @ts-ignore
            for (let territory of territories) {
                territory.classList.remove("selected-territory");
            }
            document.getElementById("territory-" + index)?.classList.add("selected-territory");
            setSelectedTerritory(territory);
            setShowModal(true);
        }
    };

    //Cache le modal quand on appuye sur le bouton Fermer
    const handleCloseModal = () => {
        let territories: HTMLCollectionOf<Element> = document.getElementsByClassName("territory");
        // @ts-ignore
        for (let territory of territories) {
            territory.classList.remove("selected-territory");
        }
        setShowModal(false);
        setSelectedTerritory(null);
    };

    return (
        <div className="Country-tab">
            {simulation["territories"].map((territory, index) => (
                <div
                    id={"territory-" + index.toString()}
                    key={index}
                    className="territory"
                    style={{
                        backgroundColor: "#" + territory.country?.color,
                        left: `${territory.x * 30}px`,
                        top: `${territory.y * 30}px`,
                        width: `${30}px`,
                        height: `${30}px`,
                        position: "absolute",
                        border: "solid 1px black",

                    }}
                    onClick={() => handleTerritoryClick(territory, index)}
                >
                </div>
            ))}

            {showModal && selectedTerritory != null && (
                <TerritoryDetailComponent
                    showModal={showModal}
                    handleCloseModal={handleCloseModal}
                    territory={selectedTerritory}
                />
            )}
        </div>
    );
}

export default TerritoryComponent;
