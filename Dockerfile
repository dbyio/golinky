FROM golang:alpine AS build 
COPY . /src
WORKDIR /src
RUN go get && CGO_ENABLED=0 go build -o golinky 

FROM alpine
COPY --from=build /src/golinky /
COPY config.yaml /
USER daemon
ENTRYPOINT ["/golinky", "config.yaml"]