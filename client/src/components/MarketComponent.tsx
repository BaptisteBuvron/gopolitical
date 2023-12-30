import React from "react";
import { Row, Col, Card } from "react-bootstrap";
import { Simulation } from "../Entity";

interface MarketComponentProps {
    simulation: Simulation | undefined;
}

const MarketComponent: React.FC<MarketComponentProps> = ({ simulation }) => {
    if (!simulation || !simulation.environment || !simulation.environment.market) {
        return <div>Loading...</div>; // Placeholder for loading state
    }

    const marketData = simulation.environment.market;
    const marketPrices = marketData.prices;
    const marketHistory = marketData.history;

    // Mise des transactions dans l'ordre décroissant de la date
    const sortedMarketHistory = marketHistory.slice().sort((a, b) => {
        const dateA = new Date(a.dateTransaction).getTime();
        const dateB = new Date(b.dateTransaction).getTime();
        return dateB - dateA;
    });

    // Liste des prix du marché des différentes ressources
    const marketPricesElements = Array.from(marketPrices.entries()).map(([resource, price], index) => (
        <Card key={`price-${index}`} className="mb-3">
            <Card.Body className="bg-dark text-light">
                <Card.Title>Current Market Price</Card.Title>
                <Card.Text>
                    <strong className="text-warning">Resource:</strong> {resource}
                    <br />
                    <strong className="text-warning">Price:</strong> {price}
                </Card.Text>
            </Card.Body>
        </Card>
    ));

    // Liste des transactions
    const sortedMarketHistoryElements = sortedMarketHistory.map((interaction, index) => (
        <Card key={index} className="mb-3">
            <Card.Body className="bg-dark text-light">
                <Card.Title className="text-warning">Day {interaction.dateTransaction}</Card.Title>
                <Card.Text>
                    <strong>Resource:</strong> {interaction.resourceType}
                    <br />
                    <strong>Amount:</strong> {interaction.amount}
                    <br />
                    <strong>Price:</strong> {interaction.price}
                    <br />
                    <strong>Buyer:</strong> {interaction.buyer.agent.name}
                    <br />
                    <strong>Seller:</strong> {interaction.seller.agent.name}
                </Card.Text>
            </Card.Body>
        </Card>
    ));

    return (
        <Row>
            <Col md={6}>
                <h2 className="text-center mb-3">Market Prices</h2>
                {marketPricesElements}
            </Col>
            <Col md={6}>
                <h2 className="text-center mb-3">Market Transactions</h2>
                {sortedMarketHistoryElements}
            </Col>
        </Row>
    );
};

export default MarketComponent;
