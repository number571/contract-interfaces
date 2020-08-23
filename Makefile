.PHONY: default build clean
default: build
build: contract.sol deploy.go gclient.go
	solc --overwrite --abi --bin contract.sol -o build
	mkdir -p contracts
	./abigen --bin=./build/WorldSkills.bin --abi=./build/WorldSkills.abi --pkg=contract --out=./contracts/Contract.go
	go build -o deploy deploy.go
	go build -o client client.go values.go
	go build -o gclient gclient.go values.go
clean: 
	rm -rf build/ contracts/
	rm deploy gclient client
