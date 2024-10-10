module pokecollect

go 1.23.0

replace internal/pokecache v0.0.0 => ../pokecache

replace internal/mapcommands v0.0.0 => ../mapcommands

require (
	internal/pokecache v0.0.0
	internal/mapcommands v0.0.0
)
