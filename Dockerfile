FROM golang:alpine AS builder
ENV GO111MODULE="on" 
RUN mkdir -p /src
WORKDIR /src
ADD go.mod go.mod
ADD go.sum go.sum
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /src/ddksm .

# Final image
FROM gcr.io/distroless/static
USER nobody
COPY --from=builder /src/ddksm  /ddksm
ENTRYPOINT [ "/ddksm" ]
