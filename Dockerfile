FROM alpine:3.9
#FROM ubuntu:trusty-20161101

COPY .dist/generic-bank /boa/bin/generic-bank
COPY ui/dist/ /boa/html/

EXPOSE 8085

WORKDIR /boa/bin

RUN adduser -D gbank
USER gbank

CMD ["/boa/bin/generic-bank"]
