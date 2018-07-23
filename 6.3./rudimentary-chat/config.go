package main


const (
	// ReadBytes delimiter 'end of text'
	EOT byte = '\x03'

	// API-like protocolbuffer messages
	//
	SUBSCRIBE             = iota
	Unsubscribe
	Publish
	DisplayText
	DisplaySubscription
	DisplayUnsubscription
)

// TODO: export some constants to config file
const (
	IpAddr string = "127.0.0.1"
	Port   string = "22365"
)
