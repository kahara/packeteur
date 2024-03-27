FROM golang:1.21.2-bullseye as build

WORKDIR /workdir

COPY go.* /workdir/
RUN go mod download

COPY *.go /workdir/
RUN go build -o pcktr .

FROM gcr.io/distroless/base-debian12 as production

COPY --from=build /workdir/pcktr /

CMD ["/pcktr"]
