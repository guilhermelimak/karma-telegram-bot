default:
	docker run --rm \
  --volume $(shell pwd):/go/src/github.com/guilhermelimak/karma-telegram-bot \
  --workdir /go/src/github.com/guilhermelimak/karma-telegram-bot \
  quay.io/deis/go-dev:latest \
  glide up && go run *.go
