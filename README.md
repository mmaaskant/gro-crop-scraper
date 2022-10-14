# Gro - Crop Scraper
> Collects and compiles data from various crop seed suppliers and other sources of crop data.

This data will be used within the Gro project to offer initial crop data, so it can be enriched with growth predictions.

## Features
Currently, the scraper supports 2 steps; crawling and filtering.
Both these steps are run by default unless specified otherwise like so:
```bash
go main.go # Run all steps
go main.go --crawl # Only run Crawl
go main.go --filter # Only run Filter
```
It is also possible to only target a specific supplier by providing the supplier's name as a flag, like so:
```bash
go main.go --burpee # Run all steps for the Burpee supplier.
```
These 2 types of flags can be combined freely.
### Steps
Steps behave according to the Chain of Responsibility pattern and pass along their processed data to the next step.
All steps support concurrency and the amount of concurrent GoRoutines for each step can be configured in the .env file.
#### Crawl
The Crawl step attempts to find all product pages within the supplier's domain and saves their contents to a MongoDB database.
Currently, it saves just the initial response and only its HTML.
#### Filter
The Filter step pulls all data that the Crawler has saved and attempts to extract relevant data.
The relevancy of this data is determined by a set of Criteria that are passed along each seed supplier's config.

## Planned Features
This project is currently still under development and the following features are currently on the roadmap.
### Map Step
Once the data has been filtered, it has to be mapped to a universal format so any following steps do not have to deal with
the complexity of different data sets.
### Compile Step
Once data has been mapped, any duplicate crops from different suppliers should be merged.
Any duplicate fields would be compared to one another and resolved accordingly.
### gRPC API
Once all data has been compiled, it should be made available through a gRPC API.
This API should offer a basic range of filters and paging.

## Deployment
This project is currently not configured for deployment, however it provides a docker-compose setup purely meant for development.
The Go container's module dependencies are synced locally to `<project_root_dir>/dev_vendor` so an IDE can access them easily.
The MongoDB container's volume is located at `<project_root_dir>/compose/mongodb/data`.

## Licensing
The code in this project is licensed under an GPL-3.0 license.