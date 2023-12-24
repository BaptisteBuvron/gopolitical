import React, { useState } from "react";
import { Modal, Button } from "react-bootstrap";
import { data } from "../../data";
import { getCountryColor } from "../../utils";

interface Territory {
    x: number;
    y: number;
    country: string;
    variations: { name: string; value: number }[];
}

function TerritoryComponent() {
    const [showModal, setShowModal] = useState(false);
    const [selectedTerritory, setSelectedTerritory] = useState<Territory | null>(null);
    const [modalPosition, setModalPosition] = useState({ x: 0, y: 0 });

    const handleTerritoryClick = (territory: Territory, x: number, y: number) => {
        //Si on reclique sur le même territoire = fermeture modal
        //Sinon, ouverture du modal
        if (selectedTerritory && selectedTerritory === territory) {
            setShowModal(false);
            setSelectedTerritory(null);
        } else {
            setSelectedTerritory(territory);
            setModalPosition({ x, y });
            setShowModal(true);
        }
    };

    //Cache le modal quand on appuye sur le bouton Fermer
    const handleCloseModal = () => {
        setShowModal(false);
        setSelectedTerritory(null);
    };

    return (
        <div className="Country-tab">
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
                        position: "absolute",
                    }}
                    //*30 pour que les coordonnées du territoire (x=1, y=3 devient x=30, y=90 sur la map)
                    //+30 pour que le modal ne s'affiche pas sur la case du pays à côté
                    onClick={() => handleTerritoryClick(territory, territory.x * 30 +30, territory.y * 30+30)}
                >
                    1
                </div>
            ))}

            {showModal && (
                <Modal
                    show={showModal}
                    onHide={handleCloseModal}
                    style={{
                        position: "absolute",
                        top: `${modalPosition.y}px`,
                        left: `${modalPosition.x}px`,
                    }}
                    className="custom-modal-container"
                >
                    <Modal.Header closeButton>
                        <Modal.Title>Country : {selectedTerritory && selectedTerritory.country}</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        {selectedTerritory && (
                            <div>
                                {`X : ${selectedTerritory.x}, Y : ${selectedTerritory.y}`}
                            </div>
                        )}
                    </Modal.Body>
                    <Modal.Footer>
                        <Button variant="secondary" onClick={handleCloseModal}>
                            Fermer
                        </Button>
                    </Modal.Footer>
                </Modal>
            )}
        </div>
    );
}

export default TerritoryComponent;
