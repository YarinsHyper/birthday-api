#build stage
FROM golang:alpine AS builder
ENV GO111MODULE=on
RUN apk add --no-cache git make
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

#final stage
FROM golang:alpine
COPY --from=builder /go/src/app/api-gateway /api-gateway
LABEL Name=api-gateway Version=0.0.1
WORKDIR /go/src/app
EXPOSE 9000
ENTRYPOINT ["/api-gateway"]