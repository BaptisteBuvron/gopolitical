import {Button, Col, Modal, Row} from "react-bootstrap";
import Image from "react-bootstrap/Image";
import React, {useEffect, useState} from "react";
import {Simulation, Territory} from "../Entity";
import StockHistoryChart from "./stockHistoryChartComponent/stockHistoryChartComponent";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";

interface TerritoryStockEvolutionProps {
    onHide: () => void;
    propsTerritory: Territory;
    simulation: Simulation;
    show: boolean;
}

function TerritoryStockEvolutionComponent({ onHide, propsTerritory, simulation, show }: TerritoryStockEvolutionProps) {
    //const [showHistoryModal, setShowHistoryModal] = useState(false);
    const [territory, setTerritory] = useState<Territory>(propsTerritory)

    useEffect(() => {
        setTerritory(propsTerritory)
    }, [simulation, propsTerritory])

    const stockHistory = territory.stockHistory;
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
                        <h3 className="card-title">
                            {country?.name} | Territoire
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
                        </h3>
                        <h4 className={"text-warning"}>Evolution des stocks</h4>
                    </div>
                    <div className="col-2">
                        <Image src={country?.flag} alt={country?.name + " flag"} fluid />
                    </div>
                </div>
            </Modal.Header>
            <Modal.Body className="bg-dark text-light">
                <Row className="justify-content-center">
                    <Col className="col-10">
                        {stockHistory && <StockHistoryChart stockHistory={stockHistory}/>}
                    </Col>
                </Row>
            </Modal.Body>
            <Modal.Footer className="bg-dark text-light">
                <Button variant="outline-warning" className="col-auto" onClick={onHide}>Retour</Button>
            </Modal.Footer>
        </Modal>
    );
}

export default TerritoryStockEvolutionComponent;