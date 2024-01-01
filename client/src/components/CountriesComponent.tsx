import React, {useEffect, useState} from "react";
import {Country, Simulation, Variation} from "../Entity";
import SimulationErrorComponent from "./SimulationErrorComponent";
import {Button, Card, Row} from "react-bootstrap";
import Container from "react-bootstrap/Container";
import CountryDetailComponent from "./CountryDetailComponent";
import OverlayTrigger from "react-bootstrap/OverlayTrigger";
import Tooltip from "react-bootstrap/Tooltip";
import {ResourceIconService} from "../services/ResourceIconService";
import {ClockHistory} from "react-bootstrap-icons";

interface CountriesComponentProps {
    simulation: Simulation | undefined;
}

function CountriesComponent({simulation}: CountriesComponentProps) {
    const [modalDetailCountryShow, setModalDetailCountryShow] = React.useState(false);
    const [selectedCountry, setSelectedCountry] = useState<Country | null>(null);

    useEffect(() => {
        if (selectedCountry && simulation) {
            let simCountry = simulation.countries.get(selectedCountry.agent.id);
            if (simCountry != undefined) {
                setSelectedCountry(simCountry);
            }
        }
    }, [simulation])

    function showDetailCountryModal(country: Country) {
        setSelectedCountry(country);
        setModalDetailCountryShow(true);
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
                                    <Card.Title className="mb-2">{country.agent.name}</Card.Title>
                                    <div className="card text-start text-bg-dark mb-2">
                                        <ul className="list-group list-group-flush">
                                            <li className="list-group-item">
                                                <strong>Habitants:</strong> {country.getCountryPopulation(simulation)}
                                            </li>
                                            <li className="list-group-item">
                                                <strong>Argent:</strong> {country?.money}
                                                <ul className="mt-1">
                                                    <li style={{listStyle: "none"}}>
                                                        <Button size="sm" variant="outline-dark" className="col-auto" onClick={() => setModalDetailCountryShow(true)}>
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
                                                            Value: {quantity}
                                                        </li>
                                                    ))}
                                                    <li style={{listStyle: "none"}}>
                                                        <Button size="sm" variant="outline-dark" className="col-auto" onClick={() => setModalDetailCountryShow(true)}>
                                                            <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des stocks
                                                        </Button>
                                                    </li>
                                                </ul>
                                            </li>
                                        </ul>
                                    </div>
                                    <Button variant="warning" onClick={() => showDetailCountryModal(country)}>
                                        <ClockHistory className="mb-1 me-1"></ClockHistory>Historique des actions
                                    </Button>
                                </Card.Body>
                            </Card>
                            <CountryDetailComponent
                                show={modalDetailCountryShow}
                                onHide ={() => setModalDetailCountryShow(false)}
                                simulation={simulation}
                                propsCountry={selectedCountry}
                            />
                        </div>
                    ))}
                </Row>
            </Container>
        )
    }

}

export default CountriesComponent;