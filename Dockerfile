FROM alpine

COPY helm-docs /usr/bin/

ENTRYPOINT ["helm-docs"]
