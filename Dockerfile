# Use an official Golang runtime as a parent image
FROM ubuntu:latest

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Install necessary dependencies
RUN apt-get update && \
    apt-get install -y \
    # curl \
    openssh-client \
    unzip \
    sudo \
    man-db \
    libgl1-mesa-dev \
    libxrandr2 \
    libx11-xcb-dev \
    libwayland-dev \
    libxkbcommon-dev \
    libxcursor-dev \
    libxi-dev \
    libxinerama-dev \
    xorg-dev \
    gcc

# Copy the current directory contents into the container at /go/src/app
COPY app .

# Install dependencies
# RUN go get fyne.io/fyne/v2@v2.0.0

# Set the environment variable for cross-compiling to Windows amd64
# ENV GOOS=linux
# ENV GOARCH=amd64
# ENV CGO_ENABLED=1
# ENV CC=mingw64
# RUN go mod tidy

# Build the application
# RUN go build -o app

# Set the entry point for the container
# CMD ["./app"]
CMD ["sleep", "infinity"]
