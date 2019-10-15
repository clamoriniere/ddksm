FROM golang:1.13 as builder

WORKDIR /go/src/app
COPY . .
RUN GO111MODULE="on" go build -o /ddksm .

# Final image
FROM gcr.io/distroless/static

USER nobody
COPY --from=builder /ddksm  /ddksm
ENTRYPOINT [ "/ddksm" ]
