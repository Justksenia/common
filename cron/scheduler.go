package cron

import (
	"fmt"
	"strings"

	"github.com/go-faster/errors"
	"github.com/robfig/cron/v3"
)

type Job interface {
	Run()
}

type Scheduler interface {
	AddJob(j Job, spec string) error

	Start()
	Stop()
}

type scheduler struct {
	crons *cron.Cron
}

func NewCronScheduler() Scheduler {
	crons := cron.New()
	return &scheduler{crons: crons}
}

func (c *scheduler) AddJob(j Job, spec string) error {
	splitted := strings.Split(spec, " ")
	option := (cron.Minute | cron.Hour | cron.Dom | cron.Month) & (cron.Minute<<(len(splitted)) - 1)
	if strings.HasPrefix(spec, "@") {
		option = cron.Descriptor
	}

	sch, err := cron.NewParser(option).Parse(spec)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("parse spec %s", spec))
	}

	c.crons.Schedule(sch, j)
	return nil
}

func (c *scheduler) Start() {
	go c.crons.Run()
}

func (c *scheduler) Stop() {
	c.crons.Stop()
}
