FROM golang:1.24.6

WORKDIR /app 

COPY . . 

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /Practicum_final

ARG TODO_PORT=7540
ENV TODO_PORT=${TODO_PORT}

CMD ["/Practicum_final"]