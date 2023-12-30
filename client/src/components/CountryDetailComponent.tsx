import CountryActionsModal, {CountryModalProps} from "./countryActionsModal/CountryActionsModal";
import {Button, Modal} from "react-bootstrap";
import Image from "react-bootstrap/Image";
import React, {useState} from "react";
import {CountryFlagService} from "../services/CountryFlagService";
import {ClockHistory} from "react-bootstrap-icons";

function CountryDetailComponent({ onHide, country, simulation, show }: CountryModalProps) {
    const [showHistoryModal, setShowHistoryModal] = useState(false);

    // Fonction pour obtenir le flag du country
    const countryFlagService = new CountryFlagService();
    const getCountryFlagById = (countryId: string | undefined): string => {
        return countryFlagService.getCountryFlagById(countryId);
    };

    return (
        <Modal
        show={show}
        onHide={onHide}
        centered
        scrollable={true}
        size="lg"
    >
        <Modal.Header className="bg-dark text-light">
            <div className="d-flex justify-content-between align-items-center col-12">
                <div className="col-10">
                    <h3 className="card-title mb-1">{country?.agent.name}</h3>
                    <h4 className={"text-warning"}>Détails du pays</h4>
                </div>
                <div className="col-2">
                    <Image src={getCountryFlagById(country?.agent.id)} alt={country?.agent.name + " flag"} fluid />
                </div>
            </div>
        </Modal.Header>
        <Modal.Body className="bg-dark text-light">
            <Button variant="warning" className="mb-3" onClick={() => setShowHistoryModal(true)}>
                <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des actions du pays
            </Button>
            <h4>Todo : afficher les infos du pays</h4>


            <CountryActionsModal
                onHide={() => setShowHistoryModal(false)}
                country={country}
                simulation={simulation}
                show={showHistoryModal}
            />
        </Modal.Body>
        <Modal.Footer className="bg-dark text-light">
            <Button variant="outline-warning" onClick={onHide} >Fermer</Button>
        </Modal.Footer>
    </Modal>
);
}

export default CountryDetailComponent;