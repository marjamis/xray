FROM alpine:3.7

RUN ln -s /dev/stdout /logging && \
    apk add --no-cache ca-certificates wget unzip libc6-compat && ln -s /lib/ /lib64 && \
    wget -O /tmp/xray.zip https://s3.dualstack.us-east-2.amazonaws.com/aws-xray-assets.us-east-2/xray-daemon/aws-xray-daemon-linux-2.x.zip && \
    unzip -d /usr/bin/ /tmp/xray.zip
COPY cfg.yaml /etc/xray/config.yaml

ENTRYPOINT [ "/usr/bin/xray", "--config", "/etc/xray/config.yaml" ]
