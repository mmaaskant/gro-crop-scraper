FROM golang:1.19.1 AS dev

RUN apt-get update && apt-get install netcat -y

WORKDIR /src/gro-crop-scraper
COPY . ./
RUN chmod 0755 ./compose/go/entrypoint.sh
CMD ["./compose/go/entrypoint.sh"]