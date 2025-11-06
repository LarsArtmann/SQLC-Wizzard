# Add TypeSpec compilation to build system
types:
	npx @typespec/compiler compile api/typespec.tsp --emit @typespec/go-server
	cd generated && go mod init sqlc-wizard-types && go mod tidy

.PHONY: types