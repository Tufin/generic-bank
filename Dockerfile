FROM alpine:3.7
#FROM ubuntu:trusty-20161101

COPY .dist/generic-bank /boa/bin/generic-bank
COPY ng-app/generic-bank/dist/ /boa/html/

EXPOSE 8085

WORKDIR /boa/bin

CMD ["/boa/bin/generic-bank"]
