import React, {useEffect, useState} from "react";
import TerritoryDetailComponent from "../territoryDetail/TerritoryDetailComponent";
import {Territory} from "../../Entity";
import {Simulation} from "../../Entity";
import Container from "react-bootstrap/Container";
import SimulationErrorComponent from "../SimulationErrorComponent";

interface TerritoriesComponentProps {
    simulation: Simulation | undefined
}

function TerritoriesComponent({ simulation }: TerritoriesComponentProps) {
    const [showModal, setShowModal] = useState(false);
    const [selectedTerritory, setSelectedTerritory] = useState<Territory | null>(null);

    //On detecte des qu'il y a un changement sur la valeur simulation
    useEffect(() => {
        // La MAJ des nouvelles valeurs se fait sur un modal ouvert
        if (selectedTerritory) {
            // On retrouve le territoires ouvert dans la simulation avec ses coordonnées
            const matchingTerritory = simulation?.territories.find(
                (simTerritory) =>
                    simTerritory.x === selectedTerritory.x && simTerritory.y === selectedTerritory.y
            );

            if (matchingTerritory) {
                // On modifie le modal avec les nouvelle valeurs
                setShowModal(false);
                setSelectedTerritory(null);
                setSelectedTerritory(matchingTerritory);
                setShowModal(true);
            }
        }
    }, [simulation?.territories, selectedTerritory]);


    if(simulation === undefined) {
        return (
            <SimulationErrorComponent />
        )
    }


    const handleTerritoryClick = (territory: Territory, index: number) => {
        //Si on reclique sur le même territoire = fermeture modal
        //Sinon, ouverture du modal
        //console.log(territory)
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
        <Container>
            <h1 style={{textAlign: "center"}}>Carte du monde</h1>
            <div className="territories">
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
                        simulation={simulation}
                        country={selectedTerritory.country}
                        consumption={simulation.environment.consumptionByHabitant}
                    />
                )}
            </div>
        </Container>
    );
}

export default TerritoriesComponent;
