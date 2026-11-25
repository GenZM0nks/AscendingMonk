# AscendingMonk
A repository to migrate and further develop the legacy project WhoKowsLegacy, and to mantain the website ascendingmonk.dk.

> Visit our search engine at: [**https://Ascendingmonk.dk**](https://ascendingmonk.dk/)

The website is a simple search engine, with the possiblility of registereing an account. We also provide an integrated weather forecast, with data supplied by openMeteo.

# Running the code:

On a *NIX system, navigate to the root directory and run
```bash
DB_PATH="test-database-name.(sqlite|db)" go run ./cmd/main
```