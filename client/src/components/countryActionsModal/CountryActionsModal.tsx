import React from "react";
import {Button, Modal} from "react-bootstrap";
import {Country} from "../../Entity";
import {CountryFlagService} from "../../services/CountryFlagService";
import Image from "react-bootstrap/Image";
import "../../App.css";

export interface CountryModalProps {
    onHide: () => void;
    country: Country | null;
    show: boolean;
}

function CountryActionsModal({ onHide, country, show }: CountryModalProps) {

    // Fonction pour obtenir le flag du country
    const countryFlagService = new CountryFlagService();
    const getCountryFlagById = (countryId: string | undefined): string => {
        return countryFlagService.getCountryFlagById(countryId);
    };

    return (
        <Modal
            show={show}
            size="lg"
            centered
            scrollable={true}
            animation={false}
        >
            <Modal.Header className="bg-dark text-light">
                <div className="d-flex justify-content-between align-items-center col-12">
                    <div className="col-10">
                        <h3 className="card-title mb-1">{country?.agent.name}</h3>
                        <h4 className={"text-warning"}>Historique des actions</h4>
                    </div>
                    <div className="col-2">
                        <Image src={getCountryFlagById(country?.agent.id)} alt={country?.agent.name + " flag"} fluid />
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
                <Button variant="warning" onClick={onHide} >Retour</Button>
            </Modal.Footer>
        </Modal>
    );
}

export default CountryActionsModal;