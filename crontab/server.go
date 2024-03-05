package crontab

import (
	"context"

	"github.com/lw396/WeComCopilot/service"
	"github.com/robfig/cron/v3"
)

type crontabServer struct {
	cron    *cron.Cron
	service *service.Service
}

func NewServer(s *service.Service) *crontabServer {
	return &crontabServer{
		service: s,
		cron: cron.New(
			cron.WithSeconds(),
			cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger),
			)),
	}
}

func (s *crontabServer) Start(ctx context.Context) error {
	if err := s.InitSyncTask(ctx); err != nil {
		return err
	}

	if _, err := s.cron.AddFunc("*/10 * * * * *", func() {
		if err := s.SyncMessage(context.Background()); err != nil {
			return
		}
	}); err != nil {
		return err
	}

	s.cron.Start()
	<-ctx.Done()
	return nil
}

func (s *crontabServer) Stop() {
	s.cron.Stop()
}
