
import React from "react";
import HeaderComponent from "./HeaderComponent";
import {Outlet} from "react-router-dom";

/*interface LayoutComponentProps {
    simulation: Simulation | undefined;
}*/

function LayoutComponent() {
    return (
        <div style={{minHeight: "100vh"}} className="text-bg-dark">
            <HeaderComponent />
            <div>
                <Outlet />
            </div>
        </div>
    );
}

export default LayoutComponent;


