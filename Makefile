.PHONY: integration-test

integration-test:
	go test -v -run ^TestFullIntegration$$
