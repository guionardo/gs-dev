// Package consts defines HTTP header names and URL path constants used by
// the pad service's REST endpoints.
package consts

const (
	// HeaderPadID is the HTTP header carrying the pad ID.
	HeaderPadID = "Gs-Dev-Pad-Id"
	// HeaderPadTTL is the HTTP header carrying the pad TTL in seconds.
	HeaderPadTTL = "Gs-Dev-Pad-TTL"
	// HeaderPadAPIKey is the HTTP header carrying the API key.
	HeaderPadAPIKey = "Gs-Dev-Pad-API-Key"
	// PadURL is the base URL path for pad operations.
	PadURL = "/api/v1/pads"
	// PadGetURL is the URL path for retrieving a pad by ID.
	PadGetURL = PadURL + "/{id}"
	// PadDeleteURL is the URL path for deleting a pad by ID.
	PadDeleteURL = PadURL + "/{id}"
	// PadCreateURL is the URL path for creating a new pad.
	PadCreateURL = PadURL
	// PadHomeURL is the URL path for the pad root.
	PadHomeURL = PadURL + "/"
	// PadReadyURL is the health check URL path.
	PadReadyURL = PadURL + "/ready"
)
