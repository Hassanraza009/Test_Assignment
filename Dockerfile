# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY . .
RUN go build

# Run stage
FROM golang:1.20-alpine
WORKDIR /app
COPY --from=builder /app/test .
# Install curl in the final image
RUN apk --no-cache add curl

EXPOSE 8001
CMD [ "/app/test" ]