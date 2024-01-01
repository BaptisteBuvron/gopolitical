import React, { useState } from "react";
import {Row, Col, Image, Table, Pagination, ButtonGroup, Dropdown} from "react-bootstrap";
import { Simulation } from "../../Entity";
import { ResourceIconService } from "../../services/ResourceIconService";
import './MarketComponent.css';
import { CountryService } from "../../services/CountryService";

interface MarketComponentProps {
    simulation: Simulation | undefined;
}

const MarketComponent: React.FC<MarketComponentProps> = ({ simulation }) => {
    const [currentPage, setCurrentPage] = useState(1);
    const [selectedResource, setSelectedResource] = useState<string | null>(null);
    const [selectedBuyer, setSelectedBuyer] = useState<string | null>(null);
    const [selectedSeller, setSelectedSeller] = useState<string | null>(null);

    if (!simulation || !simulation.environment || !simulation.environment.market) {
        return <div>Loading...</div>;
    }

    const marketData = simulation.environment.market;
    const marketPrices = marketData.prices;
    const marketHistory = marketData.history;

    const countryService = new CountryService(simulation.countries);
    const resourceIconService = new ResourceIconService();

    const itemsPerPage = 12;




    const filteredMarketHistory = marketHistory
        .filter(interaction => (
            (!selectedResource || interaction.resourceType === selectedResource) &&
            (!selectedBuyer || interaction.buyer === selectedBuyer) &&
            (!selectedSeller || interaction.seller === selectedSeller)
        ));

    const sortedMarketHistory = filteredMarketHistory
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
    const handleResourceSelect = (resource: string | null) => setSelectedResource(resource);
    const handleBuyerSelect = (buyer: string | null) => setSelectedBuyer(buyer);
    const handleSellerSelect = (seller: string | null) => setSelectedSeller(seller);

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
            <Row className="justify-content-center">
            <Col md={3} className="market-column market-filters">
                <h2 className="text-center mb-3">Filters</h2>
                <Dropdown as={ButtonGroup} className="mb-3">
                    <Dropdown.Toggle id="resource-filter-dropdown" title={selectedResource || 'Resource'}>{selectedResource || 'Resource'}
                    </Dropdown.Toggle>
                    <Dropdown.Menu>
                        <Dropdown.Item onClick={() => handleResourceSelect(null)}>All</Dropdown.Item>
                        {Array.from(marketPrices.keys()).map((resource, index) => (
                            <Dropdown.Item key={`resource-${index}`} onClick={() => handleResourceSelect(resource)}>
                                {resource}
                            </Dropdown.Item>
                        ))}
                    </Dropdown.Menu>
                </Dropdown>
                <Dropdown as={ButtonGroup} className="mb-3">
                    <Dropdown.Toggle id="buyer-filter-dropdown" title={selectedBuyer || 'Buyer'}>{selectedBuyer || 'Buyer'}</Dropdown.Toggle>
                    <Dropdown.Menu>
                        <Dropdown.Item onClick={() => handleBuyerSelect(null)}>All</Dropdown.Item>
                        {Array.from(new Set(marketHistory.map(interaction => interaction.buyer))).map((buyer, index) => (
                            <Dropdown.Item key={`buyer-${index}`} onClick={() => handleBuyerSelect(buyer)}>
                                {buyer}
                            </Dropdown.Item>
                        ))}
                        {/* Add options for buyers based on your data */}
                    </Dropdown.Menu>
                </Dropdown>
                <Dropdown as={ButtonGroup} className="mb-3">
                    <Dropdown.Toggle id="seller-filter-dropdown" title={selectedSeller || 'Seller'}>{selectedSeller || 'Seller'}</Dropdown.Toggle>
                    <Dropdown.Menu>
                        <Dropdown.Item onClick={() => handleSellerSelect(null)}>All</Dropdown.Item>
                        {Array.from(new Set(marketHistory.map(interaction => interaction.seller))).map((seller, index) => (
                            <Dropdown.Item key={`seller-${index}`} onClick={() => handleSellerSelect(seller)}>
                                {seller}
                            </Dropdown.Item>
                        ))}
                    </Dropdown.Menu>
                </Dropdown>
            </Col>
            </Row>
            {marketTransactions}
        </Row>
    );
};

export default MarketComponent;
