FROM golang:1.20
LABEL maintainer="oysteinq@stud.ntnu.no"

# Copy files to gopath
#COPY ./api /go/src/app/cmd
#COPY ./cmd /go/src/app/webhooks
#COPY ./internal /go/src/app/internal
#COPY ./res /go/src/app/res
#COPY ./go.mod /go/src/app/go.mod
COPY . /go/src/app/

# Set working directory
WORKDIR /go/src/app/

# Build project
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' ./cmd/app/main.go

# Expose port
EXPOSE 8080

# Set image entry point
CMD ["./main"]