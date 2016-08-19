package uuid

import twinj "github.com/twinj/uuid"

// GetUUID returns a cached (via file) or new UUID.
func GetUUID(_ string) string {
	// TODO: Use _ argument as a file cache with which we re-use the cached uuid
	// if present in the file or generate (and cache) a new one.  This will ensure
	// we don't generate a new uuid on each startup.
	return twinj.NewV4().String()
}
