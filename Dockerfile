
FROM golang:latest AS builder



# Copy the code from the host and compile it
WORKDIR $GOPATH/src/OpenStreetmapRouting
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM scratch
COPY --from=builder /app ./
EXPOSE 8080
ENTRYPOINT ["./app"]