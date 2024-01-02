import React from "react";
import { TransferResourceEvent } from "../../../Entity";
import { ResourceIconService } from "../../../services/ResourceIconService";

interface TransferResourceEventComponentProps {
    event: TransferResourceEvent;
    day : number;
}

function TransferResourceEventComponent({ event , day}: TransferResourceEventComponentProps) {
    const { from, to, resource, amount } = event;

    const resourceIconService = new ResourceIconService();

    const resourceIconPath = resourceIconService.getResourceIconPath(resource);

    return (
        <>
            <td>Day {day}</td>
            <td>Transfert</td>
            <td>From {from}</td>
            <td>To {to}</td>
            <td>
                <img src={resourceIconPath} alt={`${resource} icon`} height="30" width="30" />
                {resource}
            </td>
            <td>Quantity: {amount.toFixed(1)}</td>
        </>
    );
}

export default TransferResourceEventComponent;
