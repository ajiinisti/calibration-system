FROM golang:alpine as build
RUN apk update && apk add --no-cache git
WORKDIR /src
COPY . .
RUN go mod tidy
RUN go build -o calibration-system

FROM alpine
WORKDIR /app
COPY --from=build /src/calibration-system /app
COPY --from=build /src/utils/templates /app/utils/templates
ENTRYPOINT ["/app/calibration-system"]