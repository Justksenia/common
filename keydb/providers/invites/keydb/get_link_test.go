package keydb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/adstail/ts-common/entities/invites"
)

func (s *InviteLinkProviderTestSuite) TestGetLink() {
	date := time.Now().UTC()

	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		hash := "+fRwO0LahW6JkODYy"
		link := &invites.ChannelInviteLink{
			Link: "https://t.me/+fRwO0LahW6JkODYy",
			Meta: invites.InviteLinkMeta{
				ChannelID:       1,
				Name:            "name",
				ApproveRequired: true,
				CreatedAt:       date,
			},
		}

		require.NoError(t, s.instance.Set(ctx, hashLinkKey(hash), link.Meta))
		require.NoError(t, s.instance.RPush(ctx, channelIDKey(link.Meta.ChannelID), hash))

		actual, err := s.adapter.GetLink(ctx, "https://t.me/+fRwO0LahW6JkODYy")
		assert.NoError(t, err)
		assert.Equal(t, link, actual)
	})

	t.Run("not found", func(t *testing.T) {
		link := invites.InviteLink("https://t.me/+DwAO0LahW6JkODR8")
		actual, err := s.adapter.GetLink(ctx, link)
		assert.ErrorIs(t, err, ErrLinkNotFound)
		assert.Nil(t, actual)
	})
}

func (s *InviteLinkProviderTestSuite) TestGetChannelLinks() {
	date := time.Now().UTC()

	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		links := []invites.InviteLink{"https://t.me/+N855559pWTI2ZjEy", "https://t.me/+N866669pWTI2ZjEy"}
		for _, link := range links {
			l := &invites.ChannelInviteLink{
				Meta: invites.InviteLinkMeta{
					ChannelID: 999,
					CreatedAt: date,
				},
			}
			require.NoError(t, s.instance.Set(ctx, hashLinkKey(link.Hash().String()), l.Meta))
			require.NoError(t, s.instance.RPush(ctx, channelIDKey(l.Meta.ChannelID), link.Hash().String()))
		}

		actual, err := s.adapter.GetChannelInviteLinks(ctx, 999)
		assert.NoError(t, err)
		assert.Equal(t, links, actual)
	})

	t.Run("no links", func(t *testing.T) {
		_, err := s.adapter.GetChannelInviteLinks(ctx, 10000)
		assert.ErrorIs(t, err, ErrLinkNotFound)
	})
}

func (s *InviteLinkProviderTestSuite) TestGetLastLinkChannel() {
	date := time.Now().UTC()

	var (
		t   = s.T()
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		link := &invites.ChannelInviteLink{
			Link: "https://t.me/+N8iY-39pWTI2ZjEy",
			Meta: invites.InviteLinkMeta{
				ChannelID: 1,
				CreatedAt: date,
			},
		}
		require.NoError(t, s.instance.Set(ctx, hashLinkKey("+N8iY-39pWTI2ZjEy"), link.Meta))
		require.NoError(t, s.instance.RPush(ctx, channelIDKey(link.Meta.ChannelID), "+N8iY-39pWTI2ZjEy"))

		actual, err := s.adapter.GetLastLinkChannel(ctx, link.Meta.ChannelID)
		assert.NoError(t, err)
		assert.Equal(t, link, actual)
	})

	t.Run("no elements in list", func(t *testing.T) {
		actual, err := s.adapter.GetLastLinkChannel(ctx, 1000)
		assert.ErrorIs(t, err, ErrLinkNotFound)
		assert.Nil(t, actual)
	})

	t.Run("there is link, but there is no in list", func(t *testing.T) {
		link := invites.ChannelInviteLink{
			Link: "https://t.me/+DDDO0LahW6JkODYy",
			Meta: invites.InviteLinkMeta{
				ChannelID:       1,
				Name:            "name",
				ApproveRequired: true,
				CreatedAt:       date,
			},
		}
		require.NoError(t, s.instance.RPush(ctx, channelIDKey(link.Meta.ChannelID), link.Link))
		actual, err := s.adapter.GetLastLinkChannel(ctx, link.Meta.ChannelID)
		assert.ErrorIs(t, err, ErrLinkNotFound)
		assert.Nil(t, actual)
	})
}
