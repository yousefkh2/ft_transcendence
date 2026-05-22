FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./backend/cmd/api

FROM alpine:3.22

WORKDIR /app

COPY --from=build /out/api /app/api
COPY --from=build /src/frontend/sprint1-poc /app/frontend/sprint1-poc

EXPOSE 8080

CMD ["/app/api"]
