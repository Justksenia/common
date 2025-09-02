//nolint:gomnd,gochecknoglobals // default values
package retrier

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/samber/lo"
	cmnlogger "gitlab.com/adstail/ts-common/logger"
	"go.uber.org/zap"
)

type RetryPolicy struct {
	MaxAttempts        int
	StartDelay         time.Duration
	MaxDelay           *time.Duration
	BackoffCoefficient float32
}

type Retrier struct {
	policy         RetryPolicy
	excludedErrors []error
}

type Opts func(r *Retrier)

func WithRetryPolicy(rp RetryPolicy) Opts {
	return func(r *Retrier) {
		r.policy = rp
	}
}

func WithExcludedErrors(errors ...error) Opts {
	return func(r *Retrier) {
		r.excludedErrors = errors
	}
}

var defaultPolicy = RetryPolicy{
	MaxAttempts:        3,
	StartDelay:         1 * time.Second,
	MaxDelay:           lo.ToPtr(10 * time.Second),
	BackoffCoefficient: 2,
}

func NewRetrier(opts ...Opts) *Retrier {
	retrier := &Retrier{policy: defaultPolicy}

	for _, opt := range opts {
		opt(retrier)
	}
	return retrier
}

func (r *Retrier) Wrap(ctx context.Context, name string, f func() error) (err error) {
	logger := cmnlogger.FromContext(ctx).With(zap.String("method", "retrier"), zap.String("name", name))

	delay := r.policy.StartDelay
	for i := 1; i <= r.policy.MaxAttempts; i++ {
		logger.Debug("start execution", zap.Int("attempt", i))

		if err = f(); err == nil || r.checkExcludedErrors(err) {
			break
		}
		logger.Warn("error occurred during execution", zap.Error(err))

		if i != r.policy.MaxAttempts {
			time.Sleep(delay)
			delay = time.Duration(float32(delay) * r.policy.BackoffCoefficient)
			if r.policy.MaxDelay != nil && delay > *r.policy.MaxDelay {
				delay = *r.policy.MaxDelay
			}
		}
	}
	if err == nil {
		logger.Debug("execution finished")
	}
	return err
}

func (r *Retrier) checkExcludedErrors(err error) bool {
	_, ok := lo.Find(r.excludedErrors, func(item error) bool {
		return errors.Is(err, item)
	})
	return ok
}
