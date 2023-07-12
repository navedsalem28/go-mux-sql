#FROM golang:1.20-alpine
FROM golang:1.20

# Set destination for COPY
WORKDIR /app

# Print Directory Path
RUN pwd && ls

# Copy the source code.
COPY . .

# Update dependencies
RUN apt-get update

# Download Go modules
RUN go mod download

# Build
RUN go build -o /go-mux-sql

# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# Run executable file
CMD [ "/go-mux-sql" ]














