package contactability

import "github.com/sbgayhub/golem/sdk/contact"

var cache map[string]*contact.Contact

func get(username string) *contact.Contact {
	return cache[username]
}
