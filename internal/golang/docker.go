package golang

import (
	"fmt"
	"grosf-gh/internal/util"
)

func createDockerfile() error {
	path := "./" + Name + "/Dockerfile"

	// Create the Dockerfile
	dockerfile := fmt.Sprintf(`FROM golang:1.16-alpine

WORKDIR /app
	
COPY go.mod ./
COPY go.sum ./
RUN go mod download
	
COPY . ./
	
RUN go build -o /%s
	
EXPOSE 8080
	
CMD [ "/%s" ]`, Name, Name)

	// Write the Dockerfile if it doesn't exist
	return util.CreateFile(path, dockerfile)
}
