import {Country, Simulation, Territory, Variation} from "../../Entity";
import {ResourceIconService} from "../../services/ResourceIconService";
import {CountryFlagService} from "../../services/CountryFlagService";
import React, {useEffect, useState} from "react";
import {Button, Col, Modal, Row} from "react-bootstrap";
import {ClockHistory} from "react-bootstrap-icons";
import Image from "react-bootstrap/Image";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import CountryActionsModal from "../countryActionsModal/CountryActionsModal";
import TerritoryStockEvolutionComponent from "../TerritoryStockEvolutionComponent";

interface TerritoryDetailComponentProps {
    handleCloseModal(): void,
    showModal: boolean,
    territory: Territory,
    simulation: Simulation,
    country: Country
}

function TerritoryDetailComponent(props: TerritoryDetailComponentProps) {
    const [showModalCountryActions, setShowModalCountryActions] = useState(false);
    const [territory, setTerritory] = useState(props.territory);
    const [showStockEvolutionModal, setShowStockEvolutionModal] = useState(false);
    const [country, setCountry] = useState<Country | undefined>(territory.country)

    useEffect(() => {
        let territory = props.simulation.territories.find(
            (simTerritory) => props.territory.x === simTerritory.x && props.territory.y === simTerritory.y
        )
        if(territory) {
            setTerritory(territory);
            setCountry(territory.country);
        }
    }, [props.simulation, props.territory]);

    if(!territory) {
        return (
            <InvalidDataResponseComponent
                handleCloseModal={props.handleCloseModal}
                showModal={props.showModal}
                territory={props.territory}
                country={props.country}
            />
        );
    }


    // Fonction pour obtenir l'icône de ressource par nom de ressource
    const resourceIconService = new ResourceIconService();
    const getResourceIconPath = (resource: string): string => {
        return resourceIconService.getResourceIconPath(resource);
    };

    // si on trouve un pays à partir de l'id du territoire, alors on affiche le détail du territoire
    // sinon, on affiche un modal avec un message d'erreur
    if(country !== undefined) {
        return (
            <Modal
                show={props.showModal}
                onHide={props.handleCloseModal}
                animation={false}
                centered
                scrollable={true}
            >
                <Modal.Header className="bg-dark text-light">
                    <div className="d-flex justify-content-between align-items-center">
                        <div className="col-10">
                            <h4 className="card-title">
                                {country?.agent.name} | Territoire
                                <OverlayTrigger
                                    placement="right"
                                    overlay={
                                        <Tooltip>
                                            Position du territoire
                                        </Tooltip>
                                    }
                                >
                                    <span>{" (" + territory.x + "," + territory.y + ")"}</span>
                                </OverlayTrigger>
                            </h4>
                            <Button variant="warning" className="mt-2" onClick={() => setShowModalCountryActions(true)}>
                                <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des actions du pays
                            </Button>
                        </div>
                        <div className="col-2">
                            <Image src={country.flag} alt={country?.agent.name + " flag"} fluid />
                        </div>
                    </div>
                </Modal.Header>
                <Modal.Body className="bg-dark text-light">
                            <div className="card territory-card">
                                <ul className="list-group list-group-flush">
                                    {
                                        <div className="card territory-card">
                                            <ul className="list-group list-group-flush">
                                                <li className="list-group-item">
                                                    <strong>Habitants:</strong> {territory.habitants}
                                                </li>
                                                <li className="list-group-item">
                                                    <strong>Stocks:</strong>
                                                    <Row className="justify-content-center">
                                                        <Col className="col-10">
                                                            <Row className="justify-content-between">
                                                        {Array.from(territory.stock.entries()).map(([resource, quantity], index) => (
                                                            <Col key={index} className="col-5 mb-2">
                                                                <OverlayTrigger
                                                                    placement="left"
                                                                    overlay={
                                                                        <Tooltip>
                                                                            {resource.charAt(0).toUpperCase() + resource.slice(1)}
                                                                        </Tooltip>
                                                                    }
                                                                >
                                                                    <img src={getResourceIconPath(resource)} className="me-2" alt={resource + " icon"} />
                                                                </OverlayTrigger>
                                                                Value: {quantity}
                                                            </Col>
                                                        ))}
                                                    </Row>
                                                        </Col>
                                                    </Row>
                                                    <Row className="justify-content-center">
                                                        <Col className="col-auto">
                                                            <Button size="sm" variant="outline-dark" className="col-auto" onClick={() => setShowStockEvolutionModal(true)}>
                                                                <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des stocks
                                                            </Button>
                                                        </Col>
                                                    </Row>
                                                </li>
                                                <li className="list-group-item">
                                                    <strong>Variations:</strong>
                                                    <Row className="justify-content-center">
                                                        <Col className="col-10">
                                                            <Row className="justify-content-between">
                                                                {territory.variations.map((variation: Variation, index) => (
                                                                    <Col key={index} className="col-5">
                                                                        <OverlayTrigger
                                                                            placement="left"
                                                                            overlay={
                                                                                <Tooltip>
                                                                                    {variation.resource.charAt(0).toUpperCase() + variation.resource.slice(1)}
                                                                                </Tooltip>
                                                                            }
                                                                        >
                                                                            <img src={getResourceIconPath(variation.resource)} className="me-2" alt={variation.resource + " icon"} />
                                                                        </OverlayTrigger>
                                                                        Value: {variation.amount}
                                                                    </Col>
                                                                ))}
                                                            </Row>
                                                        </Col>
                                                    </Row>
                                                </li>
                                            </ul>
                                        </div>
                                    }
                                </ul>
                            </div>
                            <CountryActionsModal
                                show={showModalCountryActions}
                                onHide ={() => setShowModalCountryActions(false)}
                                simulation={props.simulation}
                                country={country}
                            />
                            <TerritoryStockEvolutionComponent
                                simulation={props.simulation}
                                onHide={() => setShowStockEvolutionModal(false)}
                                show={showStockEvolutionModal}
                                propsTerritory={territory}
                            />
                </Modal.Body>
                <Modal.Footer className="bg-dark text-light">
                    <Button variant="outline-warning" onClick={props.handleCloseModal}>
                        Fermer
                    </Button>
                </Modal.Footer>
            </Modal>

        );
    }
    else return (
        <InvalidDataResponseComponent
            handleCloseModal={props.handleCloseModal}
            showModal={props.showModal}
            territory={props.territory}
            country={props.country}
        />
    );
}

interface InvalidDataResponseComponentProps {
    handleCloseModal(): void,
    showModal: boolean,
}

function InvalidDataResponseComponent(props: InvalidDataResponseComponentProps) {
    return (
        <Modal
            show={props.showModal}
            onHide={props.handleCloseModal}
            animation={false}
            backdrop={true}
            centered
        >
            <Modal.Header className="bg-dark text-light">
                <h4 className="card-title">Erreur</h4>
            </Modal.Header>
            <Modal.Body className="bg-dark text-light">
                <div className="alert alert-danger d-flex align-items-center" role="alert">
                    <svg xmlns="http://www.w3.org/2000/svg" className="bi bi-exclamation-triangle-fill flex-shrink-0 me-2"
                         viewBox="0 0 16 16" role="img" aria-label="Warning:">
                        <path
                            d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>
                    </svg>
                    <div>
                        Aucune information trouvée sur ce territoire !
                    </div>
                </div>
            </Modal.Body>
            <Modal.Footer className="bg-dark text-light">
                <Button variant="secondary" onClick={props.handleCloseModal}>
                    Fermer
                </Button>
            </Modal.Footer>
        </Modal>
    );
}

export default TerritoryDetailComponent;