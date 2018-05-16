FROM golang:alpine as builder

ARG x=hello
RUN echo ${x} > /hello

FROM alpine as release

COPY --from=builder /hello hello
