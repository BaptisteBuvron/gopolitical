import React from "react";
import { MarketBuyEvent } from "../../../Entity";

interface MarketBuyEventComponentProps {
    event: MarketBuyEvent;
}

function MarketBuyEventComponent({ event }: MarketBuyEventComponentProps) {
    const { resource, amount, cost, from } = event;

    return (
        <>
            <td>Buy</td>
            <td>{resource}</td>
            <td>{amount}</td>
            <td>{cost}</td>
            <td>{from}</td>
        </>
    );
}

export default MarketBuyEventComponent;
