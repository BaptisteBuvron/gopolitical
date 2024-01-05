# Gopolitical

[![CI: Test](https://github.com/BaptisteBuvron/gopolitical/actions/workflows/test.yml/badge.svg)](https://github.com/BaptisteBuvron/gopolitical/actions/workflows/test.yml)

Gopolitical is a multi-agent project developped to simulate trade and diplomatic relationships between countries.

- [Subject](https://docs.google.com/document/d/1H8QpU5dTMkJEEb2nTqgMNJ84rH7QNalC8CqPTC4qPV8)
- [Repository](https://github.com/BaptisteBuvron/gopolitical)

## Installation

### Server

Install [Go](https://golang.org/doc/install).

#### With git clone

```bash
git clone https://github.com/BaptisteBuvron/gopolitical
sudo go install .
# Start-Process powershell -Verb runAs -ArgumentList "cd $(Get-Location); go install"
```

#### Without git clone

```bash
go install github.com/BaptisteBuvron/gopolitical/server@v1.1.1
```

#### Customization of the Simulation Instance

It is possible to change the `/server/resources/data.json` file that will be interpreted by the simulation by default at launch. This file instructs the simulation about various countries, territories, stock variations, consumptions, etc.

It is possible to generate a `data.json` file using the Python script located in `/server/resources/generate_data.py`.

Run the server:

```bash
cd server
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

## Loop

```mermaid
flowchart TB
    s1[Territory\ngenerate ressources and add them to country]
    s2[Countries\nsell their object by sending request to environnement]
    s3[Countries\nbuy what they need by sending request to environnement]
    s4[Environnement\nmake transaction]
    s5[Countries\nevaluate their new relations]

    s1 --> s2
    s2 --> s3
    s3 --> s4
    s4 --> s5
    s5 --> s1
```
