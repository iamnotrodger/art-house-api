MONGO_DB = art-house
ARTIST_COLLECTION = artwork
ARTWORK_COLLECTION = artwork
EXHIBITION_COLLECTION = exhibition

run: ./cmd/art-house-api/main.go 
	@go run ./cmd/art-house-api/main.go

test: 
	@go test ./...


seed: ./seed/artists.json ./seed/artworks.json ./seed/exhibitions.json
	@mongo $(MONGO_DB) --eval "db.$(ARTIST_COLLECTION).drop()"
	@mongo $(MONGO_DB) --eval "db.createCollection('$(ARTIST_COLLECTION)')"
	@mongoimport --db $(MONGO_DB) --collection $(ARTIST_COLLECTION) --file ./seed/artists.json --jsonArray

	@mongo $(MONGO_DB) --eval "db.$(ARTWORK_COLLECTION).drop()"
	@mongo $(MONGO_DB) --eval "db.createCollection('$(ARTWORK_COLLECTION)')"
	@mongoimport --db $(MONGO_DB) --collection $(ARTWORK_COLLECTION) --file ./seed/artworks.json --jsonArray

	@mongo $(MONGO_DB) --eval "db.$(EXHIBITION_COLLECTION).drop()"
	@mongo $(MONGO_DB) --eval "db.createCollection('$(EXHIBITION_COLLECTION)')"
	@mongoimport --db $(MONGO_DB) --collection $(EXHIBITION_COLLECTION) --file ./seed/exhibitions.json --jsonArray