FROM golang:1.17
WORKDIR /diary

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY ./diary ./diary
COPY ./diary/cmd/.env .env

RUN go build -o /diary-built ./diary/cmd

CMD [ "/diary-built" ]