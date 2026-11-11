.PHONY: gen-proto gen-go-proto gen-android-proto

gen-proto: gen-go-proto gen-android-proto

gen-go-proto:
	protoc -I=proto \
		--go_out=backend \
		--go_opt=module=github.com/mushibana/glacier/backend \
		--go-grpc_out=backend \
		--go-grpc_opt=module=github.com/mushibana/glacier/backend proto/*.proto

gen-android-proto:
	cd android && ./gradlew :app:generateReleaseProto
