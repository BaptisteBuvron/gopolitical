import {Button, Col, Modal, Row} from "react-bootstrap";
import Image from "react-bootstrap/Image";
import React, {useEffect, useState} from "react";
import {Country, Simulation} from "../Entity";
import HabitantsHistoryChartComponent from "./stockHistoryChartComponent/habitantsHistoryChartComponent";

interface CountryHabitantsEvolutionProps {
    onHide: () => void;
    propsCountry: Country | null;
    simulation: Simulation;
    show: boolean;
}

function CountryHabitantsEvolutionComponent({ onHide, propsCountry, simulation, show }: CountryHabitantsEvolutionProps) {
    //const [showHistoryModal, setShowHistoryModal] = useState(false);
    const [country, setCountry] = useState<Country | null>(propsCountry)

    useEffect(() => {
        if (propsCountry != null) {
            let simCountry = simulation.countries.get(propsCountry.id);
            if (simCountry !== undefined) {
                setCountry(simCountry);
            }
        }
    }, [simulation, propsCountry])

    const habitantsHistory = country?.getAllTerritoriesHabitantsHistory(simulation);

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
                    <h3 className="card-title mb-1">{country?.name}</h3>
                    <h4 className={"text-warning"}>Evolution des habitants</h4>
                </div>
                <div className="col-2">
                    <Image src={country?.flag} alt={country?.name + " flag"} fluid />
                </div>
            </div>
        </Modal.Header>
        <Modal.Body className="bg-dark text-light">
            <Row className="justify-content-center">
                <Col className="col-10">
                    {habitantsHistory && <HabitantsHistoryChartComponent habitantsHistory={habitantsHistory}/>}
                </Col>
            </Row>
        </Modal.Body>
        <Modal.Footer className="bg-dark text-light">
                <Button variant="outline-warning" className="col-auto" onClick={onHide}>Retour</Button>
        </Modal.Footer>
    </Modal>
);
}
export default CountryHabitantsEvolutionComponent;