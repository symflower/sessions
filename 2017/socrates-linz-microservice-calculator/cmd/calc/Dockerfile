FROM ubuntu:16.04

RUN addgroup --gid 1000 service
RUN adduser --uid 1000 --gid 1000 --disabled-password --gecos '' service
USER service

CMD ["/bin/calc"]

ADD calc /bin/calc
