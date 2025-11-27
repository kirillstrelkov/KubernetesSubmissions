FROM google/cloud-sdk:alpine

RUN apk add --no-cache postgresql-client

WORKDIR /usr/src/app

COPY backup.sh .
RUN chmod +x backup.sh

CMD ["./backup.sh"]
