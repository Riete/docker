package reader

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/riete/convert/str"

	set "github.com/riete/go-set"
)

type Message interface {
	Message() string
	ID() string
	Full() bool
}

type ImagePullPushMessage struct {
	Status      string `json:"status"`
	Progress    string `json:"progress"`
	Id          string `json:"id"`
	ErrorDetail struct {
		Message string `json:"message"`
	}
}

func (i ImagePullPushMessage) Message() string {
	if i.ErrorDetail.Message != "" {
		return i.ErrorDetail.Message
	}
	if i.Id != "" {
		return fmt.Sprintf("%s: %s %s", i.Id, i.Status, i.Progress)
	}
	return fmt.Sprintf("%s %s", i.Status, i.Progress)
}

func (i ImagePullPushMessage) ID() string {
	return i.Id
}

func (i ImagePullPushMessage) Full() bool {
	return i.Id != ""
}

func NewImagePullPushMessage(s string) Message {
	m := &ImagePullPushMessage{}
	_ = json.Unmarshal(str.ToBytes(s), m)
	return m
}

type ImageBuildMessage struct {
	ImagePullPushMessage
	Stream string `json:"stream"`
	Aux    struct {
		ID string `json:"ID"`
	} `json:"aux"`
}

func (i ImageBuildMessage) Message() string {
	if i.ErrorDetail.Message != "" {
		return i.ErrorDetail.Message
	}
	if i.Stream != "" {
		return i.Stream
	}
	if i.Aux.ID != "" {
		return i.Aux.ID
	}
	if i.Id != "" {
		return fmt.Sprintf("%s: %s %s", i.Id, i.Status, i.Progress)
	}
	return fmt.Sprintf("%s %s", i.Status, i.Progress)
}

func (i ImageBuildMessage) ID() string {
	return i.Id
}

func (i ImageBuildMessage) Full() bool {
	return i.Id != ""
}

func NewImageBuildMessage(s string) Message {
	m := &ImageBuildMessage{}
	_ = json.Unmarshal(str.ToBytes(s), m)
	return m
}

type MessageParser struct {
	message map[string]string
	id      []string
	s       set.Set
}

func (p *MessageParser) Message(m Message) string {
	if !m.Full() {
		return m.Message()
	}
	if !p.s.Has(m.ID()) {
		p.id = append(p.id, m.ID())
		p.s.Add(m.ID())
	}
	p.message[m.ID()] = m.Message()
	var msg []string
	for _, i := range p.id {
		msg = append(msg, p.message[i])
	}
	return strings.Join(msg, "\n")
}

func (p MessageParser) Full(m Message) bool {
	return m.Full()
}

func NewMessageParser() *MessageParser {
	return &MessageParser{message: make(map[string]string), s: set.NewSet()}
}
