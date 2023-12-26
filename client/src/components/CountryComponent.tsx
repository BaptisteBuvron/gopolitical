import React from "react";

interface CountryComponentProps {
    countries: { name: string; color: string; x: number; y: number }[];
}

const CountryComponent: React.FC<CountryComponentProps> = ({ countries }) => {
    return (
        <div className="Country-tab">
            {countries.map((country, index) => (
                <div key={index} className="Country"
                     style={{
                         backgroundColor: country.color,
                         left: `${country.x}px`,
                         top: `${country.y}px`,
                     }}>
                    <p>{country.name}</p>
                </div>
            ))}
        </div>
    )
}

export default CountryComponent;