
// Define interfaces for your data structure
export interface Resource {
    name: string;
    price: number;
}

export interface Country {
    name: string;
    id: string;
    color: string;
    money: number;
}

export interface Relation {
    country1: string;
    country2: string;
    kind: string;
}

export interface Variation {
    name: string;
    value: number;
}

export interface Territory {
    x: number;
    y: number;
    country: string;
    variations: Variation[];
}

export interface Data {
    secondByDay: number;
    resources: Resource[];
    countries: Country[];
    relations: Relation[];
    territories: Territory[];
}

