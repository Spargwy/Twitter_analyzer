FROM golang:1.17-alpine
WORKDIR /code
COPY . ./      
RUN go mod download      
RUN go get ./...
RUN go build -o /web

EXPOSE 8080

CMD ["/web"]