
import React from "react";
import HeaderComponent from "./HeaderComponent";
import {Outlet} from "react-router-dom";
import {Simulation} from "../Entity";

interface LayoutComponentProps {
    simulation: Simulation | undefined;
}

function LayoutComponent({simulation}: LayoutComponentProps) {
    return (
        <div style={{minHeight: "100vh"}} className="text-bg-dark">
            <HeaderComponent simulation={simulation}/>
            <div>
                <Outlet />
            </div>
        </div>
    );
}

export default LayoutComponent;


