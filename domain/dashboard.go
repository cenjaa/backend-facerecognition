package domain

import "time"

type DashboardStats struct {
	Menunggu    int    `json:"menunggu"`
	Hadir       int    `json:"hadir"`
	TidakHadir  int    `json:"tidak_hadir"`
	Telat       int    `json:"telat"`
	Keluar      int    `json:"keluar"`
	Cuti        int    `json:"cuti"`
	TotalHadir  int    `json:"total_hadir"`
	UpdatedAt   string `json:"updated_at"`
}

type HourlyFrequency struct {
	Labels []string `json:"labels"`
	Data   []int    `json:"data"`
}

type RecentPresence struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Time string `json:"time"`
}

type AttendanceHistoryRow struct {
	No         int    `json:"no"`
	Name       string `json:"name"`
	NIP        string `json:"nip"`
	Status     string `json:"status"`
	Confidence int    `json:"confidence"`
	Latency    string `json:"latency"`
	Time       string `json:"time"`
}

type JiraHistoryRow struct {
	ID     int    `json:"id"`
	Ticket string `json:"ticket"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Date   string `json:"date"`
}

type KpiAccumulationRow struct {
	No         int    `json:"no"`
	Name       string `json:"name"`
	TotalTiket int    `json:"total_tiket"`
}

type PaginatedResult[T any] struct {
	Records    []T `json:"records"`
	Total      int `json:"total"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
}

type BulkAttendanceRecord struct {
	UserID       int     `json:"user_id"`
	PresenceTime string  `json:"presence_time"`
	Confidence   float64 `json:"confidence"`
	LatencyMs    float64 `json:"latency_ms"`
	StatusID     int     `json:"status_id"`
}

type UserListRow struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	NIP    string `json:"nip"`
	Active bool   `json:"active"`
}

type TrainStatus struct {
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	Message  string `json:"message"`
}

func Today() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}
