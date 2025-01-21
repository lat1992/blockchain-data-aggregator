FROM alpine:3

RUN apk add --no-cache make go

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

ADD . /app/
WORKDIR /app

RUN make

FROM scratch

ADD . /app/
WORKDIR /app

COPY --from=0 /app/build/blockchain-data-aggregator /app/blockchain-data-aggregator
COPY --from=0 /app/datas /app/datas

CMD ["/bin/blockchain-data-aggregator"]
