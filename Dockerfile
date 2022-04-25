FROM alpine:3.15

COPY helm-docs /usr/bin/

WORKDIR /helm-docs

ENTRYPOINT ["helm-docs"]
