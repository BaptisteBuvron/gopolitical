import { Country } from "../Entity";

export class CountryService {
    private countries: Map<string, Country>;

    constructor(countries: Map<string, Country>) {
        this.countries = countries
    }

    getId(countryName: string): string | undefined {
        let foundId: string | undefined;

        this.countries.forEach((country, id) => {
            if (country.agent.name === countryName) {
                foundId = id;
            }
        });

        return foundId;
    }
}
