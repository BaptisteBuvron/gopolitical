# Gopolitical

[![CI: Test](https://github.com/BaptisteBuvron/gopolitical/actions/workflows/test.yml/badge.svg)](https://github.com/BaptisteBuvron/gopolitical/actions/workflows/test.yml)

Gopolitical is a multi-agent project developped to simulate trade and diplomatic relationships between countries.

- [Presentation](https://github.com/BaptisteBuvron/gopolitical/blob/527804f57558f642c1d0124dd35ae677e1a682bb/AI30%20-%20Pr%C3%A9sentation.pdf)
- [Subject](https://docs.google.com/document/d/1H8QpU5dTMkJEEb2nTqgMNJ84rH7QNalC8CqPTC4qPV8)
- [Repository](https://github.com/BaptisteBuvron/gopolitical)

## Changement depuis la dernière version

- L'environnement n'est plus lancé dans un thread.
- La simulation exécute les actions des agents après que l'ensemble des agents aient donné leurs actions
- Moins de bugs de concurrence
- Les pays choisissent en priorité les pays proches à attaquer

## Installation

### Server

Install [Go](https://golang.org/doc/install).

#### With git clone

```bash
git clone https://github.com/BaptisteBuvron/gopolitical
sudo go install .
# Start-Process powershell -Verb runAs -ArgumentList "cd $(Get-Location); go install"
```

#### Customization of the Simulation Instance

It is possible to change the `/server/resources/data.json` file that will be interpreted by the simulation by default at launch. This file instructs the simulation about various countries, territories, stock variations, consumptions, etc.

It is possible to generate a `data.json` file using the Python script located in `/server/resources/generate_data.py`.

Run the server:

```bash
cd server
# $env:GOPOLITICAL_DEBUG=0
go run .
```

Run tests:

```bash
go test '-coverprofile=coverage.txt' -v ./...
go tool cover '-html=coverage.txt'
```

### Client

```bash
cd client
npm i
npm run start
```
