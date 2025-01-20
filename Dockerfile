FROM alpine:3

RUN apk add --no-cache make go

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

ADD . /app/
WORKDIR /app

RUN make

FROM scratch

COPY --from=0 /app/build/blockchain-data-aggregator /bin/blockchain-data-aggregator

CMD ["/bin/blockchain-data-aggregator"]
