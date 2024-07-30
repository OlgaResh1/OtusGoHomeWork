FROM golang:1.22 as build

ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY . ${CODE_DIR}

RUN go mod download

CMD ["go", "test", "-tags", "integration", "./integration_tests/..."]