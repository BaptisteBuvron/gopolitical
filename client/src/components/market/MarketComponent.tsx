import React, { useState } from "react";
import { Row, Col, Image, Table, Pagination } from "react-bootstrap";
import { Simulation } from "../../Entity";
import { ResourceIconService } from "../../services/ResourceIconService";
import './MarketComponent.css';
import { CountryService } from "../../services/CountryService";

interface MarketComponentProps {
    simulation: Simulation | undefined;
}

const MarketComponent: React.FC<MarketComponentProps> = ({ simulation }) => {
    const [currentPage, setCurrentPage] = useState(1);

    if (!simulation || !simulation.environment || !simulation.environment.market) {
        return <div>Loading...</div>;
    }

    const marketData = simulation.environment.market;
    const marketPrices = marketData.prices;
    const marketHistory = marketData.history;

    const countryService = new CountryService(simulation.countries);
    const resourceIconService = new ResourceIconService();

    const itemsPerPage = 12;

    const sortedMarketHistory = marketHistory
        .slice()
        .sort((a, b) => {
            const dateA = new Date(a.dateTransaction).getTime();
            const dateB = new Date(b.dateTransaction).getTime();
            return dateB - dateA;
        });

    const pagesCount = Math.ceil(sortedMarketHistory.length / itemsPerPage);
    const visiblePages = Array.from({ length: Math.min(3, pagesCount) }, (_, index) => index + 1);

    const nextPage = () => {
        if (currentPage < pagesCount) {
            setCurrentPage(currentPage + 1);
        }
    };

    const prevPage = () => {
        if (currentPage > 1) {
            setCurrentPage(currentPage - 1);
        }
    };

    const paginate = (pageNumber: number) => setCurrentPage(pageNumber);

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
                {sortedMarketHistory
                    .slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
                    .map((interaction, index) => (
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
                                    src={countryService.getCountryByName(interaction.buyer)?.flag}
                                    alt={`${interaction.buyer} flag`}
                                    fluid
                                    className="flag-icon-market resource-icon-with-margin-flag"
                                />
                                {interaction.buyer}
                            </td>
                            <td>
                                <Image
                                    src={countryService.getCountryByName(interaction.seller)?.flag}
                                    alt={`${interaction.seller} flag`}
                                    fluid
                                    className="flag-icon-market resource-icon-with-margin-flag"
                                />
                                {interaction.seller}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </Table>

            <div className="d-flex justify-content-center">
                <Pagination>
                    <Pagination.Prev onClick={prevPage} />
                    {visiblePages.map((pageNumber) => (
                        <Pagination.Item
                            key={pageNumber}
                            active={pageNumber === currentPage}
                            onClick={() => paginate(pageNumber)}
                        >
                            {pageNumber}
                        </Pagination.Item>
                    ))}
                    {currentPage < pagesCount - 2 && (
                        <Pagination.Ellipsis disabled />
                    )}
                    {currentPage < pagesCount - 1 && (
                        <Pagination.Item onClick={() => paginate(currentPage + 1)}>
                            {currentPage + 1}
                        </Pagination.Item>
                    )}
                    {currentPage < pagesCount && (
                        <Pagination.Item onClick={() => paginate(currentPage + 2)}>
                            {currentPage + 2}
                        </Pagination.Item>
                    )}
                    <Pagination.Next onClick={nextPage} />
                </Pagination>
            </div>
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
