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
        if (data.event === "sellEvent") {
            this.eventType = new MarketSellEvent(data);
        }
        if (data.event === "buyEvent") {
            this.eventType = new MarketBuyEvent(data);
        }
    }
}

class MarketSellEvent implements EventType{
    resource: string;
    amount: number;
    gain: number;
    to : string;

    constructor(data: any) {
        this.resource = data.resourceType;
        this.amount = Math.ceil(data.amountExecuted);
        this.gain = data.gain;
        this.to = data.to;
    }
}

class MarketBuyEvent implements EventType{
    resource: string;
    amount: number;
    cost: number;
    from : string;

    constructor(data: any) {
        this.resource = data.resourceType;
        this.amount = Math.ceil(data.amountExecuted);
        this.cost = data.cost;
        this.from = data.from;
    }
}

// Implement the TransferResourceEvent class
class TransferResourceEvent implements EventType{
    from: string;
    to: string;
    resource: string;
    amount: number;

    constructor(data: any, country?: Country) {
        this.from = data.from;
        this.to = data.to;
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
    id: string;
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
    moneyHistory: Map<string, number>;
    flag: string;

    constructor(data: any) {
        this.agent = new Agent(data.agent);
        this.color = data.color;
        this.money = Math.ceil(data.money)
        this.history = data.history.map((eventData: any) => new CountryEvent(eventData, this));
        this.moneyHistory = new Map<string, number>(Object.entries(data.moneyHistory));
        this.flag = data.flag;
    }


    getCountryPopulation(simulation: Simulation): number {
        let territories = this.getTerritories(simulation);
        return territories.reduce((accumulator, currentTerritory) => {
            return accumulator + currentTerritory.habitants
        },0);
    }

    getTotalStocks(simulation: Simulation): Map<string, number> {
        let territories = this.getTerritories(simulation);
        const result: Map<string, number> = new Map();

        // Iterate over each territory
        territories.forEach((territory) => {
            //console.log(territory.stock);
            // Iterate over each entry in the territory's stock map
            territory.stock.forEach((value, key) => {
                // Add the value to the result map or update if the key already exists
                result.set(key, (result.get(key) || 0) + value);
            });
        });

        return result;
    }

    getTerritories(simulation: Simulation): Territory[] {
        return simulation.territories.filter(
            (territory) => territory.country?.agent.id === this.agent.id
        )
    }

    // Méthode pour récupérer le stockHistory de tous les territoires
    getAllTerritoriesStockHistory(simulation: Simulation): Map<number, Map<string, number>> {
        const allTerritoriesStockHistory = new Map<number, Map<string, number>>();

        // Parcourir tous les territoires du pays
        this.getTerritories(simulation).forEach((territory) => {
            // Ajouter le stockHistory du territoire à la carte globale
            territory.stockHistory.forEach((history, timestamp) => {
                if (!allTerritoriesStockHistory.has(timestamp)) {
                    allTerritoriesStockHistory.set(timestamp, new Map<string, number>());
                }

                const existingHistory = allTerritoriesStockHistory.get(timestamp);
                if (existingHistory) {
                    // Mettre à jour le stockHistory global en ajoutant les valeurs du territoire
                    history.forEach((value, key) => {
                        existingHistory.set(key, (existingHistory.get(key) || 0) + value);
                    });
                }
            });
        });

        return allTerritoriesStockHistory;
    }

    getAllTerritoriesHabitantsHistory(simulation: Simulation): Map<string, number> {
        const allTerritoriesHabitantsHistory = new Map<string, number>();

        // Parcourir tous les territoires du pays
        this.getTerritories(simulation).forEach((territory) => {
            // Ajouter le stockHistory du territoire à la carte globale
            territory.habitantsHistory.forEach((value, key) => {
                allTerritoriesHabitantsHistory.set(key, (allTerritoriesHabitantsHistory.get(key) || 0) + value);
            });
        });

        return allTerritoriesHabitantsHistory;
    }

}

class Territory {
    x: number;
    y: number;
    name: string;
    variations: Variation[];
    stock: Map<string, number>;
    stockHistory: Map<number, Map<string, number>>;
    habitants: number;
    habitantsHistory: Map<string, number>;
    country: Country;

    constructor(data: any) {
        this.x = data.x;
        this.y = data.y;
        this.name = data.name;
        this.variations = data.variations.map((variationData: any) => new Variation(variationData));
        this.stock = new Map<string, number>(Object.entries(data.stock));
        this.habitants = data.habitants;
        this.country = new Country(data.country);
        this.stockHistory = new Map<number, Map<string, number>>();
        for (const key in data.stockHistory) {
            const innerMap = new Map<string, number>(Object.entries(data.stockHistory[key]));
            this.stockHistory.set(Number(key), innerMap);
        }
        this.habitantsHistory = new Map<string, number>(Object.entries(data.habitantsHistory));
    }
}

class MarketInteraction {
    dateTransaction: string;
    resourceType: string;
    amount: number;
    price: number;
    buyer: string ;
    seller: string;

    constructor(data: any) {
        this.dateTransaction = data.dateTransaction;
        this.resourceType = data.resourceType;
        this.amount = data.amount;
        this.price = data.price;
        this.buyer = data.buyer.agent.name;
        this.seller = data.seller.agent.name;
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
    consumptionByHabitant: Map<string, number>;

    constructor(data: any) {
        this.market = new Market(data.market);
        this.consumptionByHabitant = new Map<string, number>(Object.entries(data.consumptionByHabitant))
    }
}


class Simulation {
    secondByDay: number;
    environment: Environment;
    territories: Territory[];
    countries: Map<string, Country>;
    currentDay: number;

    constructor(data: any) {
        console.log(data)
        this.secondByDay = data.secondByDay;
        this.environment = new Environment(data.environment);
        this.territories = data.environment.world.territories.map((territoryData: any) => new Territory(territoryData));
        this.countries = new Map<string, Country>(
            Object.entries(data.environment.countries).map(([countryKey, countryData]: [string, any]) => [
                countryKey,
                new Country(countryData),
            ])
        );
        this.currentDay = data.currentDay;
    }
}

export { Simulation, Territory, Country, Resource, Variation, Agent, MarketInteraction, Market, Environment, TransferResourceEvent, MarketBuyEvent, MarketSellEvent };
