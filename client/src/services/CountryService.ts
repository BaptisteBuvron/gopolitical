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

    getCountryById(countryId: string | undefined): Country | undefined {
        return this.countries.get(countryId || "");
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
