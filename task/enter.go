package task

import (
	"context"
	"server/global"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func RegisterScheduledTasks(c *cron.Cron) error {
	//更新浏览量定时任务
	if _, err := c.AddFunc("@hourly", func() {
		if err := UpdateArticleViewsSyncTask(context.Background()); err != nil {
			global.Log.Error("Failed to update article views:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	//更新日历信息定时任务
	if _, err := c.AddFunc("@daily", func() {
		if err := GetCalendarSyncTask(context.Background()); err != nil {
			global.Log.Error("Failed to get calendar:", zap.Error(err))
		}
	}); err != nil {
		return err
	}
	return nil
}
