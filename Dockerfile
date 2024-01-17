FROM --platform=linux/amd64 golang:1.18

WORKDIR /app

COPY . ./
RUN go mod download


RUN go build -o /server main.go

CMD [ "/server" ]
