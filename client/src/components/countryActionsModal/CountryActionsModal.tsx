import React, { useState } from "react";
import { Button, Modal, Table } from "react-bootstrap";
import { Country, MarketBuyEvent, MarketSellEvent, TransferResourceEvent } from "../../Entity";
import { CountryFlagService } from "../../services/CountryFlagService";
import Image from "react-bootstrap/Image";
import MarketBuyEventComponent from "./actions/MarketBuyEventComponent";
import TransferResourceEventComponent from "./actions/TransferResourceEventComponent";
import MarketSellEventComponent from "./actions/MarketSellEventComponent";

const ACTIONS_PER_PAGE = 10;

export interface CountryModalProps {
    onHide: () => void;
    country: Country | null;
    show: boolean;
}

function CountryActionsModal({ onHide, country, show }: CountryModalProps) {
    const [currentPage, setCurrentPage] = useState(1);

    // Fonction pour obtenir le flag du country
    const countryFlagService = new CountryFlagService();
    const getCountryFlagById = (countryId: string | undefined): string => {
        return countryFlagService.getCountryFlagById(countryId);
    };

    const sortedActions = country?.history ? [...country.history].reverse() : [];

    const totalPageCount = Math.ceil(sortedActions.length / ACTIONS_PER_PAGE);

    const paginatedActions = sortedActions.slice(
        (currentPage - 1) * ACTIONS_PER_PAGE,
        currentPage * ACTIONS_PER_PAGE
    );

    return (
        <Modal show={show} size="lg" centered scrollable={true} animation={false}>
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
                <h4>Historique des actions</h4>
                <Table striped bordered hover variant="dark">
                    <tbody>
                    {paginatedActions.map((action, index) => (
                        <tr key={index}>
                            {action.eventType && (
                                <>
                                    {action.eventType.constructor === MarketSellEvent && (
                                        <MarketSellEventComponent event={action.eventType as MarketSellEvent} />
                                    )}
                                    {action.eventType.constructor === MarketBuyEvent && (
                                        <MarketBuyEventComponent event={action.eventType as MarketBuyEvent} />
                                    )}
                                    {action.eventType.constructor === TransferResourceEvent && (
                                        <TransferResourceEventComponent event={action.eventType as TransferResourceEvent} />
                                    )}
                                </>
                            )}
                        </tr>
                    ))}
                    </tbody>
                </Table>
                {totalPageCount > 1 && (
                    <div className="d-flex justify-content-center mt-3">
                        <Button
                            variant="light"
                            disabled={currentPage === 1}
                            onClick={() => setCurrentPage((prevPage) => prevPage - 1)}
                        >
                            Précédent
                        </Button>
                        <span className="mx-2">
              Page {currentPage} sur {totalPageCount}
            </span>
                        <Button
                            variant="light"
                            disabled={currentPage === totalPageCount}
                            onClick={() => setCurrentPage((prevPage) => prevPage + 1)}
                        >
                            Suivant
                        </Button>
                    </div>
                )}
            </Modal.Body>
            <Modal.Footer className="bg-dark text-light">
                <Button variant="warning" onClick={onHide}>
                    Retour
                </Button>
            </Modal.Footer>
        </Modal>
    );
}

export default CountryActionsModal;
