# Gopolitical

[![CI: Test](https://github.com/BaptisteBuvron/gopolitical/actions/workflows/test.yml/badge.svg)](https://github.com/BaptisteBuvron/gopolitical/actions/workflows/test.yml)

- [Subject](https://docs.google.com/document/d/1H8QpU5dTMkJEEb2nTqgMNJ84rH7QNalC8CqPTC4qPV8)
- [Repository](https://github.com/BaptisteBuvron/gopolitical)

## Installation

### Server

Install [Go](https://golang.org/doc/install).

Clone the repository:

```bash
go install github.com/BaptisteBuvron/gopolitical/server@v1.0.0
```

Run the server:

```bash
server
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
