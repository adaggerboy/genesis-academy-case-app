FROM golang:1.21-alpine AS build
LABEL stage=gobuild
RUN apk update --no-cache
WORKDIR /src
COPY . .
WORKDIR /src/cmd/currency-api
RUN go build -o /bin/currency-api

FROM scratch
COPY --from=build /bin/currency-api /bin/currency-api
CMD ["/bin/currency-api"]