FROM alpine

WORKDIR /app

COPY ../../cmd/web /app/interface

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' > /etc/timezone

ENV LANG C.UTF-8

EXPOSE 9000

ENTRYPOINT ["./interface"]
