# alpine:3.6
FROM ubuntu:trusty-20161101

COPY .dist/generic-bank /boa/bin/generic-bank
COPY html/ /boa/html/

EXPOSE 8085

WORKDIR /boa/bin

CMD ["/boa/bin/generic-bank"]
