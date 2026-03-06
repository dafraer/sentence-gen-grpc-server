generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
run:
	go run cmd/main.go
push:
	git add . && git commit -m "$(m)" && git push
docker:
	#Pass version using v variable
	sudo docker build  --platform linux/amd64 -t dafraer/sentence-gen-grpc-server:$(v) .
	docker push dafraer/sentence-gen-grpc-server:$(v)