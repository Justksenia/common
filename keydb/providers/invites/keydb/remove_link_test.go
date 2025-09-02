package keydb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/Justksenia/common/entities/invites"
	"github.com/Justksenia/common/keydb/redis"
)

func (s *InviteLinkProviderTestSuite) TestRemoveLink() {
	var (
		t    = s.T()
		ctx  = context.Background()
		date = time.Now().UTC()
	)

	t.Run("success", func(t *testing.T) {
		link := invites.ChannelInviteLink{
			Link: "https://t.me/+AAAAAAAAAAAAAAAA",
			Meta: invites.InviteLinkMeta{
				ChannelID: 10,
				CreatedAt: date,
			},
		}

		require.NoError(t, s.instance.RPush(ctx, channelIDKey(link.Meta.ChannelID), link.Link.Hash().String()))
		require.NoError(t, s.instance.Set(ctx, hashLinkKey(link.Link.Hash().String()), link.Meta))

		assert.NoError(t, s.adapter.RemoveLink(ctx, link.Meta.ChannelID, link.Link))

		var actualLink string
		assert.ErrorIs(t, s.instance.Get(ctx, hashLinkKey(link.Link.Hash().String()), &actualLink), redis.ErrNoData)

		var actualMeta []invites.InviteLinkMeta
		assert.ErrorIs(t, s.instance.GetList(ctx, channelIDKey(link.Meta.ChannelID), &actualMeta), redis.ErrNoData)
	})

	t.Run("remove not existed link", func(t *testing.T) {
		assert.NoError(t, s.adapter.RemoveLink(ctx, 1000, "https://t.me/+BBBBBBBBBBBBBBBB"))
	})
}
