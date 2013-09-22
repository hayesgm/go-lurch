# Lurch API Server
#
# Run a Lurch API server

FROM centos
MAINTAINER Geoffrey Hayes <hayesgm@gmail.com>

RUN mkdir -p /srv/lurch
RUN wget https://github.com/hayesgm/lurch/releases/download/0.0.1/lurch.linux -O /srv/lurch/lurch
RUN chmod +x /srv/lurch/lurch

ENTRYPOINT ["/srv/lurch/lurch"]

EXPOSE 9119