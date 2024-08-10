  
FROM golang:1.22 as build

ENV BIN_FILE /opt/calendar/calendar-sender-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY . ${CODE_DIR}

RUN go mod download

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar_sender/*

FROM alpine:3.20

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar-sender"
LABEL MAINTAINERS="4olgamail@gmail.com"

ENV BIN_FILE "/opt/calendar/calendar-sender-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/sender_config.toml
COPY ./configs/sender_config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} --config ${CONFIG_FILE}
