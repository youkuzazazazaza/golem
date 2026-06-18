package sync

import (
	contactability "github.com/sbgayhub/golem/host/ability/contact"
	messageapi "github.com/sbgayhub/golem/host/api/message"
	"github.com/sbgayhub/golem/host/plugin"
	pluginsdk "github.com/sbgayhub/golem/sdk/plugin"
)

func handleContact(contacts []*messageapi.ModifyContact) {
	publishChangeEvents(contactability.HandleModifyContacts(contacts))
}

func publishChangeEvents(events []*pluginsdk.ChangeEvent) {
	for _, event := range events {
		if event == nil {
			continue
		}
		plugin.Publish(&pluginsdk.Event{
			Topic:   event.Domain + "::" + event.Action,
			Payload: &pluginsdk.Event_Change{Change: event},
		})
		for _, change := range event.Changes {
			if change.GetPath() == "" {
				continue
			}
			plugin.Publish(&pluginsdk.Event{
				Topic:   event.Domain + "::" + event.Action + "::" + change.GetPath(),
				Payload: &pluginsdk.Event_Change{Change: event},
			})
		}
	}
}
