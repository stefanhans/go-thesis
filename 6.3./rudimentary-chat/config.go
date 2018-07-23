package main


const (
	// ReadBytes delimiter 'end of text'
	EOT byte = '\x03'

	// API-like protocolbuffer messages
	//
	SUBSCRIBE             = iota
	UNSUBSCRIBE
	PUBLISH
	DISPLAY_TEXT
	DISPLAY_SUBSCRIPTION
	DISPLAY_UNSUBSCRIPTION
)

// TODO: export some constants to config file
const (
	IpAddr string = "127.0.0.1"
	Port   string = "22365"
)

