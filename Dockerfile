FROM --platform=$BUILDPLATFORM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETARCH TARGETOS

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o metasploit-db ./cmd

FROM metasploitframework/metasploit-framework:latest

RUN apk update && apk add --no-cache \
    tcpdump \
    tshark

WORKDIR /app

COPY --from=builder /app/metasploit-db /app/metasploit-db

COPY --from=builder /app/scripts/ /app/

RUN mkdir -p /app/results

RUN pip3 install -r /app/requirements.txt

ENTRYPOINT ["/bin/sh", "-c", "ruby /usr/src/metasploit-framework/msfrpcd -U msf -P dL0rHLep -p 55552 -S false -f & ./metasploit-db"]