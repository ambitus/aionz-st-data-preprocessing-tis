FROM icr.io/ibmz/golang:1.21.4-bullseye

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY *.go ./

RUN go build -o /go_preprocess

EXPOSE 8080

ENTRYPOINT ["/go_preprocess"]