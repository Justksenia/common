package keydb

import (
	"context"

	"github.com/go-faster/errors"
	"gitlab.com/adstail/ts-common/entities/invites"
	"gitlab.com/adstail/ts-common/keydb/redis"
	cmnlogger "gitlab.com/adstail/ts-common/logger"
	"gitlab.com/adstail/ts-common/tracer"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (p *InviteLinksKeyDBProvider) GetLink(ctx context.Context, link invites.InviteLink) (*invites.ChannelInviteLink, error) {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	hash := link.Hash().String()

	var invite invites.InviteLinkMeta
	if err := p.client.Get(ctx, hashLinkKey(hash), &invite); err != nil {
		if errors.Is(err, redis.ErrNoData) {
			return nil, ErrLinkNotFound
		}
		return nil, span.Error(errors.Wrap(err, "get link"))
	}

	return &invites.ChannelInviteLink{
		Link: link,
		Meta: invite,
	}, nil
}

func (p *InviteLinksKeyDBProvider) GetChannelInviteLinks(ctx context.Context, channelID int64) ([]invites.InviteLink, error) {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	var links []invites.InviteLink
	if err := p.client.GetList(ctx, channelIDKey(channelID), &links); err != nil && !errors.Is(err, redis.ErrNoData) {
		return nil, errors.Wrap(err, "get list")
	}

	if len(links) == 0 {
		return nil, ErrLinkNotFound
	}

	fullLinks := make([]invites.InviteLink, len(links))
	for i, link := range links {
		fullLinks[i] = invites.InviteLink(link.Full())
	}

	return fullLinks, nil
}

// GetLastLinkChannel - optimistic way, just take the last link in list.
func (p *InviteLinksKeyDBProvider) GetLastLinkChannel(ctx context.Context, channelID int64) (*invites.ChannelInviteLink, error) {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	logger := cmnlogger.FromContext(ctx)

	var link invites.InviteLink
	if err := p.client.GetElementByPosition(ctx, channelIDKey(channelID), redis.ListElementLastPosition, &link); err != nil {
		if errors.Is(err, redis.ErrNoData) {
			return nil, ErrLinkNotFound
		}
		return nil, span.Error(errors.Wrap(err, "get last channel link"))
	}

	if link == "" {
		return nil, ErrLinkNotFound
	}

	var meta invites.InviteLinkMeta
	if err := p.client.Get(ctx, hashLinkKey(link.String()), &meta); err != nil {
		if errors.Is(err, redis.ErrNoData) {
			logger.Warn(
				"there is link in channel list but there is no link meta",
				zap.String("link", link.Full()),
				zap.Int64("channel_id", channelID),
			)
			return nil, ErrLinkNotFound
		}
		return nil, span.Error(errors.Wrap(err, "get link's meta information"))
	}

	return &invites.ChannelInviteLink{
		Link: invites.InviteLink(link.Full()),
		Meta: meta,
	}, nil
}
