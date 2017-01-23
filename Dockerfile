FROM ubuntu:xenial

RUN apt-get update && \
apt-get -y upgrade && \
rm -rf /var/cache/apt/*

COPY squidbot-flowdock-frontend /usr/bin/squidbot-flowdock-frontend

CMD ["/usr/bin/squidbot-flowdock-frontend"]
