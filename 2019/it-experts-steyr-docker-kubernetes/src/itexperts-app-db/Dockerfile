FROM ubuntu:18.04

RUN addgroup --gid 1000 service
RUN adduser --uid 1000 --gid 1000 --disabled-password --gecos '' service
USER service

CMD ["/bin/itexperts-app-db"]

ADD itexperts-app-db /bin/itexperts-app-db
