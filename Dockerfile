#
# Build:
#

FROM golang:1.15-buster AS build

WORKDIR /app

COPY ./src/go.mod  .
COPY ./src/go.sum  .

RUN go mod download

COPY ./src  ./

RUN CGO_ENABLED=0    \
    GOOS=linux       \
    GOARCH=amd64     \
    go build -a -o /server

#
# Dist:
#

FROM gcr.io/distroless/static-debian10 AS dist

COPY --from=build /server /server
COPY ./assets     ./assets
COPY ./templates  ./templates

ENV HOST=0.0.0.0
ENV PORT=8080
ENV GIN_MODE=release

CMD ["/server"]
