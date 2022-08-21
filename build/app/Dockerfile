# Собираем в гошке
FROM golang:1.17.3 as build

ENV BIN_FILE /opt/banner/app
ENV CODE_DIR /go/src/
ENV CGO_ENABLED=0

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/rotator/*


FROM alpine:3.9

LABEL SERVICE="bannersRotator"
LABEL MAINTAINERS="yur.nazarov@x5.ru"

ENV BIN_FILE "/opt/banner/app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}