import {data} from "./data";
import {Country} from "./models/types";
import {countryFlags} from "./countryFlags";
import {resourceIcons} from "./resourceIcons";

export const getCountryColor = (countryId: string): string | undefined => {
    const country = getCountryById(countryId);
    const color = country ? country.color : "D9D9D9";
    return "#" + color;
};

export const getCountryByTerritory = (x: number, y:number): Country | undefined => {
    const territory = data.territories.find((t) => t.x === x && t.y === y);
    if (territory !== undefined) {
        return data.countries.find((c) => c.id === territory.country);
    }
    return undefined;
}

export const getCountryById = (countryId: string): Country | undefined => {
    return data.countries.find((c) => c.id === countryId);
}

export const getCountryFlagById = (countryId: string): string => {
    const country = countryFlags.find((c) => c.country === countryId);
    return country ? country.flag : "";
}

export const getResourceIconPath = (resource: string): string => {
    const resourceIcon = resourceIcons.find((r) => r.resource === resource);
    return resourceIcon ? resourceIcon.iconPath : "";
}