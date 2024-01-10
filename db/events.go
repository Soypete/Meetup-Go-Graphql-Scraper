package db

import (
	"fmt"
	"time"
)

func (w Warehouse) InsertEvent(id, title, groupID, groupName string, date time.Time, going, waiting int) error {
	query := `
	INSERT OR REPLACE INTO events (id, title, date, going, waiting, group_id, group_name)
	VALUES (?,?,?,?,?,?,?);
	`
	_, err := w.db.Exec(query, id, title, date, going, waiting, groupID, groupName)
	if err != nil {
		return fmt.Errorf("error inserting event, %w", err)
	}
	return nil
}
