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

func (p *InviteLinksKeyDBProvider) AddLink(ctx context.Context, inviteLink invites.ChannelInviteLink) error {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	if err := inviteLink.Validate(); err != nil {
		return span.Error(errors.Wrap(err, "validation"))
	}

	tx, err := p.client.Begin(ctx)
	if err != nil {
		return span.Error(errors.Wrap(err, "begin transaction"))
	}

	if err = p.addLink(ctx, inviteLink, tx); err != nil {
		_ = tx.Rollback(ctx)
		return span.Error(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return span.Error(errors.Wrap(err, "commit transaction"))
	}
	return nil
}

func (p *InviteLinksKeyDBProvider) AddLinks(ctx context.Context, inviteLinks []invites.ChannelInviteLink) error {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	logger := cmnlogger.FromContext(ctx)

	tx, err := p.client.Begin(ctx)
	if err != nil {
		return span.Error(errors.Wrap(err, "begin transaction"))
	}

	for _, link := range inviteLinks {
		if err = link.Validate(); err != nil {
			logger.Error("AddLinks", zap.String("link", link.Link.String()), zap.Error(err))
			continue
		}

		if err = p.addLink(ctx, link, tx); err != nil && !errors.Is(err, ErrLinkAlreadyExists) {
			_ = tx.Rollback(ctx)
			return span.Error(err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return span.Error(errors.Wrap(err, "commit transaction"))
	}
	return nil
}

func (p *InviteLinksKeyDBProvider) addLink(ctx context.Context, link invites.ChannelInviteLink, tx *redis.Instance) error {
	hash := link.Link.Hash().String()
	exists, err := p.client.IsExist(ctx, hashLinkKey(hash))
	if err != nil {
		return errors.Wrap(err, "check existence")
	}

	if exists {
		return ErrLinkAlreadyExists
	}

	cl := p.client
	if tx != nil {
		cl = tx
	}

	if err = cl.RPush(ctx, channelIDKey(link.Meta.ChannelID), hash); err != nil {
		return errors.Wrap(err, "add link for the channel")
	}

	if err = cl.Set(ctx, hashLinkKey(hash), link.Meta); err != nil {
		return errors.Wrap(err, "save link")
	}
	return nil
}
