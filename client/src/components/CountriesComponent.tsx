import React, {useEffect, useState} from "react";
import {Country, Simulation} from "../Entity";
import SimulationErrorComponent from "./SimulationErrorComponent";
import {Button, Card, Row} from "react-bootstrap";
import Container from "react-bootstrap/Container";
import CountryStockEvolutionComponent from "./CountryStockEvolutionComponent";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import {ResourceIconService} from "../services/ResourceIconService";
import {ClockHistory} from "react-bootstrap-icons";
import CountryActionsModal from "./countryActionsModal/CountryActionsModal";
import CountryMoneyEvolutionComponent from "./CountryMoneyEvolutionComponent";
import CountryHabitantsEvolutionComponent from "./CountryHabitantsEvolutionComponent";

interface CountriesComponentProps {
    simulation: Simulation | undefined;
}

function CountriesComponent({simulation}: CountriesComponentProps) {
    const [showModalStockEvolutionCountry, setShowModalStockEvolutionCountry] = useState(false);
    const [showModalMoneyEvolutionCountry, setShowModalMoneyEvolutionCountry] = useState(false);
    const [showModalHabitantsEvolutionCountry, setShowModalHabitantsEvolutionCountry] = useState(false);
    const [showCountryActionsModal, setShowCountryActionsModal] = useState(false);
    const [selectedCountry, setSelectedCountry] = useState<Country | null>(null);

    useEffect(() => {
        if (selectedCountry && simulation) {
            let simCountry = simulation.countries.get(selectedCountry.id);
            if (simCountry !== undefined) {
                setSelectedCountry(simCountry);
            }
        }
    }, [simulation, selectedCountry])

    function showCountryActions(country: Country) {
        setSelectedCountry(country);
        setShowCountryActionsModal(true);
    }

    function hideCountryActions() {
        setSelectedCountry(null);
        setShowCountryActionsModal(false);
    }

    function showStockEvolution(country: Country) {
        setSelectedCountry(country);
        setShowModalStockEvolutionCountry(true);
    }

    function showHabitantsEvolution(country: Country) {
        setSelectedCountry(country);
        setShowModalHabitantsEvolutionCountry(true);
    }

    function hideStockEvolution() {
        setShowModalStockEvolutionCountry(false);
        setSelectedCountry(null);
    }

    function showMoneyEvolution(country: Country) {
        setSelectedCountry(country);
        setShowModalMoneyEvolutionCountry(true);
    }

    function hideMoneyEvolution() {
        setShowModalMoneyEvolutionCountry(false);
    }
        
    function hideHabitantsEvolution() {
        setShowModalHabitantsEvolutionCountry(false);
        setSelectedCountry(null);
    }

    if(simulation === undefined) {
        return (
            <SimulationErrorComponent />
        )
    }
    else {
        // Fonction pour obtenir l'icône de ressource par nom de ressource
        const resourceIconService = new ResourceIconService();
        const getResourceIconPath = (resource: string): string => {
            return resourceIconService.getResourceIconPath(resource);
        };
        return (
            <Container>
                <h1 style={{textAlign: "center"}}>Liste des pays</h1>
                <Row className="gx-5">
                    {Array.from(simulation.countries.values()).map((country, index) => (
                        <div key={index} className="col-4">
                            <Card key={index} className="text-center m-3">
                                <Card.Header style={{
                                    backgroundColor: "#" + country.color,
                                }}></Card.Header>
                                <Card.Body>
                                    <Card.Title className="mb-2">{country.name}</Card.Title>
                                    <div className="card text-start text-bg-dark mb-2">
                                        <ul className="list-group list-group-flush">
                                            <li className="list-group-item">
                                                <strong>Habitants:</strong> {country.getCountryPopulation(simulation)}
                                                <Button size="sm" variant="outline-dark" className="col-auto" onClick={() => showHabitantsEvolution(country)}>
                                                    <ClockHistory className="mb-1 me-1"></ClockHistory>Historique du nombre d'habitants
                                                </Button>
                                            </li>
                                            <li className="list-group-item">
                                                <strong>Argent:</strong> {country?.money}
                                                <ul className="mt-1">
                                                    <li style={{listStyle: "none"}}>
                                                        <Button size="sm" variant="outline-dark" className="col-auto" onClick={() => showMoneyEvolution(country)}>
                                                            <ClockHistory className="mb-1 me-1"></ClockHistory>Historique de l'argent
                                                        </Button>
                                                    </li>
                                                </ul>
                                            </li>
                                            <li className="list-group-item">
                                                <strong>Stock global:</strong>
                                                <ul className="mt-1">
                                                    {Array.from(country.getTotalStocks(simulation).entries()).map(([resource, quantity], index) => (
                                                        <li key={index} style={{listStyle: "none"}} className="mb-2">
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
                                                            Value: {Math.floor(quantity)}
                                                        </li>
                                                    ))}
                                                    <li style={{listStyle: "none"}}>
                                                        <Button size="sm" variant="outline-dark" className="col-auto" onClick={() => showStockEvolution(country)}>
                                                            <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des stocks
                                                        </Button>
                                                    </li>
                                                </ul>
                                            </li>
                                        </ul>
                                    </div>
                                    <Button variant="warning" onClick={() => showCountryActions(country)}>
                                        <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des actions
                                    </Button>
                                </Card.Body>
                            </Card>
                            <CountryStockEvolutionComponent
                                show={showModalStockEvolutionCountry}
                                onHide ={() => hideStockEvolution()}
                                simulation={simulation}
                                propsCountry={selectedCountry}
                            />
                            <CountryMoneyEvolutionComponent
                                show={showModalMoneyEvolutionCountry}
                                onHide ={() => hideMoneyEvolution()}
                                propsCountry={selectedCountry}
                                simulation={simulation}
                            />

                            <CountryHabitantsEvolutionComponent
                                show={showModalHabitantsEvolutionCountry}
                                onHide ={() => hideHabitantsEvolution()}
                                simulation={simulation}
                                propsCountry={selectedCountry}
                            />
                            <CountryActionsModal
                                show={showCountryActionsModal}
                                onHide ={() => hideCountryActions()}
                                country={selectedCountry}
                            />
                        </div>
                    ))}
                </Row>
            </Container>
        )
    }

}

export default CountriesComponent;