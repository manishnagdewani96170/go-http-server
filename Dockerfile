# build stage
FROM golang:latest as builder

RUN mkdir /app

ADD *.go /app/

# Set the Current Working Directory inside the container
WORKDIR /app


# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a  -o server .


# final stage
FROM alpine:latest  

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/server .

ENTRYPOINT [ "./server" ]

# arguments that can be overridden
CMD [ "3", "300" ]

RUN CGO_ENABLED=0 GOOS=linux go build -a -o client .


