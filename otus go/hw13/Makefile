.PHONY: test
PROJECTNAME="hw13"

help: Makefile
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

## build: Build application
build:
	@for type in "api" "scheduler" "notifier" ; do \
		CGO_ENABLED=0 go build -o ./bin/$$type ./cmd/$$type/ ; \
	done

## build-api: Build api
build-api:
		CGO_ENABLED=0 go build -o ./bin/api ./cmd/api/ ;

## build-notifier: Build notifier
build-notifier:
		CGO_ENABLED=0 go build -o ./bin/notifier ./cmd/notifier/ ;

## build-scheduler: Build scheduler
build-scheduler:
		CGO_ENABLED=0 go build -o ./bin/scheduler ./cmd/scheduler/ ;

## test: run integration tests
test:
	cd deployment ;\
	test_status=0 ;\
	docker-compose -f docker-compose.test.yaml up --build --exit-code-from tests --abort-on-container-exit|| test_status=$$? ;\
	cd .. ;\
	exit $$test_status ;\
