# Builder

FROM golang:1.23-alpine AS builder
RUN apk add --update curl git make && apk cache clean

ARG MODULE_NAME=github.com/Kapeland/task-Avito

COPY Makefile /home/${MODULE_NAME}/Makefile
COPY go.mod /home/${MODULE_NAME}/go.mod
COPY go.sum /home/${MODULE_NAME}/go.sum

WORKDIR /home/${MODULE_NAME}

COPY . /home/${MODULE_NAME}

RUN make build

# Service

FROM alpine:latest AS server
RUN apk add --update curl && apk cache clean
ARG MODULE_NAME=github.com/Kapeland/task-Avito
WORKDIR /root/

COPY --from=builder /home/${MODULE_NAME}/bin/main .
COPY --from=builder /home/${MODULE_NAME}/configs ./configs

COPY --from=builder /home/${MODULE_NAME}/internal/storage/db/migrations ./postgres/migrations

RUN chown root:root main

CMD ["sh", "-c", " ./main"]