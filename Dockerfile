FROM golang:1.21.2 AS builder
ARG VERSION
ARG GIT_HASH

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
	-ldflags="-X 'mmlabel.gitlab.com/mm-printing-backend/version.Version=${VERSION}' -X 'mmlabel.gitlab.com/mm-printing-backend/version.GitHash=${GIT_HASH}' -X 'mmlabel.gitlab.com/mm-printing-backend/version.GoVersion=1.16.5'" \
	-a -o /server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /server ./

COPY ./resources/ws ./ws-docs
COPY ./migrations /migrations

RUN chmod +x ./server
ENTRYPOINT ["./server"]
