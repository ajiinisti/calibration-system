FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o calibration-system

FROM alpine
WORKDIR /app
COPY --from=build /src/calibration-system /app
ENTRYPOINT ["/app/calibration-system"]