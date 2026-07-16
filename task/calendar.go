package task

import (
	"context"
	"encoding/json"
	"server/global"
	"server/utils"
	"time"
)

func GetCalendarSyncTask(ctx context.Context) error {
	dateStr := time.Now().Format("2006/0102")
	calendar, err := utils.GetCalendar(dateStr)
	if err != nil {
		return err
	}

	data, err := json.Marshal(calendar)
	if err != nil {
		return err
	}
	if err := global.Redis.Set(ctx, "calendar-"+dateStr, data, time.Hour*24).Err(); err != nil {
		return err
	}
	return nil
}
