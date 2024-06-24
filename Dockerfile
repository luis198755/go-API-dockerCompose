# Use an official Go runtime as a parent image compatible with Raspberry Pi
#FROM golang:latest
# Start from the golang base image
FROM golang:alpine

# Set the working directory in the container
WORKDIR /usr/src/app

# Copy the local code to the container
COPY ./api.go ./api.go

COPY ./program.json ./program.json

COPY ./index.html ./index.html

# Download any dependencies
RUN go mod init example.com/api

# (Assuming you have a go.mod and go.sum file in your project)
#RUN go mod download

# Copy the go.mod and go.sum files
#COPY go.mod go.sum ./

# Compile the program
RUN go build -o api

# Expose port 8080 to the outside world
EXPOSE 8080

# Step 7: Command to run the executable
CMD ["./api"]
