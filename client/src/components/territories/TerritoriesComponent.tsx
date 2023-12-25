import React, { useState } from "react";
import { Modal, Button } from "react-bootstrap";
import { data } from "../../data";
import {getCountryById, getCountryColor, getCountryFlagById} from "../../utils";
import TerritoryDetailComponent from "../territoryDetail/TerritoryDetailComponent";
import {Territory} from "../../models/types";
import {ClockHistory} from "react-bootstrap-icons";
import Image from "react-bootstrap/Image";
import TestComponent from "../TestComponent";

function TerritoryComponent() {
    const [showModal, setShowModal] = useState(false);
    const [selectedTerritory, setSelectedTerritory] = useState<Territory | null>(null);
    const [modalPosition, setModalPosition] = useState({ x: 0, y: 0 });

    const handleTerritoryClick = (territory: Territory, x: number, y: number, index: number) => {
        //Si on reclique sur le même territoire = fermeture modal
        //Sinon, ouverture du modal
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
            setModalPosition({ x, y });
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
            {data["territories"].map((territory, index) => (
                <div
                    id={"territory-" + index.toString()}
                    key={index}
                    className="territory"
                    style={{
                        backgroundColor: getCountryColor(territory.country),
                        left: `${territory.x * 30}px`,
                        top: `${territory.y * 30}px`,
                        width: `${30}px`,
                        height: `${30}px`,
                        position: "absolute",
                        border: "solid 1px black",

                    }}
                    //*30 pour que les coordonnées du territoire (x=1, y=3 devient x=30, y=90 sur la map)
                    //+30 pour que le modal ne s'affiche pas sur la case du pays à côté
                    onClick={() => handleTerritoryClick(territory, territory.x*30, territory.y*30, index)}
                >
                </div>
            ))}

            {showModal && (
                <TestComponent showModal={showModal} handleCloseModal={handleCloseModal} territory={selectedTerritory}></TestComponent>
                /*<Modal
                    show={showModal}
                    onHide={handleCloseModal}
                    style={{
                        /!*position: "absolute",
                        top: `${modalPosition.y}px`,
                        left: `${modalPosition.x}px`,*!/
                    }}
                    animation={false}
                    backdrop={true}
                    centered
                    className="bg-dark text-light"
                >
                    <Modal.Header >

                    </Modal.Header>
                    <Modal.Body className="bg-dark text-light">
                        { selectedTerritory && (
                            <TerritoryDetailComponent {...selectedTerritory} />
                        )}
                    </Modal.Body>
                    <Modal.Footer className="bg-dark text-light">
                        <Button variant="secondary" onClick={handleCloseModal}>
                            Fermer
                        </Button>
                    </Modal.Footer>
                </Modal>*/
            )}
        </div>
    );
}

export default TerritoryComponent;
