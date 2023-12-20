import React from "react";
import {Button, Modal, Row} from "react-bootstrap";
import {Country} from "../../models/types";
import Image from "react-bootstrap/Image";
import {getCountryFlagById} from "../../utils";
import "../../App.css";

interface CountryActionsModalProps {
    onHide: () => void;
    country: Country; // Remplacez YourTerritoryType par le type exact de votre objet Territory
    show: boolean;
}

function CountryActionsModal({ onHide, country, show }: CountryActionsModalProps) {
    return (
        <Modal
            show={show}
            size="lg"
            aria-labelledby="contained-modal-title-vcenter"
            centered
            scrollable={true}
        >
            <Modal.Header className="bg-dark text-light">
                <div className="d-flex justify-content-between align-items-center">
                    <div className="col-10">
                        <h3 className="card-title mb-3">{country.name}</h3>
                        <h4 className={"text-warning"}>Historique des actions</h4>
                    </div>
                    <div className="col-2">
                        <Image src={getCountryFlagById(country.id)} alt={country.name + " flag"} fluid />
                    </div>
                </div>
            </Modal.Header>
            <Modal.Body className="bg-dark text-light">

                    <h4>Centered Modal</h4>
                    <p>
                        Cras mattis consectetur purus sit amet fermentum. Cras justo odio,
                        dapibus ac facilisis in, egestas eget quam. Morbi leo risus, porta ac
                        consectetur ac, vestibulum at eros.
                    </p>
                    <h4>Centered Modal</h4>
                    <p>
                        Cras mattis consectetur purus sit amet fermentum. Cras justo odio,
                        dapibus ac facilisis in, egestas eget quam. Morbi leo risus, porta ac
                        consectetur ac, vestibulum at eros.
                    </p>
                    <h4>Centered Modal</h4>
                    <p>
                        Cras mattis consectetur purus sit amet fermentum. Cras justo odio,
                        dapibus ac facilisis in, egestas eget quam. Morbi leo risus, porta ac
                        consectetur ac, vestibulum at eros.
                    </p>
                    <h4>Centered Modal</h4>
                    <p>
                        Cras mattis consectetur purus sit amet fermentum. Cras justo odio,
                        dapibus ac facilisis in, egestas eget quam. Morbi leo risus, porta ac
                        consectetur ac, vestibulum at eros.
                    </p>
                    <h4>Centered Modal</h4>
                    <p>
                        Cras mattis consectetur purus sit amet fermentum. Cras justo odio,
                        dapibus ac facilisis in, egestas eget quam. Morbi leo risus, porta ac
                        consectetur ac, vestibulum at eros.
                    </p>

            </Modal.Body>
            <Modal.Footer className="bg-dark text-light">
                <Button variant="warning" onClick={onHide} size="lg">Fermer</Button>
            </Modal.Footer>
        </Modal>
    );
}

export default CountryActionsModal;