module github.com/Brycicle777/pokedexcli

go 1.23.0

replace internal/mapcommands v0.0.0 => ./internal/mapcommands

replace internal/pokecache v0.0.0 => ./internal/pokecache

replace internal/pokecollect v0.0.0 => ./internal/pokecollect

require (
	internal/mapcommands v0.0.0
	internal/pokecache v0.0.0
	internal/pokecollect v0.0.0
)
