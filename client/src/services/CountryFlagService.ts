
import {countryFlags} from "../countryFlags";

export class CountryFlagService {
    private countryFlags: {country: string, flag: string}[];

    constructor() {
        this.countryFlags = countryFlags;
    }

    getCountryFlagById(countryId: string | undefined): string {
        const country = this.countryFlags.find((c) => c.country === countryId);
        return country ? country.flag : "";
    }
}