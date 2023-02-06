FROM alpine:3.9

RUN apk add --no-cache bash

ENV GOOSE_PATH /usr/local/bin/goose
RUN wget -q -O $GOOSE_PATH https://github.com/pressly/goose/releases/latest/download/goose_linux_x86_64 \
    && chmod +x $GOOSE_PATH

COPY wait-for-it.sh wait-for-it.sh
RUN chmod +x wait-for-it.sh

CMD ["goose"]