
.PHONY: all
all: setup-ngrok i-issuer-fe i-verifier-fe r-issuer-be r-issuer-fe r-verifier-fe open-browser

i-issuer-fe: # install issuer frontend
	yarn --cwd examples/kycAge/issuerClient

i-verifier-fe: # install verifier frontend
	yarn --cwd examples/kycAge/verifierClient

r-issuer-fe: # run issuer frontend
	sleep 2;
	yarn --cwd examples/kycAge/issuerClient dev

r-verifier-fe: # run verifier frontend
	yarn --cwd examples/kycAge/verifierClient dev

r-issuer-be: # run issuer backend
	cd issuer/ && go run cmd/main.go

open-browser:
	sleep 4;
	./scripts/open_browser.sh;

setup-ngrok:
	./scripts/run_ngrok_update_cfg.sh

stop-ngrok:
	killall ngrok




