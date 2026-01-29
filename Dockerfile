# use official Golang image
FROM golang:1.17-alpine3.13

# set working directory
WORKDIR /app

# Copy the source code
COPY . . 

# Download and install the dependencies + Build the Go app
RUN go get -d -v ./... \
    && go build -o api .

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./api"]