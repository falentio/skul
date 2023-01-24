FROM node:16-bullseye as node

WORKDIR /web

COPY web/package.json .
COPY web/package-lock.json .
RUN npm ci

COPY web/static static
COPY web/src src
COPY web/.* .
COPY web/*.js .
COPY web/*.json .
RUN npm run sync
RUN npm run build

FROM golang:1.19.2-bullseye as go

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY internal internal
COPY cmd cmd
COPY web/web.go web/
COPY --from=node /web/dist web/dist
RUN CGO_ENABLED=0 go build -o skul ./cmd/skul

FROM debian:bullseye as prod

WORKDIR /app
COPY --from=go /app/skul ./skul

EXPOSE 8080
CMD "./skul"

