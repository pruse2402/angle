FROM golang:1.13 as builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /go/workspace
COPY go.mod go.sum ./
RUN echo '[url "ssh://git@github.com/"]' >> /root/.gitconfig \
   && echo 'insteadOf = https://github.com/' >> /root/.gitconfig \
   && go mod download
COPY . .
RUN go build -o app -a -installsuffix cgo .

# ----------------

FROM alpine:latest

ARG APP_PATH=/go/src/angle
WORKDIR ${APP_PATH}
RUN addgroup -S angle_group \
    && adduser -D -S -h ${APP_PATH} -s /sbin/nologin -G angle_group angle_user \
    && apk --no-cache add ca-certificates
USER angle_user

COPY --from=builder /go/workspace/app .
COPY config ./config

EXPOSE 8085
ENTRYPOINT ["./app", "serve"]