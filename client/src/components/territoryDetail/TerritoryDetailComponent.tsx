import {Territory, Variation} from "../../models/types";
import React from "react";
import {getCountryById, getCountryFlagById, getResourceIconPath} from "../../utils";
import {Button, Modal} from "react-bootstrap";
import {ClockHistory} from "react-bootstrap-icons";
import Image from "react-bootstrap/Image";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import CountryActionsModal from "../countryActionsModal/CountryActionsModal";

interface TerritoryDetailComponentProps {
    handleCloseModal(): void,
    showModal: boolean,
    territory: Territory | null,
}

function TerritoryDetailComponent(props: TerritoryDetailComponentProps) {
    const [modalShow, setModalShow] = React.useState(false);
    let territory = props.territory;
    if(!territory) {
        return (
            <InvalidDataResponseComponent
                handleCloseModal={props.handleCloseModal}
                showModal={props.showModal}
                territory={props.territory}
            />
        );
    }
    const country = getCountryById(territory.country);

    // si on trouve un pays à partir de l'id du territoire, alors on affiche le détail du territoire
    // sinon, on affiche un modal avec un message d'erreur
    if(country !== undefined) {
        return (
            <Modal
                show={props.showModal}
                onHide={props.handleCloseModal}
                animation={false}
                backdrop={true}
                centered
            >
                <Modal.Header className="bg-dark text-light">
                    <div className="d-flex justify-content-between align-items-center">
                        <div className="col-10">
                            <h4 className="card-title">{country.name + " [" + territory.x + "-" + territory.y+"]"} </h4>
                            <Button variant="warning" className="mt-2" onClick={() => setModalShow(true)}>
                                <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des actions
                            </Button>
                        </div>
                        <div className="col-2">
                            <Image src={getCountryFlagById(country.id)} alt={country.name + " flag"} fluid />
                        </div>
                    </div>
                </Modal.Header>
                <Modal.Body className="bg-dark text-light">
                            <div className="card territory-card">
                                <ul className="list-group list-group-flush">
                                    {
                                        territory.variations.map((variation: Variation, index) => (
                                            <li key={index} className="list-group-item">
                                                <OverlayTrigger
                                                    key="bottom"
                                                    placement="bottom"
                                                    overlay={
                                                        <Tooltip>
                                                            {variation.name.charAt(0).toUpperCase() + variation.name.slice(1)}
                                                        </Tooltip>
                                                    }
                                                >
                                                    <img src={getResourceIconPath(variation.name)} className="me-2" alt={variation.name + " icon"} />
                                                </OverlayTrigger>
                                                Value: {variation.value}
                                            </li>
                                        ))
                                    }
                                </ul>
                            </div>
                            <CountryActionsModal
                                show={modalShow}
                                onHide ={() => setModalShow(false)}
                                country={country}
                            />
                </Modal.Body>
                <Modal.Footer className="bg-dark text-light">
                    <Button variant="secondary" onClick={props.handleCloseModal}>
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
        />
    );
}

function InvalidDataResponseComponent(props: TerritoryDetailComponentProps) {
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