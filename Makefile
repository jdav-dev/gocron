all: vendor test vet

clean:
	rm -rf vendor

test:
	go test -v $$(glide novendor)

vendor:
	glide up --strip-vendor

vet:
	go vet $$(glide novendor)
