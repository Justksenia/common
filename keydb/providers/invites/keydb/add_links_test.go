package keydb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/adstail/ts-common/entities/invites"
)

func (s *InviteLinkProviderTestSuite) TestAddLink() {
	date := time.Now().UTC()

	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		input := invites.ChannelInviteLink{
			Link: "https://t.me/+5V23yMex8GY5ZWFi",
			Meta: invites.InviteLinkMeta{
				ChannelID:       1,
				Name:            "name",
				ApproveRequired: false,
				CreatedAt:       date,
			},
		}

		err := s.adapter.AddLink(ctx, input)
		assert.NoError(t, err)
		s.compareResults(t, input)
	})

	t.Run("already exists", func(t *testing.T) {
		input := invites.ChannelInviteLink{
			Link: "https://t.me/+6hlhkJIshkxmZjIy",
			Meta: invites.InviteLinkMeta{
				ChannelID:       2,
				Name:            "name",
				ApproveRequired: false,
				CreatedAt:       date,
			},
		}

		require.NoError(t, s.adapter.AddLink(ctx, input))
		assert.ErrorIs(t, s.adapter.AddLink(ctx, input), ErrLinkAlreadyExists)
	})

	t.Run("invalid link", func(t *testing.T) {
		input := invites.ChannelInviteLink{
			Link: "https://t.me/+6hlhkJIshkxmZj",
			Meta: invites.InviteLinkMeta{
				ChannelID:       2,
				Name:            "name",
				ApproveRequired: false,
				CreatedAt:       date,
			},
		}
		require.ErrorIs(t, s.adapter.AddLink(ctx, input), invites.ErrInvalidInviteLink)
	})
}

func (s *InviteLinkProviderTestSuite) compareResults(t *testing.T, invite invites.ChannelInviteLink) {
	t.Helper()
	var ctx = context.Background()

	hash := invite.Link.Hash().String()
	var inviteLinks []string
	require.NoError(t, s.instance.GetList(ctx, channelIDKey(invite.Meta.ChannelID), &inviteLinks))

	var linkInList bool
	for _, i := range inviteLinks {
		if i == hash {
			linkInList = true
		}
	}

	if !linkInList {
		t.Errorf("there is no link in redis list")
		return
	}

	var meta invites.InviteLinkMeta
	require.NoError(t, s.instance.Get(ctx, hashLinkKey(hash), &meta))

	assert.Equal(t, invite.Meta, meta)
}
