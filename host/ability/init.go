package ability

import (
	_ "github.com/sbgayhub/golem/host/ability/cdn"
	_ "github.com/sbgayhub/golem/host/ability/chatroom"
	_ "github.com/sbgayhub/golem/host/ability/contact"
	_ "github.com/sbgayhub/golem/host/ability/favor"
	_ "github.com/sbgayhub/golem/host/ability/label"
	_ "github.com/sbgayhub/golem/host/ability/message"
)

import contactability "github.com/sbgayhub/golem/host/ability/contact"
import chatroomability "github.com/sbgayhub/golem/host/ability/chatroom"

// Inject 注入能力层
func Inject() {
}

func Destroy() {
	contactability.Destroy()
	chatroomability.Destroy()
}
