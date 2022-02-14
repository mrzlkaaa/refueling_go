FROM golang:1.17
WORKDIR /refuel

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./refueling ./refueling
COPY ./refueling/cmd/.env .env

RUN go build -o /refueling ./refueling/cmd
CMD [ "/refueling" ]