import {Button, Col, Modal, Row} from "react-bootstrap";
import Image from "react-bootstrap/Image";
import React, {useEffect, useState} from "react";
import { Simulation, Territory} from "../Entity";
import HabitantsHistoryChartComponent from "./stockHistoryChartComponent/habitantsHistoryChartComponent";

interface TerritoryHabitantsEvolutionProps {
    onHide: () => void;
    propsTerritory: Territory;
    simulation: Simulation;
    show: boolean;
}

function TerritoryHabitantsEvolutionComponent({onHide, propsTerritory, simulation, show}: TerritoryHabitantsEvolutionProps) {
    //const [showHistoryModal, setShowHistoryModal] = useState(false);
    const [territory, setTerritory] = useState<Territory>(propsTerritory)

    useEffect(() => {
        setTerritory(propsTerritory)
    }, [simulation, propsTerritory])

    const habitantsHistory = territory?.habitantsHistory;
    const country = territory.country;


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
                        <h4 className={"text-warning"}>Evolution des stocks</h4>
                    </div>
                    <div className="col-2">
                        <Image src={country?.flag} alt={country?.name + " flag"} fluid/>
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

export default TerritoryHabitantsEvolutionComponent;