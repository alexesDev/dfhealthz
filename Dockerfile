FROM golang:1.16.2-alpine3.13 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -extldflags '-static'" -o dfhealthz .

FROM scratch
COPY --from=0 /app/* ./
CMD ["./dfhealthz"]
