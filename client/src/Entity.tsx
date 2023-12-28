
class Resource {
    id: number;
    name: string;
    quantity: number;

    constructor(data: any) {
        this.id = data.id;
        this.name = data.name;
        this.quantity = data.quantity;
    }
}

class Variation {
    resource: string;
    amount: number;

    constructor(data: any) {
        this.resource = data.resource;
        this.amount = data.amount;
    }

}



class Agent {
    id: number;
    name: string;

    constructor(data: any) {
        this.id = data.id;
        this.name = data.name;
    }
}

class Country {
    agent: Agent;
    color: string;
    money: number;

    constructor(data: any) {
        this.agent = new Agent(data.agent);
        this.color = data.color;
        this.money = data.money;
    }
}

class Territory {
    x: number;
    y: number;
    variations: Variation[];
    stock: Map<string, number>;
    country: Country | null;
    habitants: number;

    constructor(data: any) {
        this.x = data.x;
        this.y = data.y;
        this.variations = data.variations.map((variationData: any) => new Variation(variationData));
        this.stock = new Map<string, number>(Object.entries(data.stock));
        this.country = data.country ? new Country(data.country) : null;
        this.habitants = data.habitants;
    }
}

class MarketInteraction {
    dateTransaction: string; // Adjust the type based on your actual data type for time.Time
    resourceType: string;
    amount: number;
    price: number;
    buyer: Country ;
    seller: Country;

    constructor(data: any) {
        this.dateTransaction = data.dateTransaction; // Adjust the field name based on the actual JSON structure
        this.resourceType = data.resourceType;
        this.amount = data.amount;
        this.price = data.price;
        this.buyer = new Country(data.buyer);
        this.seller = new Country(data.seller);
    }
}

class Market {
    history: MarketInteraction[];
    prices: Map<string, number>;

    constructor(data: any) {
        this.history = data.history.map((interactionData: any) => new MarketInteraction(interactionData));
        this.prices = new Map<string, number>(Object.entries(data.prices));
    }
}

class Environment {
    market : Market;

    constructor(data: any) {
        this.market = new Market(data.market);
    }
}



class Simulation {
    secondByDay: number;
    environment: Environment;
    territories: Territory[];
    countries: Map<string, Country>;

    constructor(data: any) {
        this.secondByDay = data.secondByDay;
        this.environment = new Environment(data.environment);
        this.territories = data.territories.map((territoryData: any) => new Territory(territoryData));
        this.countries = new Map<string, Country>(
            Object.entries(data.countries).map(([countryKey, countryData]: [string, any]) => [
                countryKey,
                new Country(countryData),
            ])
        );
    }
}


class CountryFlagService {
    private countryFlags: any[];

    constructor(countryFlags: any[]) {
        this.countryFlags = countryFlags;
    }

    getCountryFlagById(countryId: string): string {
        const country = this.countryFlags.find((c) => c.country === countryId);
        return country ? country.flag : "";
    }
}

class ResourceIconService {
    private resourceIcons: any[];

    constructor(resourceIcons: any[]) {
        this.resourceIcons = resourceIcons;
    }

    getResourceIconPath(resource: string): string {
        const resourceIcon = this.resourceIcons.find((r) => r.resource === resource);
        return resourceIcon ? resourceIcon.iconPath : "";
    }
}



export { Simulation, Territory, Country, Resource, Variation, Agent, MarketInteraction, Market, Environment, CountryFlagService, ResourceIconService};
