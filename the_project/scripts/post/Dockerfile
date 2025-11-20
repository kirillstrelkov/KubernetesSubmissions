FROM alpine:latest

RUN apk add --no-cache bash curl grep gawk

COPY create_new_post.sh .

ENTRYPOINT ["./create_new_post.sh"]
