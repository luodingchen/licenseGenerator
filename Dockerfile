FROM golang:alpine as dev
MAINTAINER dev
WORKDIR /app
ADD . /app
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o app

FROM alpine:latest as prod
WORKDIR /app
ENV LISTEN_ADDRESS ":80"
EXPOSE 80
COPY --from=dev /app/configs /app/configs
COPY --from=dev /app/files   /app/files
COPY --from=dev /app/app /app/
CMD /app/app
