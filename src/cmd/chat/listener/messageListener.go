package listener

import (
	"fmt"

	authTypes "local/auth/types"
	"local/message/types"
)

type MessageListener authTypes.UserCap

func (user MessageListener) MessageReceived(msg types.Message) {
	fmt.Printf("%s: %s\n", msg.From, msg.Text)
}
