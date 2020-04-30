FROM golang:1.13

WORKDIR /go/src/app

COPY . .

RUN make tester

ENTRYPOINT ["/go/src/app/build/tester"]