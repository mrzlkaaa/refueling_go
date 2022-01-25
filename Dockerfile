FROM golang:1.13
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /refuelDiary ./cmd/

CMD [ "/refuelDiary" ]