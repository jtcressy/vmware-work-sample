FROM golang:1.17-alpine AS build-env

WORKDIR /go/src/vmware-work-sample

COPY go.sum go.mod ./
RUN go mod download

COPY . .

ARG TARGETARCH

RUN CGO_ENABLED=0 GOARCH=$TARGETARCH go install

FROM alpine:3.15
COPY --from=build-env /go/bin/* /usr/local/bin
ENTRYPOINT ["/usr/local/bin/vmware-work-sample"]