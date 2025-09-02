package keydb

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/adstail/ts-common/entities/invites"
)

func (s *InviteLinkProviderTestSuite) TestUpdateLink() {
	var (
		t    = s.T()
		ctx  = context.Background()
		date = time.Now().UTC()
	)

	t.Run("success", func(t *testing.T) {
		link := invites.ChannelInviteLink{
			Link: "https://t.me/+AAAAAAAAAAAAAANN",
			Meta: invites.InviteLinkMeta{
				ChannelID: 300,
				CreatedAt: date,
			},
		}
		hash := link.Link.Hash().String()

		require.NoError(t, s.instance.Set(ctx, hashLinkKey(hash), link.Meta))
		require.NoError(t, s.instance.RPush(ctx, channelIDKey(link.Meta.ChannelID), hash))

		linkUp := invites.ChannelInviteLinkUpdateModel{
			Link:      "https://t.me/+AAAAAAAAAAAAAANN",
			ChannelID: 300,
			ValidTo:   lo.ToPtr(date.Add(48 * time.Hour)),
			UserLimit: nil,
		}

		assert.NoError(t, s.adapter.UpdateLink(ctx, linkUp))

		var actualValue invites.InviteLinkMeta
		assert.NoError(t, s.instance.Get(ctx, hashLinkKey(hash), &actualValue))

		expectedLink := link
		expectedLink.Meta.ValidTo = *linkUp.ValidTo
		assert.Equal(t, expectedLink.Meta, actualValue)
	})

	t.Run("not existed link", func(t *testing.T) {
		linkUp := invites.ChannelInviteLinkUpdateModel{
			Link:      "https://t.me/+AAAAAAAAAAAADDDD",
			ChannelID: 300,
		}
		assert.ErrorIs(t, s.adapter.UpdateLink(ctx, linkUp), ErrLinkNotFound)
	})

	t.Run("invalid link", func(t *testing.T) {
		linkUp := invites.ChannelInviteLinkUpdateModel{
			Link:      "https://t.me/+AAAAAAAAAAAA",
			ChannelID: 300,
		}
		assert.ErrorIs(t, s.adapter.UpdateLink(ctx, linkUp), invites.ErrInvalidInviteLink)
	})
}
