import React, { useState } from "react";
import { Button, Modal, Table, Pagination } from "react-bootstrap";
import { Country, TransferResourceEvent } from "../../Entity";
import Image from "react-bootstrap/Image";
import TransferResourceEventComponent from "./actions/TransferResourceEventComponent";

const ACTIONS_PER_PAGE = 11;

export interface CountryModalProps {
    onHide: () => void;
    country: Country | null;
    show: boolean;
}

function CountryActionsModal({ onHide, country, show }: CountryModalProps) {
    const [currentPage, setCurrentPage] = useState(1);

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
                        <Image src={country?.flag} alt={country?.agent.name + " flag"} fluid />
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
                    <div className="d-flex justify-content-center align-items-center">
                        <Pagination>
                            <Pagination.Prev
                                disabled={currentPage === 1}
                                onClick={() => setCurrentPage((prevPage) => prevPage - 1)}
                            />
                            {Array.from({ length: Math.min(totalPageCount, 5) }, (_, index) => {
                                const startPage = Math.max(1, currentPage - 2);
                                return (
                                    <Pagination.Item
                                        key={startPage + index}
                                        active={startPage + index === currentPage}
                                        onClick={() => setCurrentPage(startPage + index)}
                                    >
                                        {startPage + index}
                                    </Pagination.Item>
                                );
                            })}
                            <Pagination.Next
                                disabled={currentPage === totalPageCount}
                                onClick={() => setCurrentPage((prevPage) => prevPage + 1)}
                            />
                        </Pagination>
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
