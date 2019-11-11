FROM gcr.io/kaniko-project/executor:debug-v0.10.0

ENV HOME /root
ENV USER root
ENV SSL_CERT_DIR=/kaniko/ssl/certs
ENV DOCKER_CONFIG /kaniko/.docker/

# add the wrapper which acts as a drone plugin
COPY drone-kaniko /kaniko/drone-kaniko
COPY plugin.sh /kaniko/plugin.sh
ENTRYPOINT [ "/kaniko/drone-kaniko" ]
