FROM golang:alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN cd cmd && go build -o diary.be .

FROM alpine:latest
RUN apk add --no-cache sqlite-libs ca-certificates tzdata
# Set the timezone (optional)
ENV TZ=Europe/Prague
WORKDIR /root/
COPY --from=builder /app/cmd/diary.be .
COPY --from=builder /app/webapp ./webapp
EXPOSE 8080
CMD ["./diary.be", "server"]
