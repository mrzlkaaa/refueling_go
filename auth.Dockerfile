FROM golang:1.17
WORKDIR /auth

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./auth ./auth
COPY ./auth/cmd/.env .env

RUN go build -o /authserv ./auth/cmd
CMD [ "/authserv" ]