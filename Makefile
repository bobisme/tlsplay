CERTS := $(shell find certs -name "*.pem")

.PHONY: run
run: certs
	go run main.go & sleep 3; ./main.rb

certs: $(CERTS)
	./make-certs.sh
