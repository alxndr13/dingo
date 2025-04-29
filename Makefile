.PHONY: build run test integration_tests test_google_decryptor

build:
	@echo "🏗️ Building the app and storing in ./bin/dingo"
	mkdir -p bin
	go build -o bin/dingo

run:
	@echo "🏗️ Running the app with defaults"
	go run . --logmode human

test:
	@echo "🏗️ Running unit tests"
	go test -v ./... -cover

integration_tests: test_google_decryptor
	@echo "✅ Finished Integration Testing"

test_google_decryptor:
	@if gcloud auth application-default print-access-token > /dev/null 2>&1; then \
		echo "✅ Logged in with Application Default Credentials"; \
	else \
		echo "❌ Not logged in with Application Default Credentials, login using 'gcloud auth application-default login'"; \
	fi

	@echo "🏗️ Running Google decryptor integration tests"
	@go run . --logmode human --basepath ./test/google_decryptor/data/base \
		--overlaypath ./test/google_decryptor/data/overlays/dev \
		--templatepath ./test/google_decryptor/templates \
		--decryptor google
