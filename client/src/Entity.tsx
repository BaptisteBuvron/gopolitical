// Define EventType as an interface
interface EventType {}

// Implement the CountryEvent class
class CountryEvent {
    day: number;
    eventType: EventType | undefined;

    constructor(data: any, country?: Country) {
        this.day = data.day;
        //Check the type of the event TraitementEvent or TransferResourceEvent
        if (data.event === "transferResource") {
            this.eventType = new TransferResourceEvent(data, country);
        }
    }
}

// Implement the TransferResourceEvent class
class TransferResourceEvent implements EventType{
    from: Territory;
    to: Territory;
    resource: string;
    amount: number;

    constructor(data: any, country?: Country) {
        this.from = new Territory(data.from, country);
        this.to = new Territory(data.to, country);
        this.resource = data.resource;
        this.amount = data.amount;
    }
}

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
    history: CountryEvent[];
    territories: Territory[];
    moneyHistory: Map<string, number>;

    constructor(data: any) {
        this.agent = new Agent(data.agent);
        this.color = data.color;
        this.money = data.money;
        this.history = data.history.map((eventData: any) => new CountryEvent(eventData, this));
        this.territories = data.territories.map((territoryData: any) => new Territory(territoryData, this));
        this.moneyHistory = new Map<string, number>(Object.entries(data.moneyHistory));
    }
}

class Territory {
    x: number;
    y: number;
    variations: Variation[];
    stock: Map<string, number>;
    stockHistory: Map<number, Map<string, number>>;
    habitants: number;
    habitantsHistory: Map<string, number>;
    country: Country | undefined;

    constructor(data: any, country?: Country) {
        this.x = data.x;
        this.y = data.y;
        this.variations = data.variations.map((variationData: any) => new Variation(variationData));
        this.stock = new Map<string, number>(Object.entries(data.stock));
        this.habitants = data.habitants;
        this.country = country;
        this.stockHistory = new Map<number, Map<string, number>>();
        for (const key in data.stockHistory) {
            const innerMap = new Map<string, number>(Object.entries(data.stockHistory[key]));
            this.stockHistory.set(Number(key), innerMap);
        }
        this.habitantsHistory = new Map<string, number>(Object.entries(data.habitantsHistory));
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
    currentDay: number;

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
        for (let territory of this.territories) {
            for (let country of this.countries.values()) {
                for (let countryTerritory of country.territories) {
                    if (territory.x === countryTerritory.x && territory.y === countryTerritory.y) {
                        territory.country = country;
                    }
                }
            }
        }
        this.currentDay = data.currentDay;
    }
}

export { Simulation, Territory, Country, Resource, Variation, Agent, MarketInteraction, Market, Environment };
