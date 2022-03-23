FROM golang:latest as buildcontainer

WORKDIR /build
RUN git clone https://github.com/Eeshaan-rando/sipb.git
WORKDIR /build/sipb/server
RUN git checkout prithvi-dev
RUN go build -o sipb

FROM alpine:latest as sipb
LABEL maintainer="Prithvi Vishak <prithvivishak@gmail.com>"

COPY --from=buildcontainer /build/sipb/server/sipb /
COPY --from=buildcontainer /build/sipb/webpages/ /var/www/html/
RUN apk add --no-cache gcompat

EXPOSE 80/tcp

WORKDIR /etc/sipb
CMD [ "/sipb" ]
