run: ./cmd/art-house-api/main.go 
	@go run ./cmd/art-house-api/main.go

seed: ./cmd/seed/seed.go ./cmd/seed/data/artworks.json ./cmd/seed/data/artists.json ./cmd/seed/data/exhibitions.json
	@go run ./cmd/seed/seed.go
	@echo "Data base seeded"

test: 
	@go test ./...
