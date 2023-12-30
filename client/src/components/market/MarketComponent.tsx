import React from "react";
import { Row, Col, Image, Table } from "react-bootstrap";
import { Simulation } from "../../Entity";
import { ResourceIconService } from "../../services/ResourceIconService";
import { CountryFlagService } from "../../services/CountryFlagService";
import './MarketComponent.css';

interface MarketComponentProps {
    simulation: Simulation | undefined;
}

const MarketComponent: React.FC<MarketComponentProps> = ({ simulation }) => {
    if (!simulation || !simulation.environment || !simulation.environment.market) {
        return <div>Loading...</div>;
    }

    const marketData = simulation.environment.market;
    const marketPrices = marketData.prices;
    const marketHistory = marketData.history;

    const countryFlagService = new CountryFlagService();
    const resourceIconService = new ResourceIconService();

    const getCountryFlagById = (countryId: string | undefined): string => {
        return countryFlagService.getCountryFlagById(countryId);
    };

    const sortedMarketHistory = marketHistory
        .slice()
        .sort((a, b) => {
            const dateA = new Date(a.dateTransaction).getTime();
            const dateB = new Date(b.dateTransaction).getTime();
            return dateB - dateA;
        })
        .slice(0, 12);

    const marketPricesElements = (
        <Table striped bordered hover>
            <tbody>
            <tr>
                {Array.from(marketPrices.entries()).map(([resource, price], index) => (
                    <td key={`price-${index}`} className="">
                        <div className="flex-lg-column m-lg-2 text-center">
                            <img
                                src={resourceIconService.getResourceIconPath(resource)}
                                alt={`${resource} icon`}
                                className="img-fluid resource-icon"
                            />
                            <span className="ms-2"> {price}$/Unit</span>
                        </div>
                    </td>
                ))}
            </tr>

            </tbody>
        </Table>
    );


    const marketTransactions = (
        <Col md={9} className="market-column market-transactions">
            <h2 className="text-center mb-3">Market Transactions</h2>
            <Table striped bordered hover className="table-market-transactions">
                <thead>
                <tr className="text-center">
                    <th>Day</th>
                    <th>Resource</th>
                    <th>Amount</th>
                    <th>Price</th>
                    <th>Cost</th>
                    <th>Buyer</th>
                    <th>Seller</th>
                </tr>
                </thead>
                <tbody>
                {sortedMarketHistory.map((interaction, index) => (
                    <tr key={index}>
                        <td className="text-center">{interaction.dateTransaction}</td>
                        <td className="text-center">
                            <img
                                src={resourceIconService.getResourceIconPath(interaction.resourceType)}
                                alt={`${interaction.resourceType} icon`}
                                className="resource-icon-with-margin"
                            />
                            {interaction.resourceType}
                        </td>
                        <td className="text-center">{interaction.amount}</td>
                        <td className="text-center">{interaction.price}$/Unit</td>
                        <td className="text-center">{interaction.price * interaction.amount}$</td>
                        <td>
                            <Image
                                src={countryFlagService.getCountryFlagById(interaction.buyer.agent.id)}
                                alt={`${interaction.buyer.agent.name} flag`}
                                fluid
                                className="flag-icon-market resource-icon-with-margin-flag"
                            />
                            {interaction.buyer.agent.name}
                        </td>
                        <td>
                            <Image
                                src={countryFlagService.getCountryFlagById(interaction.seller.agent.id)}
                                alt={`${interaction.seller.agent.name} flag`}
                                fluid
                                className="flag-icon-market resource-icon-with-margin-flag"
                            />
                            {interaction.seller.agent.name}
                        </td>
                    </tr>
                ))}
                </tbody>
            </Table>
        </Col>
    );

    return (
        <Row className="market-display">
            {marketPricesElements}
            {marketTransactions}
        </Row>

    );
};

export default MarketComponent;
