FROM golang:1.15.6 AS tools-yj
RUN CGO_ENABLED=0 GOPROXY=https://goproxy.cn go get github.com/sclevine/yj && go install github.com/sclevine/yj
WORKDIR /app
COPY remove-duplicated-runner /app
RUN CGO_ENABLED=0 GOPROXY=https://goproxy.cn go build -o /go/bin/remove-duplicated-runner

FROM ubuntu:18.04 AS tools-envsubst
RUN apt-get update -y && apt-get install -y gettext-base

FROM gitlab/gitlab-runner:v13.6.0
RUN apt-get update -y && apt-get install -y jq
COPY --from=tools-yj /go/bin/yj /bin/yj
COPY --from=tools-yj /go/bin/remove-duplicated-runner /bin/remove-duplicated-runner
COPY --from=tools-envsubst /usr/bin/envsubst /bin/envsubst
WORKDIR /app
COPY . .
ENTRYPOINT  /app/entrypoint.sh
