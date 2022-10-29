package event

import (
	"github.com/go-sonic/sonic/consts"
	"github.com/go-sonic/sonic/model/entity"
)

type Event interface {
	EventType() string
}

const (
	LogEventName              = "LogEvent"
	StartEventName            = "StartEvent"
	UserUpdateEventName       = "UserUpdateEvent"
	ThemeUpdateEventName      = "ThemeUpdateEvent"
	OptionUpdateEventName     = "OptionUpdateEvent"
	ThemeActivatedEventName   = "ThemeActivatedEvent"
	ThemeFileUpdatedEventName = "ThemeFileUpdatedEvent"
	PostUpdateEventName       = "PostUpdateEvent"
	CommentNewEventName       = "CommentNewEvent"
	CommentReplyEventName     = "CommentReplayEvent"
)

type LogEvent struct {
	LogKey    string
	LogType   consts.LogType
	Content   string
	IpAddress string
}

func (*LogEvent) EventType() string {
	return LogEventName
}

type StartEvent struct{}

func (*StartEvent) EventType() string {
	return StartEventName
}

type UserUpdateEvent struct {
	UserID int32
}

func (*UserUpdateEvent) EventType() string {
	return UserUpdateEventName
}

type ThemeUpdateEvent struct{}

func (*ThemeUpdateEvent) EventType() string {
	return ThemeUpdateEventName
}

type OptionUpdateEvent struct{}

func (o *OptionUpdateEvent) EventType() string {
	return OptionUpdateEventName
}

type ThemeActivatedEvent struct{}

func (t *ThemeActivatedEvent) EventType() string {
	return ThemeActivatedEventName
}

type ThemeFileUpdatedEvent struct{}

func (t *ThemeFileUpdatedEvent) EventType() string {
	return ThemeFileUpdatedEventName
}

type PostUpdateEvent struct {
	PostID int32
}

func (p *PostUpdateEvent) EventType() string {
	return PostUpdateEventName
}

type CommentNewEvent struct {
	Comment *entity.Comment
}

func (c *CommentNewEvent) EventType() string {
	return CommentNewEventName
}

type CommentReplyEvent struct {
	Comment *entity.Comment
}

func (c *CommentReplyEvent) EventType() string {
	return CommentReplyEventName
}
