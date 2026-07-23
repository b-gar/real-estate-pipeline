# ----------------------------------------------------
# Stage 1: Build the Go Binary
# ----------------------------------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compile the binary.
# CGO_ENABLED=0 ensures it is a statically linked binary (no external C dependencies)
# -ldflags="-w -s" strips debugging symbols to make the file size drastically smaller!
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o pipeline .


# ----------------------------------------------------
# Stage 2: Create the Micro-Container (Distroless)
# ----------------------------------------------------
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

COPY --from=builder /app/pipeline /pipeline

CMD ["/pipeline"]
