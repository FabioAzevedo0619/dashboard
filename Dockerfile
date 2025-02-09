FROM golang:1.23-alpine
 
WORKDIR /usr/src/app
 
# Effectively tracks changes within your go.mod file
COPY ./src/go.mod ./
RUN go mod download

COPY ./src/ ./
RUN go build -v -o /usr/local/bin/app ./...

# Expose the port the service listens on
EXPOSE 8080
 
CMD ["app"]