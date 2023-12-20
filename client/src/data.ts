

export const data = {
  "secondByDay": 24.0,
  "resources": [
    { "name": "petrol", "price": 30},
    { "name": "water", "price": 1},
    { "name": "food", "price": 2}
  ],
  "countries": [
    { "name": "Russie", "id": "ru", "color": "FF0000", "money": 100.0 },
    { "name": "United-state", "id": "us", "color": "0000FF", "money": 100.0 }
  ],
  "relations": [
    {"country1": "ru", "country2": "us", "kind": "enemy"}
  ],
  "territories": [
    {
      "x": 0,
      "y": 0,
      "country": "ru",
      "variations": [
        { "name": "petrol", "value": 10 },
        { "name": "water", "value": 5 }
      ]
    },
    {
      "x": 0,
      "y": 1,
      "country": "ru",
      "variations": [
        { "name": "petrol", "value": 10 },
        { "name": "food", "value": 5 }
      ]
    },
    {
      "x": 1,
      "y": 1,
      "country": "ru",
      "variations": [
        { "name": "petrol", "value": 10 },
        { "name": "water", "value": 5 }
      ]
    },
    {
      "x": 1,
      "y": 0,
      "country": "us",
      "variations": [
        { "name": "petrol", "value": 10 },
        { "name": "water", "value": 5 }
      ]
    },
    {
      "x": 2,
      "y": 0,
      "country": "us",
      "variations": [
        { "name": "food", "value": 10 },
        { "name": "water", "value": 5 }
      ]
    },
    {
      "x": 2,
      "y": 1,
      "country": "us",
      "variations": [
        { "name": "food", "value": 10 },
        { "name": "water", "value": 5 }
      ]
    },
  ]
};