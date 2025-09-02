package invites

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-faster/errors"
)

const (
	inviteLinkSchema = "https"
	inviteLinkHost   = "t.me"
	ChannelLinkBase  = "https://t.me/"
	joinchatPrefix   = "joinchat/"
	minHashLinkLen   = 16
)

var (
	ErrInvalidInviteLink      = errors.New("invalid invite link")
	ErrNotEnoughLenInviteLink = errors.Wrap(
		ErrInvalidInviteLink,
		fmt.Sprintf("invite link hash can't be less than %d", minHashLinkLen),
	)

	ErrInvalidInviteLinkMeta   = errors.New("invalid meta of invite link")
	ErrChannelIDIsEmpty        = errors.Wrap(ErrInvalidInviteLinkMeta, "empty channel id")
	ErrInvalidCreatedAt        = errors.Wrap(ErrInvalidInviteLinkMeta, "invalid created at")
	ErrValidToEarlierCreatedAt = errors.Wrap(ErrInvalidInviteLinkMeta, "valid to field should be later than crated at")
)

type (
	InviteLink string
	Hash       string
)

func (h Hash) Validate() error {
	if len(strings.TrimPrefix(h.String(), "+")) < minHashLinkLen {
		return ErrNotEnoughLenInviteLink
	}
	return nil
}

func (h Hash) String() string {
	return string(h)
}

func (l InviteLink) String() string {
	return string(l)
}

func (l InviteLink) Hash() Hash {
	linkParts := strings.Split(l.String(), "/")
	return Hash(linkParts[len(linkParts)-1])
}

func (l InviteLink) Full() string {
	if strings.HasPrefix(l.String(), ChannelLinkBase) {
		return l.String()
	}
	if strings.HasPrefix(l.String(), "+") {
		return ChannelLinkBase + l.String()
	}
	return ChannelLinkBase + joinchatPrefix + l.String()
}

func (l InviteLink) Validate() error {
	if err := l.Hash().Validate(); err != nil {
		return err
	}

	if l.IsHash() {
		return nil
	}

	link, err := url.Parse(l.String())
	if err != nil {
		return ErrInvalidInviteLink
	}

	if link.Scheme != "" && link.Scheme != inviteLinkSchema {
		return errors.Wrap(ErrInvalidInviteLink, "invalid schema")
	}
	if link.Host != inviteLinkHost {
		return errors.Wrap(ErrInvalidInviteLink, "invalid host")
	}

	return nil
}

func (l InviteLink) IsHash() bool {
	return len(strings.Split(l.String(), "/")) == 1
}

type InviteLinkMeta struct {
	ChannelID       int64
	Name            string
	ApproveRequired bool
	CreatedAt       time.Time
	ValidTo         time.Time
	UserLimit       uint32
}

func (m InviteLinkMeta) Validate() error {
	if m.ChannelID == 0 {
		return ErrChannelIDIsEmpty
	}
	if m.CreatedAt.IsZero() {
		return ErrInvalidCreatedAt
	}

	if !m.ValidTo.IsZero() && m.ValidTo.Before(m.CreatedAt) {
		return ErrValidToEarlierCreatedAt
	}
	return nil
}

type ChannelInviteLink struct {
	Link InviteLink
	Meta InviteLinkMeta
}

func (cil ChannelInviteLink) Validate() error {
	if err := cil.Link.Validate(); err != nil {
		return err
	}

	return cil.Meta.Validate()
}

type ChannelInviteLinkUpdateModel struct {
	Link            InviteLink
	ChannelID       int64
	Name            *string
	ApproveRequired *bool
	ValidTo         *time.Time
	UserLimit       *uint32
}

func (c ChannelInviteLinkUpdateModel) Validate() error {
	if err := c.Link.Validate(); err != nil {
		return err
	}

	if c.ChannelID == 0 {
		return errors.New("channel id can't be zero")
	}
	return nil
}
