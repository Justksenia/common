package keydb

import (
	"context"

	"github.com/go-faster/errors"
	"gitlab.com/adstail/ts-common/entities/invites"
	"gitlab.com/adstail/ts-common/tracer"
	"go.opentelemetry.io/otel/trace"
)

func (p *InviteLinksKeyDBProvider) RemoveLink(ctx context.Context, channelID int64, link invites.InviteLink) error {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	if err := link.Validate(); err != nil {
		return span.Error(errors.Wrap(err, "validation"))
	}

	tx, err := p.client.Begin(ctx)
	if err != nil {
		return span.Error(errors.Wrap(err, "begin transaction"))
	}

	hash := link.Hash().String()

	if err = tx.Delete(ctx, hashLinkKey(hash)); err != nil {
		_ = tx.Rollback(ctx)
		return span.Error(errors.Wrap(err, "delete link meta"))
	}

	if err = tx.RemoveFromList(ctx, channelIDKey(channelID), hash); err != nil {
		_ = tx.Rollback(ctx)
		return span.Error(errors.Wrap(err, "remove from list"))
	}

	if err = tx.Commit(ctx); err != nil {
		return span.Error(errors.Wrap(err, "commit transaction"))
	}

	return nil
}
