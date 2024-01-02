import { Country } from "../Entity";

export class CountryService {
    private countries: Map<string, Country>;

    constructor(countries: Map<string, Country>) {
        this.countries = countries
    }

    getCountryByName(name: string) {
        let foundCountry: Country | undefined;

        this.countries.forEach((country) => {
            if (country.agent.name === name) {
                foundCountry = country;
            }
        });

        return foundCountry;

    }
}
