FROM golang:alpine
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY go.mod go.sum ./
COPY app.env .
# Download all the dependencies
RUN go get -d -v ./...
RUN go install -v ./...
# This container exposes port 8080 to the outside world
EXPOSE 9000
CMD ["go-api-gateway"]