FROM alpine:3 AS build

RUN apk add --no-cache make go

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

ADD . /app/
WORKDIR /app

RUN make

FROM alpine:3

ADD . /app/
WORKDIR /app

COPY --from=build /app/build/blockchain-data-aggregator /app/
COPY --from=build /app/datas /app/

CMD ["/app/blockchain-data-aggregator"]
