FROM alpine:latest

RUN mkdir /app

COPY lisneterApp /app

CMD [ "/app/lisneterApp" ]