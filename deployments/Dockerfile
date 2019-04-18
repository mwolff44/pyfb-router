# Build Stage
FROM lacion/alpine-golang-buildimage:1.11 AS build-stage

LABEL app="build-pyfb-router"
LABEL REPO="https://github.com/mwolff44/pyfb-router"

ENV PROJPATH=/go/src/github.com/mwolff44/pyfb-router

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/mwolff44/pyfb-router
WORKDIR /go/src/github.com/mwolff44/pyfb-router

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/mwolff44/pyfb-router"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/pyfb-router/bin

WORKDIR /opt/pyfb-router/bin

COPY --from=build-stage /go/src/github.com/mwolff44/pyfb-router/bin/pyfb-router /opt/pyfb-router/bin/
RUN chmod +x /opt/pyfb-router/bin/pyfb-router

# Create appuser
RUN adduser -D -g '' pyfb-router
USER pyfb-router

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/pyfb-router/bin/pyfb-router"]
