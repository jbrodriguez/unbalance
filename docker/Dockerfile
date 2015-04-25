FROM gliderlabs/alpine:3.1
MAINTAINER Juan B. Rodriguez <jbrodriguez.io@apertoire.net>

RUN apk --update add wget

RUN wget https://raw.github.com/jbrodriguez/unbalance/downloads/unbalance-0.2.tar.gz -O - | tar -xzf - -C /unbalance && \
	chmod +x /unbalance/unbalance

WORKDIR /unbalance

ENTRYPOINT ["./unbalance"]