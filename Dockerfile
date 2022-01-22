FROM alpine:latest

COPY helm-docs /usr/bin/

WORKDIR /helm-docs

ENTRYPOINT ["helm-docs"]
