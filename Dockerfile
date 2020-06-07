FROM golang:latest

RUN mkdir /home/treee
WORKDIR /home/treee
COPY go.mod /home/treee
COPY go.sum /home/treee
RUN go mod download
COPY . /home/treee

RUN go build
CMD [ "./treee" ]