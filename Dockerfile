FROM golang:alpine as builder

RUN echo xx > /hello

FROM alpine as release

COPY --from=builder /hello hello
