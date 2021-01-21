FROM golang:1.15.2-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /app .
FROM scratch AS bin
COPY --from=build /app /app
CMD ["./battleship", "start"]