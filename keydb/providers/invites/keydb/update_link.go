package keydb

import (
	"context"

	"github.com/go-faster/errors"
	"gitlab.com/adstail/ts-common/entities/invites"
	"gitlab.com/adstail/ts-common/keydb/redis"
	"gitlab.com/adstail/ts-common/tracer"
	"go.opentelemetry.io/otel/trace"
)

func (p *InviteLinksKeyDBProvider) UpdateLink(ctx context.Context, inviteLink invites.ChannelInviteLinkUpdateModel) error {
	ctx, span := tracer.StartSpan(ctx, tracer.AutoFillName(), trace.SpanKindInternal)
	defer span.End()

	if err := inviteLink.Validate(); err != nil {
		return errors.Wrap(err, "validation")
	}

	hash := inviteLink.Link.Hash().String()

	var actualInviteLink invites.InviteLinkMeta

	if err := p.client.Get(ctx, hashLinkKey(hash), &actualInviteLink); err != nil {
		if errors.Is(err, redis.ErrNoData) {
			return ErrLinkNotFound
		}
		return span.Error(errors.Wrap(err, "get actual version of link"))
	}

	if inviteLink.Name != nil {
		actualInviteLink.Name = *inviteLink.Name
	}

	if inviteLink.ApproveRequired != nil {
		actualInviteLink.ApproveRequired = *inviteLink.ApproveRequired
	}

	if inviteLink.ValidTo != nil {
		actualInviteLink.ValidTo = *inviteLink.ValidTo
	}

	if inviteLink.UserLimit != nil {
		actualInviteLink.UserLimit = *inviteLink.UserLimit
	}

	if err := p.client.Set(ctx, hashLinkKey(hash), actualInviteLink); err != nil {
		return span.Error(errors.Wrap(err, "update link"))
	}

	return nil
}
