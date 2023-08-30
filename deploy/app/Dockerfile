FROM golang:1.19 as build
WORKDIR /app
COPY . .
RUN go build -o /bin/spiralMatrix ./cmd/*

FROM golang:1.19
COPY --from=build /bin/spiralMatrix spiralMatrix
COPY --from=build /app/.env .env
CMD ["./spiralMatrix"]
