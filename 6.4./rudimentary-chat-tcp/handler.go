package main

import (
	"fmt"
	"log"
	"net"

	"bitbucket.org/stefanhans/go-thesis/6.4./rudimentary-chat-tcp/chat-group"
)

// Map as link between message and its handler
var actionMap = map[chatgroup.Message_MessageType]func(*chatgroup.Message, net.Addr) error{
	chatgroup.Message_SUBSCRIBE_REQUEST:   handleSubscribeRequest,
	chatgroup.Message_SUBSCRIBE_REPLY:     handleSubscribeReply,
	chatgroup.Message_UNSUBSCRIBE_REQUEST: handleUnsubscribeRequest,
	chatgroup.Message_UNSUBSCRIBE_REPLY:   handleUnsubscribeReply,
	chatgroup.Message_PUBLISH_REQUEST:     handlePublishRequest,
	chatgroup.Message_PUBLISH_REPLY:       handlePublishReply,
	chatgroup.Message_MEMBERLIST_REQUEST:  handleMemberlistRequest,
	chatgroup.Message_MEMBERLIST_REPLY:    handleMemberlistReply,
	chatgroup.Message_LEADERLIST_REQUEST:  handleLeaderlistRequest,
	chatgroup.Message_LEADERLIST_REPLY:    handleLeaderlistReply,
	chatgroup.Message_ECHO_REQUEST:        handleEchoRequest,
	chatgroup.Message_ECHO_REPLY:          handleEchoReply,
}

func handleSubscribeRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	// Check subscriber for uniqueness
	for _, recipient := range selfMemberList {
		if recipient.Name == message.Sender.Name {
			return fmt.Errorf("name %q already used", message.Sender.Name)
		}
		if recipient.Ip == message.Sender.Ip && recipient.Port == message.Sender.Port {
			return fmt.Errorf("address %s:%s already used by %s", recipient.Ip, recipient.Port, recipient.Name)
		}
	}

	// Add subscriber
	log.Printf("Add subscriber: %v\n", message.Sender)
	selfMemberList = append(selfMemberList, message.Sender)
	log.Printf("Current members registered: %v\n", selfMemberList)

	err := publishMessage(message, chatgroup.Message_SUBSCRIBE_REPLY)
	if err != nil {
		return fmt.Errorf("failed to publish Message_SUBSCRIBE_REPLY: %v", err)
	}

	return nil
}

// Display new member
func handleSubscribeReply(message *chatgroup.Message, addr net.Addr) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s (%s:%s) has joined>", message.Sender.Name, message.Sender.Ip, message.Sender.Port))

	return nil
}

func handleUnsubscribeRequest(message *chatgroup.Message, addr net.Addr) error {

	log.Printf("Unregister: %v\n", message.Sender)

	// Remove subscriber
	for i, s := range selfMemberList {
		if s.Name == message.Sender.Name {
			selfMemberList = append(selfMemberList[:i], selfMemberList[i+1:]...)
			break
		}
	}
	log.Printf("Current members registered: %v\n", selfMemberList)

	err := publishMessage(message, chatgroup.Message_UNSUBSCRIBE_REPLY)
	if err != nil {
		return fmt.Errorf("failed to publish Message_UNSUBSCRIBE_REPLY: %v", err)
	}

	return nil
}

func handleUnsubscribeReply(message *chatgroup.Message, addr net.Addr) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("<%s has left>", message.Sender.Name))

	return nil
}

func handlePublishRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	log.Printf("Publish from %v: %q\n", message.Sender.Name, message.Text)

	err := publishMessage(message, chatgroup.Message_PUBLISH_REPLY)
	if err != nil {
		return fmt.Errorf("failed to publish Message_Message_PUBLISH_REPLY: %v", err)
	}

	return nil
}

func handlePublishReply(message *chatgroup.Message, addr net.Addr) error {

	// Append text message in "messages" view
	displayText(fmt.Sprintf("%s: %s", message.Sender.Name, message.Text))

	return nil
}

func handleMemberlistRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	// Address string to reply requestor
	requestor := message.Sender.Ip + ":" + message.Sender.Port

	// Change message for reply
	message.MsgType = chatgroup.Message_MEMBERLIST_REPLY
	message.Sender = selfMember
	message.MemberList.Member = selfMemberList

	// Send message to requestor
	err := sendMessage(message, requestor)
	if err != nil {
		return fmt.Errorf("failed send reply: %v", err)
	}
	return nil
}

func handleMemberlistReply(message *chatgroup.Message, addr net.Addr) error {

	selfMemberList = message.MemberList.Member

	displayText(fmt.Sprintf("<CMD_SYNC_REPLY>: reply received from %v:%v - call \\list",
		message.Sender.Ip, message.Sender.Port))
	return nil
}

func handleLeaderlistRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	// Address string to reply requestor
	requestor := message.Sender.Ip + ":" + message.Sender.Port

	// Change message for reply
	message.MsgType = chatgroup.Message_LEADERLIST_REPLY
	message.Sender = selfMember
	message.MemberList.Member = getLeaderlist()

	// Send message to requestor
	err := sendMessage(message, requestor)
	if err != nil {
		return fmt.Errorf("failed send reply: %v", err)
	}
	return nil
}

func getLeaderlist() []*chatgroup.Member {

	leaderlist := make([]*chatgroup.Member, len(selfMemberList))

	for _, member := range selfMemberList {
		if member.Leader {
			leaderlist = append(leaderlist, member)
		}
	}

	if len(leaderlist) == 0 {
		selfMember.Leader = true
		leaderlist = append(leaderlist, selfMember)
	}

	return leaderlist
}

func handleLeaderlistReply(message *chatgroup.Message, addr net.Addr) error {

	selfMemberList = message.MemberList.Member

	displayText(fmt.Sprintf("<CMD_SYNC_REPLY>: reply received from %v:%v - call \\list",
		message.Sender.Ip, message.Sender.Port))
	return nil
}

func handleEchoRequest(message *chatgroup.Message, addr net.Addr) error {

	// Update remote IP address, if changed
	updateRemoteIP(message, addr)

	// Address string to reply requestor
	requestor := message.Sender.Ip + ":" + message.Sender.Port

	// Change message for reply
	message.MsgType = chatgroup.Message_LEADERLIST_REPLY
	message.Sender = selfMember
	message.MemberList.Member = selfMemberList

	// Send message to requestor
	err := sendMessage(message, requestor)
	if err != nil {
		return fmt.Errorf("failed send reply: %v", err)
	}
	return nil
}

func handleEchoReply(message *chatgroup.Message, addr net.Addr) error {

	selfMemberList = message.MemberList.Member

	displayText(fmt.Sprintf("<CMD_SYNC_REPLY>: reply received from %v:%v - call \\list",
		message.Sender.Ip, message.Sender.Port))
	return nil
}
