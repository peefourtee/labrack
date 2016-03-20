FROM golang:1.6

RUN mkdir -p $GOPATH/src/github.com/peefourtee/labrack
ADD . $GOPATH/src/github.com/peefourtee/labrack

WORKDIR $GOPATH/src/github.com/peefourtee/labrack
RUN go get github.com/tools/godep

#RUN go install github.com/peefourtee/labrack/cmd/labrack
EXPOSE 8000
CMD go run cmd/labrack/* --mock-devices=1 --sample=2s
