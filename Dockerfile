FROM golang:latest as buildcontainer

WORKDIR /build
# RUN git clone https://github.com/Eeshaan-rando/sipb.git
COPY ./webpages /build/sipb/webpages/
COPY ./server /build/sipb/server/
WORKDIR /build/sipb/server
RUN go build -o sipb

FROM alpine:latest as sipb
LABEL maintainer="Prithvi Vishak <prithvivishak@gmail.com>"

COPY --from=buildcontainer /build/sipb/server/sipb /
COPY --from=buildcontainer /build/sipb/server/config.yaml.docker /etc/sipb/config.yaml
COPY --from=buildcontainer /build/sipb/webpages/ /var/www/html/
RUN mkdir -p /var/www/bin && ln -s /var/www/bin /var/www/html/bin
RUN apk add --no-cache gcompat

EXPOSE 80/tcp

WORKDIR /etc/sipb
CMD [ "/sipb" ]
