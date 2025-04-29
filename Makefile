.PHONY: build run

build:
	mkdir -p bin
	go build -o bin/dingo

run:
	go run . --logmode human

test:
	go test -v ./... -cover

test_google_decryptor:
	@if gcloud auth application-default print-access-token > /dev/null 2>&1; then \
		echo "Logged in with Application Default Credentials"; \
	else \
		echo "Not logged in with Application Default Credentials"; \
	fi
	@go run . --logmode human --basepath ./test/google_decryptor/data/base \
		--overlaypath ./test/google_decryptor/data/overlays/dev \
		--templatepath ./test/google_decryptor/templates \
		--decryptor google

