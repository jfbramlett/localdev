FROM 118139069697.dkr.ecr.us-west-1.amazonaws.com/hub/golang:1.16.8-alpine3.14

# Install git+ SSL c certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints
RUN apk update && apk add --no-cache git ca-certificates tzdata bash curl jq && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/github.com/splice/platform/

COPY ./go.mod $GOPATH/src/github.com/splice/platform/go.mod
COPY ./go.sum $GOPATH/src/github.com/splice/platform/go.sum

RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get -u github.com/cosmtrek/air

COPY ./infra/cmd/localdev/scripts/build.sh /tmp/build.sh
COPY ./infra/cmd/localdev/scripts/run.sh /tmp/run.sh
COPY ./infra/cmd/localdev/scripts/start.sh /tmp/start.sh

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 40000
EXPOSE 8080
EXPOSE 8085

ARG SERVICE

WORKDIR $GOPATH/src/github.com/splice/platform/infra/cmd/localdev

HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD nc -zv localhost 8080 || exit 1

CMD "/tmp/start.sh"
