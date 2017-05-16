FROM alpine

MAINTAINER Renan Tateoka <renan@vendhq.com>

# RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /opt/lab-go-api

EXPOSE 8080

# @todo update volumne to only point to web / dist folder
VOLUME ["/opt/lab-go-api/app"]

COPY ./etc /opt/lab-go-api/etc

COPY lab-go-api /bin/lab-go-api

ENTRYPOINT ["/bin/lab-go-api"]
