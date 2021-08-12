FROM golang:alpine

# RUN mkdir /app

WORKDIR /app
COPY . .
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon

# ADD go.mod .
# ADD go.sum .

# RUN go mod download
# ADD . .

# RUN go get github.com/githubnemo/CompileDaemon

EXPOSE ${PORT}

ENTRYPOINT CompileDaemon --build="go build -o main_app" --command=./main_app
