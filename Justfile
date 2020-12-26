help:
  @just --list


# Build Docker image
build-image:
  docker build -t hello-gin:latest .


# Build Docker image
run-image:
  docker run --rm -p 8000:8080 hello-gin:latest


# Start development setup
start:
  ( cd src; nodemon --exec 'go run *.go || exit 1' -e go,html,css,js )


# TODO: The command below leaves few duplicate keys, remove those.
# Generate file-ext -> mime-type mapping source file
mimetypes:
  echo "package assets" > ./src/assets/mime.go
  echo "" >> ./src/assets/mime.go
  echo "var MimeTypes = map[string]string{" >> ./src/assets/mime.go
  wget -qO- http://svn.apache.org/repos/asf/httpd/httpd/trunk/docs/conf/mime.types     \
    | egrep -v ^#                                                                      \
    | awk '{ for (i=2; i<=NF; i++) {print "    \"."$i"\": \""$1"\","} }'               \
    | sort                                                                             \
    >> ./src/assets/mime.go
  echo "}" >> ./src/assets/mime.go

build-dev:
  docker build --tag hello-gin-server:dev -f Dockerfile-dev .


dev +args="":
  docker run                                         \
    --rm                                             \
    --name dev                                       \
    -p 8080:8080                                     \
    -e HOST=0.0.0.0                                  \
    -e PORT=8080                                     \
    -e GIN_MODE=debug                                \
    -e RESOURCES=/app                                \
    -w /app/src                                      \
    -v $(pwd)/src:/app/src:cached                    \
    -v $(pwd)/assets:/app/assets:cached              \
    -v $(pwd)/templates:/app/templates:cached        \
    hello-gin-server:dev {{ args }}


sh:
  docker run                                         \
    --rm                                             \
    -it                                              \
    -p 8080:8080                                     \
    -e HOST=0.0.0.0                                  \
    -e PORT=8080                                     \
    -e GIN_MODE=debug                                \
    -e RESOURCES=/app                                \
    -w /app/src                                      \
    -v $(pwd)/src:/app/src:cached                    \
    -v $(pwd)/assets:/app/assets:cached              \
    -v $(pwd)/templates:/app/templates:cached        \
    hello-gin-server:dev bash
