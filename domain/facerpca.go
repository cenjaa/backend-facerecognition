package domain

import (
	"context"
	"mime/multipart"
	"time"
)

type AdminUser struct {
	IDAdmin  int64  `json:"id_admin"`
	IDUser   int    `json:"id_user"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	PinCode  string `json:"pin_code"`
	Active   bool   `json:"active"`
}

type User struct {
	IDUser       int    `json:"id_user"`
	Name         string `json:"name"`
	NIP          string `json:"nip"`
	Email        string `json:"email"`
	Active       bool   `json:"active"`
	CreatedBy    string `json:"created_by"`
	DatasetPath  string `json:"dataset_path"`
	AkumulasiKPI int    `json:"akumulasi_kpi"`
}

type AttendanceRecord struct {
	IDUser       int			`json:"id_user"`
	PresenceTime time.Time		`json:"precence_time"`
	Confidence   Confidence		`json:"confidence"`
	LatencyMS    float64		`json:"latency_ms"`
	Status       int			`json:"status_id"`
}

type KPIHarian struct {
	IDTiketJira   string    `json:"id_tiket_jira"`
	StoryJira     string    `json:"story_jira"`
	IDStatusTiket int       `json:"id_status_tiket"`
	NamaTiket     string    `json:"nama_tiket"`
	IDUser        int       `json:"id_user"`
	Date          time.Time `json:"date"`
}

type JiraTask struct {
	TicketKey string
	StoryJira string
	Status    int
}

type FaceRPCASQLRepository interface {
	IsAdmin(ctx context.Context, userID int) (bool, error)
	VerifyAdminPin(ctx context.Context, userID int, pin string) (bool, error)
	GetByID(ctx context.Context, userID int) (*User, error)
	CreateUser(ctx context.Context, name, nip, email, createdBy, datasetPath string) (int, error)
	GetUserNameByID(ctx context.Context, userID int) (string, error)
	GetAllActiveUsers(ctx context.Context) ([]User, error)
	GetUserEmailByID(ctx context.Context, userID int) (string, error)
	GetAll(ctx context.Context) ([]UserListRow, error)
	Update(ctx context.Context, userID int, name, nip string, active bool) error
	Delete(ctx context.Context, userID int) error

	SaveAttendance(ctx context.Context, attendance AttendanceRecord) error
	GetLoggedInUsers(ctx context.Context, today time.Time) ([]int, error)
	GetUserStatusToday(ctx context.Context, userID int, today time.Time) (*int, error)
	MarkAbsentUsers(ctx context.Context, today time.Time, cutoffTime time.Time) error
	ForceLogoffUsers(ctx context.Context, today time.Time) ([]int, error)
	InitializeDailyStatus(ctx context.Context, today time.Time) error
	SetUserCuti(ctx context.Context, userID int, today time.Time) error
	GetUsersWhoWorkedOnDate(ctx context.Context, date time.Time) ([]int, error)

	UpsertKpiHarian(ctx context.Context, kpi KPIHarian) error
	GetDoneTasksCount(ctx context.Context, userID int, today time.Time) (int, error) 
	HasDoneTasks(ctx context.Context, userID int, today time.Time) (bool, error)
	UpdateAkumulasiKPI(ctx context.Context, userID int, akumulasiKPI int) error
	GetRecentJira(ctx context.Context, limit int) ([]JiraHistoryRow, error)
	GetJiraHistory(ctx context.Context, page, perPage int) (*PaginatedResult[JiraHistoryRow], error)
	GetKpiAccumulation(ctx context.Context, page, perPage int) (*PaginatedResult[KpiAccumulationRow], error)

	SaveRaw(ctx context.Context, rec BulkAttendanceRecord) error
	SaveBulk(ctx context.Context, records []BulkAttendanceRecord) (int, error)
	GetRecentPresence(ctx context.Context, limit int) ([]RecentPresence, error)
	GetHistory(ctx context.Context, page, perPage int) (*PaginatedResult[AttendanceHistoryRow], error)
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetHourlyFrequency(ctx context.Context) (*HourlyFrequency, error)
	GetTotalUserCount(ctx context.Context) (int, error)
	GetTodayAttendeeCount(ctx context.Context) (int, error)
}

type FaceRPCAMinIORepository interface {
	DownloadModels(ctx context.Context, localDir string) error
	UploadModel(ctx context.Context, localDir string) error
	UploadDataset(ctx context.Context, datasetDir string) error
	DownloadDataset(ctx context.Context, localDir string) error
	DeleteUserFaces(ctx context.Context, userID int) error
	UploadUserFaces(ctx context.Context, userID int, localUserDir string) (int, error)
	GetModelTimestamp(ctx context.Context) (float64, error)
	GetDatasetTimestamp(ctx context.Context) (float64, error)
}

type FaceRPCAJiraRepository interface {
	GetTasksForUser(ctx context.Context, email string) ([]JiraTask, error)
}

type FaceRPCAUsecase interface {
	CreateUser(ctx context.Context, name, nip, email, createdBy, datasetPath string) (int, error)
	UpdateUser(ctx context.Context, userID int, name, nip string, active bool) error
	DeleteUser(ctx context.Context, userID int) error
	VerifyAdminPin(ctx context.Context, userID int, pin string) (bool, string, error)
	RecordAttendance(ctx context.Context, userID int, mode string, confidence, latency float64) error
	LogAttendanceRaw(ctx context.Context, rec BulkAttendanceRecord) error
	LogAttendanceBulk(ctx context.Context, records []BulkAttendanceRecord) (int, error)
	UploadDataset(ctx context.Context, userID int, fileHeaders []*multipart.FileHeader) error
	TrainModel(ctx context.Context) error
	GetTrainStatus(ctx context.Context) (*TrainStatus, error)
	HandleDailyInit(ctx context.Context) error
	HandleMorningCutoff(ctx context.Context) error
	HandleEveningCutoff(ctx context.Context) error
	PollJiraTasks(ctx context.Context) error
	PollJiraTasksForDate(ctx context.Context, date time.Time) error
	SetUserCuti(ctx context.Context, userID int) error
	AccumulateKPI(ctx context.Context, userID int, today time.Time) error
	GetAllUsers(ctx context.Context) ([]UserListRow, error)
	GetUserName(ctx context.Context, userID int) (string, error)
	GetUserStatusToday(ctx context.Context, userID int) (*int, error)
	GetRecentPresence(ctx context.Context) ([]RecentPresence, error)
	GetAttendanceHistory(ctx context.Context, page, perPage int) (*PaginatedResult[AttendanceHistoryRow], error)
	GetDashboardStats(ctx context.Context) (*DashboardStats, error)
	GetAttendanceFrequency(ctx context.Context) (*HourlyFrequency, error)
	GetRecentJira(ctx context.Context) ([]JiraHistoryRow, error)
	GetJiraHistory(ctx context.Context, page, perPage int) (*PaginatedResult[JiraHistoryRow], error)
	GetKpiAccumulation(ctx context.Context, page, perPage int) (*PaginatedResult[KpiAccumulationRow], error)
	GetAdminDashboard(ctx context.Context, adminID int) (*AdminDashboard, error)
}

type AdminDashboard struct {
	AdminName      string  `json:"admin_name"`
	TotalUsers     int     `json:"total_users"`
	TodayAttendees int     `json:"today_attendees"`
	NeedsRetrain   bool    `json:"needs_retrain"`
	ModelTimestamp float64 `json:"model_timestamp"`
}