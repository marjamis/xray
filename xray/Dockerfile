FROM alpine:3.5

RUN ln -s /dev/stdout /logging
RUN apt-get update -y && \
  apt-get -y install ca-certificates wget unzip
RUN wget -O /tmp/xray.zip https://s3.amazonaws.com/aws-xray-assets.us-east-1/xray-daemon/aws-xray-daemon-linux-1.x.zip && unzip -d /usr/bin/ /tmp/xray.zip

COPY cfg.yaml /usr/bin/config.yaml

CMD /usr/bin/xray --log-file /logging --config /usr/bin/config.yaml
