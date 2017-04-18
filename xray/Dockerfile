FROM alpine:3.5

RUN ln -s /dev/stdout /logging
RUN apk add --no-cache ca-certificates wget unzip libc6-compat && ln -s /lib/ /lib64
RUN wget -O /tmp/xray.zip https://s3.amazonaws.com/aws-xray-assets.us-east-1/xray-daemon/aws-xray-daemon-linux-1.x.zip && unzip -d /usr/bin/ /tmp/xray.zip

COPY cfg.yaml /etc/xray/config.yaml

CMD /usr/bin/xray --log-file /logging --config /etc/xray/config.yaml
