import json
import random
from colorama import Fore, Style  # Ajout de colorama

def generate_flashy_color():
    # Génération d'une couleur plus flashi avec une intensité plus élevée
    return f"#{random.choice('3456789ABCDEF')}{random.choice('3456789ABCDEF')}{random.choice('3456789ABCDEF')}"

def generate_data(num_countries, territories_per_country=10):
    data = {
        "secondByDay": 2.0,
        "resources": [
            {"name": "petrol", "price": 10},
            {"name": "water", "price": 1.5},
            {"name": "food", "price": 2},
            {"name": "armement", "price": 5}
        ],
        "consumptionsByHabitant": [
            {"name": "petrol", "value": 0.2},
            {"name": "water", "value": 0.5},
            {"name": "food", "value": 0.5}
        ],
        "countries": [],
        "relations": [],
        "territories": []
    }

    with open("flags_code.json", "r") as f:
        flags_code = json.load(f)

        unique_coordinates = set()
        selected_country_names = set()

        for i in range(num_countries):
            # Sélection d'un nom de pays au hasard depuis le fichier JSON
            available_country_names = set(flags_code) - selected_country_names
            if not available_country_names:
                # Si tous les noms de pays ont été utilisés, vous pouvez choisir une autre approche, ou lever une exception
                raise ValueError("Insufficient unique country names available.")

            country_id = random.choice(list(available_country_names))
            country_name = flags_code[country_id]
            selected_country_names.add(country_name)

            country_color = generate_flashy_color()[1:]
            country_money = random.uniform(200, 400)
            data["countries"].append({
                "name": country_name,
                "id": country_id,
                "color": country_color,
                "money": country_money,
                "flag": get_country_flag_url(country_name)
            })

            # Liste des coordonnées des territoires du pays en cours
            country_territory_coordinates = []

            for j in range(territories_per_country):
                territory_name = f"Territory{j + 1}"

                # Generate unique coordinates (adjacent if possible)
                while True:
                    # Sélection aléatoire des coordonnées d'un territoire existant du même pays
                    existing_territory_coords = random.choice(
                        country_territory_coordinates) if country_territory_coordinates else None
                    if existing_territory_coords:
                        # Sélection aléatoire d'une direction parmi les voisins possibles (Nord, Sud, Est, Ouest)
                        direction = random.choice(["North", "South", "East", "West"])
                        dx, dy = {"North": (0, 1), "South": (0, -1), "East": (1, 0), "West": (-1, 0)}[direction]
                        territory_x, territory_y = existing_territory_coords[0] + dx, existing_territory_coords[1] + dy
                    else:
                        territory_x = random.randint(0, 20)
                        territory_y = random.randint(0, 20)

                    # Vérification d'unicité des coordonnées et validité des coordonnées
                    if (territory_x, territory_y) not in unique_coordinates and 0 <= territory_x <= 20 and 0 <= territory_y <= 20:
                        unique_coordinates.add((territory_x, territory_y))
                        break

                country_territory_coordinates.append((territory_x, territory_y))

                territory_habitants = random.randint(5, 50)
                territory_stocks = [
                    {"name": resource["name"], "value": random.randint(5, 50)}
                    for resource in data["resources"]
                ]
                territory_variations = [
                    {"name": resource["name"], "value": random.randint(5, 15)}
                    for resource in data["resources"]
                ]

                data["territories"].append({
                    "x": territory_x,
                    "y": territory_y,
                    "country": country_id,
                    "name": territory_name,
                    "habitants": territory_habitants,
                    "stocks": territory_stocks,
                    "variations": territory_variations
                })

    return data

def get_country_flag_url(country_name):
    with open("flags_code.json", "r") as f:
        flags_code = json.load(f)
        for code, name in flags_code.items():
            if name == country_name:
                return f"https://www.flagcdn.com/256x192/{code}.png"
                break
    return "https://upload.wikimedia.org/wikipedia/commons/2/2e/Unknown_flag_-_European_version.png"

def save_to_json(data, filename="data.json"):
    with open(filename, "w") as f:
        json.dump(data, f, indent=2)

if __name__ == "__main__":
    num_countries = int(input("Enter the number of countries: "))
    territories_per_country = 10  # You can change this if needed

    simulation_data = generate_data(num_countries, territories_per_country)
    save_to_json(simulation_data)
    print(f"{Fore.GREEN}Data saved to data.json{Style.RESET_ALL}")
