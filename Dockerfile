FROM golang:1.17-alpine
# Note: This app should run fine on an M1 Mac,
#       if that's what you're using (I can't verify), but...
#       we may need to set a '--platform=linux/amd64' flag here if it doesn't run

WORKDIR /usr/src/github.com/kellyfay94/fetch

COPY . .

RUN go mod vendor && \
    go build

ENTRYPOINT ["sh","run_checks.sh"]