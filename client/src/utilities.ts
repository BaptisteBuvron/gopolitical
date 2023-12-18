import {data} from "./data";

export const getCountryColor = (countryId: string | null): string | undefined => {
    const country = data.countries.find((c) => c.id === countryId);
    return "#" + country?.color;
};