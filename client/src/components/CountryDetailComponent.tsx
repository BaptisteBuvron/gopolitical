import CountryActionsModal, {CountryModalProps} from "./countryActionsModal/CountryActionsModal";
import {Button, Modal, Row} from "react-bootstrap";
import Image from "react-bootstrap/Image";
import React, {useEffect, useState} from "react";
import {CountryFlagService} from "../services/CountryFlagService";
import {ClockHistory} from "react-bootstrap-icons";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import {Country, Simulation, Variation} from "../Entity";

interface CountryDetailProps {
    onHide: () => void;
    propsCountry: Country | null;
    simulation: Simulation;
    show: boolean;
}

function CountryDetailComponent({ onHide, propsCountry, simulation, show }: CountryDetailProps) {
    const [showHistoryModal, setShowHistoryModal] = useState(false);
    const [country, setCountry] = useState<Country | null>(propsCountry)
    const [countryPopulation, setCountryPopulation] = useState(country?.getCountryPopulation(simulation));

    useEffect(() => {
        if (propsCountry != null) {
            let simCountry = simulation.countries.get(propsCountry.agent.id);
            if (simCountry != undefined) {
                setCountry(simCountry);
                setCountryPopulation(country?.getCountryPopulation(simulation))
            }
        }
    }, [simulation])

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
            <div className="card">
                <ul className="list-group list-group-flush">
                    <li className="list-group-item">
                        <strong>Habitants:</strong> {countryPopulation}
                    </li>
                    <li className="list-group-item">
                        <strong>Money:</strong> {country?.money}
                    </li>
                </ul>
            </div>

            <CountryActionsModal
                onHide={() => setShowHistoryModal(false)}
                country={country}
                simulation={simulation}
                show={showHistoryModal}
            />
        </Modal.Body>
        <Modal.Footer className="bg-dark text-light justify-content-center">
            <Row className="justify-content-between col-12">
                <Button variant="warning" className="col-auto" onClick={() => setShowHistoryModal(true)}>
                    <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des actions du pays
                </Button>
                <Button variant="outline-warning" className="col-auto" onClick={onHide} >Fermer</Button>
            </Row>
        </Modal.Footer>
    </Modal>
);
}
export default CountryDetailComponent;