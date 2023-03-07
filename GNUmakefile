default: build

locals:
	echo "locals {" > locals.tf
	cat oktawave/config.go|grep -A1000 "<export>"|grep -B1000 "</export>"|head -n -1|tail -n +2|sed "s/DICT_//"|sed "s/\/\//#/" >> locals.tf
	echo "}" >> locals.tf

test:
	go test -v ./...

testacc:
	TF_ACC=1 go test -v ./...

build:
	go build -o terraform-provider-oktawave

.PHONY: build test testacc locals
