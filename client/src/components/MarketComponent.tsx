import React from 'react';
import { Simulation, Market } from '../Entity';

interface MarketComponentProps {
    simulation: Simulation | undefined;
}

const MarketComponent: React.FC<MarketComponentProps> = ({ simulation }) => {
    if (!simulation) {
        return <div>Loading...</div>;
    }

    const marketData: Market = simulation.environment?.market;

    return (
        <div>
            <h1>Market Component</h1>
            {marketData && (
                <div>
                    <h2>Market Prices</h2>
                    <ul>
                        {Object.entries(marketData.prices).map(([resource, price]) => (
                            <li key={resource}>
                                {resource}: {price}
                            </li>
                        ))}
                    </ul>
                    <h2>Market History</h2>
                    <ul>
                        {marketData.history.map((interaction, index) => (
                            <li key={index}>
                                <p>Date: {interaction.dateTransaction}</p>
                                <p>Resource: {interaction.resourceType}</p>
                                <p>Amount: {interaction.amount}</p>
                                <p>Price: {interaction.price}</p>
                                <p>Buyer: {interaction.buyer.agent.name}</p>
                                <p>Seller: {interaction.seller.agent.name}</p>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </div>
    );
};

export default MarketComponent;
