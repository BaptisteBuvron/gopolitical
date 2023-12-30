import React, {useState} from "react";
import {Country, Simulation} from "../Entity";
import SimulationErrorComponent from "./SimulationErrorComponent";
import {Button, Card, Row} from "react-bootstrap";
import Container from "react-bootstrap/Container";
import CountryDetailComponent from "./CountryDetailComponent";

interface CountryComponentProps {
    simulation: Simulation | undefined;
}

function CountryComponent({simulation}: CountryComponentProps) {
    const [modalDetailCountryShow, setModalDetailCountryShow] = React.useState(false);
    const [selectedCountry, setSelectedCountry] = useState<Country | null>(null);

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
        return (
            <Container>
                <h1 style={{textAlign: "center"}}>Liste des pays</h1>
                <Row className="gx-5">
                    {Array.from(simulation.countries.values()).map((country, index) => (
                        <div className="col-3">
                            <Card key={index} className="text-center m-3">
                                <Card.Header style={{
                                    backgroundColor: "#" + country.color,
                                }}></Card.Header>
                                <Card.Body>
                                    <Card.Title className="mb-3">{country.agent.name}</Card.Title>
                                    <Button variant="warning" onClick={() => showDetailCountryModal(country)}>Consulter</Button>
                                </Card.Body>
                            </Card>
                            <CountryDetailComponent
                                show={modalDetailCountryShow}
                                onHide ={() => setModalDetailCountryShow(false)}
                                simulation={simulation}
                                country={selectedCountry}
                            />
                        </div>
                    ))}
                </Row>
            </Container>
        )
    }

}

export default CountryComponent;