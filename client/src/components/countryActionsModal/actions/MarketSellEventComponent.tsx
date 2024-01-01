import React from "react";
import { MarketSellEvent } from "../../../Entity";

interface MarketSellEventComponentProps {
    event: MarketSellEvent;
}

function MarketSellEventComponent({ event }: MarketSellEventComponentProps) {
    const { resource, amount, gain, to } = event;

    return (
        <>
            <td>Sell</td>
            <td>{resource}</td>
            <td>{amount}</td>
            <td>{gain}</td>
            <td>{to}</td>
        </>
    );
}

export default MarketSellEventComponent;
