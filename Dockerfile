FROM golang:1.21-alpine AS builder

RUN apk update && apk --no-cache add bash git make

WORKDIR /usr/src

COPY ["go.mod","go.sum","./"]

RUN go mod tidy && go mod download

COPY . .

# build
RUN go build -o ./bin/app ./cmd/app/main.go 

# create alpine
FROM alpine:latest AS runner

RUN apk update

# copy binary from builder
COPY --from=builder /usr/src/bin/app .
COPY ./config ./config

CMD ["./app", "--config=./config/dev.yaml"]