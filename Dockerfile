FROM alpine:3.19

COPY helm-docs /usr/bin/

WORKDIR /helm-docs

ENTRYPOINT ["helm-docs"]
