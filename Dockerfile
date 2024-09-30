FROM golang:1.23 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /SimpleAPRSBot

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /SimpleAPRSBot /SimpleAPRSBot

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/SimpleAPRSBot"]