FROM golang:1.23.4 AS builder

WORKDIR /build
COPY . .
RUN go build .

FROM scratch
COPY --from=builder /build/openrct2-webui /openrct2-webui

CMD [ "/openrct2-webui" ]
