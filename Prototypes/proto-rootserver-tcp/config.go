package main

const (
	// TODO: Debug '\x04' as delimiter
	// ReadBytes delimiter
	EOF byte = '\x08'

	// API-like protocolbuffer messages for UDP and TCP
	//
	Join    = iota // TCP
	Members        // TCP
	Update         // UDP
	Leave          // UDP
)

// TODO: export some constants to config file
const (
	IpAddr string = "127.0.0.1"
	Port   string = "22365"
)
