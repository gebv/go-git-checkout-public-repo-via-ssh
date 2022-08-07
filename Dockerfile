FROM golang:1.17-buster as builder
    WORKDIR /go/src/github.com/gebv/go-git-checkout-public-repo-via-ssh
    COPY go.mod .
    COPY go.sum .
    RUN go mod download
    COPY . ./
    # if using docker experimental mode
    # RUN --mount=type=cache,mode=0777,target=/root/.cache/go-build go build -v -o ./bin/app
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -v -o ./bin/app.bin ./

FROM alpine:3.11 as configured-env

    RUN apk update && apk add --no-cache openssh-client git ca-certificates tzdata make curl && update-ca-certificates

    RUN mkdir -p ~/.ssh
    RUN chmod 700 ~/.ssh
    RUN echo -e "Host *\n\tStrictHostKeyChecking no\n\tPasswordAuthentication no\n" > ~/.ssh/config
    RUN ssh-keygen -q -t rsa -N '' -f ~/.ssh/id_rsa
    RUN chmod 600 ~/.ssh/id_rsa
    # RUN ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts


FROM configured-env as try-checkout-via-ssh
    COPY --from=builder /go/src/github.com/gebv/go-git-checkout-public-repo-via-ssh/bin/app.bin /app.bin

    CMD eval `ssh-agent -s` && /app.bin
