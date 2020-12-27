FROM alpine:3.12

RUN apk add --no-cache \
	nginx \
	bash
WORKDIR /src
COPY . .

ENTRYPOINT ["./start.sh"]
