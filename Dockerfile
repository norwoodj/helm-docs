FROM alpine:3.11.5

COPY helm-docs /usr/bin/

WORKDIR /helm-docs

ENTRYPOINT ["helm-docs"]
