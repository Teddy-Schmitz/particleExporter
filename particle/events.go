package particle

import "time"

type Event struct {
	Name        string    `json:"event"`
	Data        string    `json:"data"`
	DeviceID    string    `json:"coreid"`
	PublishedAt time.Time `json:"published_at"`
}
