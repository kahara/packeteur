FROM golang:1.21.2-bullseye as build

RUN mkdir /workdir
COPY go.* /workdir/
COPY *.go /workdir/

WORKDIR /workdir
RUN go build -o pcktr .

FROM gcr.io/distroless/base-debian12 as production

COPY --from=build /workdir/pcktr /

CMD ["/pcktr"]
