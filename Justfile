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
