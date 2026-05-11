package sql

import (
	"context"
	"database/sql"
	"time"
	"math"
	"strconv"
	"fmt"
	"face-recognition-fyp/domain"
	"face-recognition-fyp/constant"
)

type FaceRPCASQLRepository struct {
	DB *sql.DB
}

func NewFaceRPCASQLRepository(db *sql.DB) *FaceRPCASQLRepository {
	return &FaceRPCASQLRepository{
		DB: db,
	}
}

//Admin Functions
func (r *FaceRPCASQLRepository) IsAdmin(ctx context.Context, userID int) (bool, error) {
	query := `SELECT count(*) FROM ms_admin WHERE id_user = $1 AND active = TRUE`

	var count int
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FaceRPCASQLRepository) VerifyAdminPin(ctx context.Context, userID int, pin string) (bool, error) {
	query := `SELECT count(*) FROM ms_admin 
	          WHERE id_user = $1 AND pin_code = $2 AND active = TRUE`

	var count int
	err := r.DB.QueryRowContext(ctx, query, userID, pin).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

//User Functions
func (r *FaceRPCASQLRepository) GetByID(ctx context.Context, userID int) (*domain.User, error) {
	query := `SELECT id_user, nama, nip, COALESCE(email, ''), active, COALESCE(created_by, ''), COALESCE(dataset_path, '')
	          FROM ms_user WHERE id_user = $1`

	var u domain.User
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(
		&u.IDUser, &u.Name, &u.NIP, &u.Email, &u.Active, &u.CreatedBy, &u.DatasetPath,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *FaceRPCASQLRepository) CreateUser(ctx context.Context, name, nip, email, createdBy, datasetPath string) (int, error) {
	query := `INSERT INTO ms_user (nama, nip, email, active, created_by, dataset_path)
	          VALUES ($1, $2, $3, TRUE, $4, $5)
	          RETURNING id_user`

	var newID int
	err := r.DB.QueryRowContext(ctx, query, name, nip, email, createdBy, datasetPath).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}



func (r *FaceRPCASQLRepository) GetUserNameByID(ctx context.Context, userID int) (string, error) {
	query := `SELECT nama FROM ms_user WHERE id_user = $1 AND active = TRUE`

	var name string
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&name)
	if err != nil {
		return "Unknown", nil
	}
	return name, nil
}

func (r *FaceRPCASQLRepository) GetAllActiveUsers(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id_user, nama, nip, COALESCE(email, ''), active, COALESCE(created_by, ''), COALESCE(dataset_path, '')
	          FROM ms_user WHERE active = TRUE`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.IDUser, &u.Name, &u.NIP, &u.Email, &u.Active, &u.CreatedBy, &u.DatasetPath); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *FaceRPCASQLRepository) GetUserEmailByID(ctx context.Context, userID int) (string, error) {
	query := `SELECT COALESCE(email, '') FROM ms_user WHERE id_user = $1 AND active = TRUE`

	var email string
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&email)
	if err != nil {
		return "", nil
	}
	return email, nil
}

func (r *FaceRPCASQLRepository) GetAll(ctx context.Context) ([]domain.UserListRow, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id_user, nama, nip, active FROM ms_user ORDER BY id_user ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserListRow
	for rows.Next() {
		var u domain.UserListRow
		if err := rows.Scan(&u.ID, &u.Name, &u.NIP, &u.Active); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *FaceRPCASQLRepository) Update(ctx context.Context, userID int, name, nip string, active bool) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE ms_user SET nama=$1, nip=$2, active=$3, updated_at=NOW() WHERE id_user=$4`,
		name, nip, active, userID)
	return err
}

func (r *FaceRPCASQLRepository) Delete(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM ms_user WHERE id_user=$1`, userID)
	return err
}


//Attendance Functions
func (r *FaceRPCASQLRepository) SaveAttendance(ctx context.Context, record domain.AttendanceRecord) error {
	query := `INSERT INTO tbl_log_attendance 
	          (id_user, presence_time, confidence, latency_ms, id_status)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.DB.ExecContext(ctx, query,
		record.IDUser,
		record.PresenceTime,
		record.Confidence,
		record.LatencyMS,
		record.Status,
	)
	return err
}

// GetLoggedInUsers returns users who are currently active today (Hadir or Telat, not yet Keluar/Cuti).
func (r *FaceRPCASQLRepository) GetLoggedInUsers(ctx context.Context, today time.Time) ([]int, error) {
	query := `SELECT DISTINCT a.id_user 
	          FROM tbl_log_attendance a
	          WHERE DATE(a.presence_time) = $1 
	          AND a.id_status IN ($2, $3)
	          AND a.id_user NOT IN (
	              SELECT id_user FROM tbl_log_attendance 
	              WHERE DATE(presence_time) = $1 AND id_status IN ($4, $5)
	          )`

	rows, err := r.DB.QueryContext(ctx, query, today,
		constant.StatusHadir, constant.StatusTelat,
		constant.StatusKeluar, constant.StatusCuti)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, rows.Err()
}

// GetUsersWhoWorkedOnDate returns users who had status Hadir or Telat on a given date
// (regardless of whether they were later forced to Keluar). Used for overnight Jira catch-up.
func (r *FaceRPCASQLRepository) GetUsersWhoWorkedOnDate(ctx context.Context, date time.Time) ([]int, error) {
	query := `SELECT DISTINCT a.id_user 
	          FROM tbl_log_attendance a
	          WHERE DATE(a.presence_time) = $1 
	          AND a.id_status IN ($2, $3)`

	rows, err := r.DB.QueryContext(ctx, query, date,
		constant.StatusHadir, constant.StatusTelat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, rows.Err()
}

func (r *FaceRPCASQLRepository) GetUserStatusToday(ctx context.Context, userID int, today time.Time) (*int, error) {
	query := `SELECT id_status FROM tbl_log_attendance 
	          WHERE id_user = $1 AND DATE(presence_time) = $2
	          ORDER BY presence_time DESC
	          LIMIT 1`

	var status int
	err := r.DB.QueryRowContext(ctx, query, userID, today).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}

// InitializeDailyStatus inserts status 0 (Menunggu Kehadiran) for all active users
// who don't yet have any attendance record today. Called at 05:00 AM.
func (r *FaceRPCASQLRepository) InitializeDailyStatus(ctx context.Context, today time.Time) error {
	query := `INSERT INTO tbl_log_attendance (id_user, presence_time, confidence, latency_ms, id_status)
	          SELECT u.id_user, $1, 0, 0, $2
	          FROM ms_user u
	          WHERE u.active = TRUE
	          AND u.id_user NOT IN (
	              SELECT DISTINCT id_user FROM tbl_log_attendance 
	              WHERE DATE(presence_time) = $3
	          )`

	_, err := r.DB.ExecContext(ctx, query, time.Now(), constant.StatusMenunggu, today)
	return err
}

// MarkAbsentUsers changes status 0 (Menunggu) → status 2 (Tidak Hadir) at 09:00 AM cutoff.
// Users who already have Hadir, Telat, or Cuti are not affected.
func (r *FaceRPCASQLRepository) MarkAbsentUsers(ctx context.Context, today time.Time, cutoffTime time.Time) error {
	query := `INSERT INTO tbl_log_attendance (id_user, presence_time, confidence, latency_ms, id_status)
	          SELECT u.id_user, $1, 0, 0, $2
	          FROM ms_user u
	          WHERE u.active = TRUE
	          AND u.id_user NOT IN (
	              SELECT DISTINCT id_user FROM tbl_log_attendance 
	              WHERE DATE(presence_time) = $3
	              AND id_status IN ($4, $5, $6, $7)
	          )`

	_, err := r.DB.ExecContext(ctx, query, cutoffTime, constant.StatusTidakHadir, today,
		constant.StatusHadir, constant.StatusTelat, constant.StatusCuti, constant.StatusTidakHadir)
	return err
}

// ForceLogoffUsers forces all users with status Hadir (1) or Telat (3) to Keluar (4) at 20:00.
func (r *FaceRPCASQLRepository) ForceLogoffUsers(ctx context.Context, today time.Time) ([]int, error) {
	findQuery := `SELECT DISTINCT a.id_user 
	              FROM tbl_log_attendance a
	              WHERE DATE(a.presence_time) = $1 
	              AND a.id_status IN ($2, $3)
	              AND a.id_user NOT IN (
	                  SELECT id_user FROM tbl_log_attendance 
	                  WHERE DATE(presence_time) = $1 AND id_status IN ($4, $5)
	              )`

	rows, err := r.DB.QueryContext(ctx, findQuery, today,
		constant.StatusHadir, constant.StatusTelat,
		constant.StatusKeluar, constant.StatusCuti)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return nil, nil
	}

	now := time.Now()
	insertQuery := `INSERT INTO tbl_log_attendance 
	                (id_user, presence_time, confidence, latency_ms, id_status)
	                VALUES ($1, $2, 0, 0, $3)`

	for _, uid := range userIDs {
		if _, err := r.DB.ExecContext(ctx, insertQuery, uid, now, constant.StatusKeluar); err != nil {
			return userIDs, err // return partial list on error
		}
	}

	return userIDs, nil
}

// SetUserCuti inserts a Cuti (5) status record for a user today.
func (r *FaceRPCASQLRepository) SetUserCuti(ctx context.Context, userID int, today time.Time) error {
	query := `INSERT INTO tbl_log_attendance (id_user, presence_time, confidence, latency_ms, id_status)
	          VALUES ($1, $2, 0, 0, $3)`

	_, err := r.DB.ExecContext(ctx, query, userID, time.Now(), constant.StatusCuti)
	return err
}

func (r *FaceRPCASQLRepository) UpsertKpiHarian(ctx context.Context, kpi domain.KPIHarian) error {
	query := `INSERT INTO tbl_kpi_harian 
	          (id_tiket_jira, story_jira, id_status_tiket, nama_tiket, id_user, date)
	          VALUES ($1, $2, $3, $4, $5, $6)
	          ON CONFLICT (id_tiket_jira, id_user, date) 
	          DO UPDATE SET id_status_tiket = EXCLUDED.id_status_tiket, story_jira = EXCLUDED.story_jira`

	_, err := r.DB.ExecContext(ctx, query,
		kpi.IDTiketJira,
		kpi.StoryJira,
		kpi.IDStatusTiket,
		kpi.NamaTiket,
		kpi.IDUser,
		kpi.Date,
	)
	return err
}

func (r *FaceRPCASQLRepository) GetDoneTasksCount(ctx context.Context, userID int, today time.Time) (int, error) {
	query := `SELECT COUNT(*) FROM tbl_kpi_harian 
	          WHERE id_user = $1 AND date = $2 AND id_status_tiket = 4`

	var count int
	err := r.DB.QueryRowContext(ctx, query, userID, today).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FaceRPCASQLRepository) HasDoneTasks(ctx context.Context, userID int, today time.Time) (bool, error) {
	count, err := r.GetDoneTasksCount(ctx, userID, today)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FaceRPCASQLRepository) UpdateAkumulasiKPI(ctx context.Context, userID int, akumulasiKPI int) error {
	query := `UPDATE ms_user SET akumulasi_kpi = $1, updated_at = NOW() WHERE id_user = $2`
	_, err := r.DB.ExecContext(ctx, query, akumulasiKPI, userID)
	return err
}

func (r *FaceRPCASQLRepository) GetRecentJira(ctx context.Context, limit int) ([]domain.JiraHistoryRow, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT k.id_kpi_harian, k.id_tiket_jira, k.nama_tiket, s.deskripsi, k.date
		 FROM tbl_kpi_harian k
		 JOIN ms_status_kpi s ON k.id_status_tiket=s.id_status_tiket
		 ORDER BY k.date DESC, k.id_kpi_harian DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.JiraHistoryRow
	for rows.Next() {
		var row domain.JiraHistoryRow
		var d time.Time
		if rows.Scan(&row.ID, &row.Ticket, &row.Name, &row.Status, &d) != nil {
			continue
		}
		row.Date = d.Format("02/01/2006")
		records = append(records, row)
	}
	return records, rows.Err()
}

// GetJiraHistory returns paginated Jira ticket history.
func (r *FaceRPCASQLRepository) GetJiraHistory(ctx context.Context, page, perPage int) (*domain.PaginatedResult[domain.JiraHistoryRow], error) {
	offset := (page - 1) * perPage

	var total int
	r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM tbl_kpi_harian`).Scan(&total)
	tp := int(math.Ceil(float64(total) / float64(perPage)))

	rows, err := r.DB.QueryContext(ctx,
		`SELECT k.id_kpi_harian, k.id_tiket_jira, k.nama_tiket, s.deskripsi, k.date
		 FROM tbl_kpi_harian k
		 JOIN ms_status_kpi s ON k.id_status_tiket=s.id_status_tiket
		 ORDER BY k.date DESC, k.id_kpi_harian DESC LIMIT $1 OFFSET $2`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.JiraHistoryRow
	for rows.Next() {
		var row domain.JiraHistoryRow
		var d time.Time
		if rows.Scan(&row.ID, &row.Ticket, &row.Name, &row.Status, &d) != nil {
			continue
		}
		row.Date = d.Format("02/01/2006 15:04:05")
		records = append(records, row)
	}
	return &domain.PaginatedResult[domain.JiraHistoryRow]{
		Records: records, Total: total,
		Page: page, PerPage: perPage, TotalPages: tp,
	}, nil
}

// GetKpiAccumulation returns paginated KPI accumulation data.
func (r *FaceRPCASQLRepository) GetKpiAccumulation(ctx context.Context, page, perPage int) (*domain.PaginatedResult[domain.KpiAccumulationRow], error) {
	offset := (page - 1) * perPage

	var total int
	r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM ms_user WHERE akumulasi_kpi > 0`).Scan(&total)
	tp := int(math.Ceil(float64(total) / float64(perPage)))

	rows, err := r.DB.QueryContext(ctx,
		`SELECT id_user, nama, akumulasi_kpi
		 FROM ms_user WHERE akumulasi_kpi > 0
		 ORDER BY akumulasi_kpi DESC LIMIT $1 OFFSET $2`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.KpiAccumulationRow
	for rows.Next() {
		var row domain.KpiAccumulationRow
		if rows.Scan(&row.No, &row.Name, &row.TotalTiket) != nil {
			continue
		}
		records = append(records, row)
	}
	return &domain.PaginatedResult[domain.KpiAccumulationRow]{
		Records: records, Total: total,
		Page: page, PerPage: perPage, TotalPages: tp,
	}, nil
}


func (r *FaceRPCASQLRepository) SaveRaw(ctx context.Context, rec domain.BulkAttendanceRecord) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO tbl_log_attendance (id_user, presence_time, confidence, latency_ms, id_status)
		 VALUES ($1,$2,$3,$4,$5)`,
		rec.UserID, rec.PresenceTime, rec.Confidence, rec.LatencyMs, rec.StatusID)
	return err
}

func (r *FaceRPCASQLRepository) SaveBulk(ctx context.Context, records []domain.BulkAttendanceRecord) (int, error) {
	synced := 0
	for _, rec := range records {
		_, err := r.DB.ExecContext(ctx,
			`INSERT INTO tbl_log_attendance (id_user, presence_time, confidence, latency_ms, id_status)
			 VALUES ($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING`,
			rec.UserID, rec.PresenceTime, rec.Confidence, rec.LatencyMs, rec.StatusID)
		if err == nil {
			synced++
		}
	}
	return synced, nil
}

func (r *FaceRPCASQLRepository) GetRecentPresence(ctx context.Context, limit int) ([]domain.RecentPresence, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT a.id_log, u.nama, a.presence_time
		 FROM tbl_log_attendance a JOIN ms_user u ON a.id_user=u.id_user
		 ORDER BY a.presence_time DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.RecentPresence
	for rows.Next() {
		var r domain.RecentPresence
		var t time.Time
		if err := rows.Scan(&r.ID, &r.Name, &t); err != nil {
			continue
		}
		r.Time = t.Format("15:04:05")
		records = append(records, r)
	}
	return records, rows.Err()
}

func (r *FaceRPCASQLRepository) GetHistory(ctx context.Context, page, perPage int) (*domain.PaginatedResult[domain.AttendanceHistoryRow], error) {
	offset := (page - 1) * perPage

	var total int
	r.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM tbl_log_attendance`).Scan(&total)
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	rows, err := r.DB.QueryContext(ctx,
		`SELECT a.id_log, u.nama, u.nip, s.deskripsi, a.confidence, a.latency_ms, a.presence_time
		 FROM tbl_log_attendance a
		 JOIN ms_user u ON a.id_user=u.id_user
		 JOIN ms_status s ON a.id_status=s.id_status
		 ORDER BY a.presence_time DESC LIMIT $1 OFFSET $2`, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []domain.AttendanceHistoryRow
	for rows.Next() {
		var row domain.AttendanceHistoryRow
		var conf, lat float64
		var t time.Time
		if err := rows.Scan(&row.No, &row.Name, &row.NIP, &row.Status, &conf, &lat, &t); err != nil {
			continue
		}
		row.Confidence = int(math.Round(conf * 100))
		row.Latency = strconv.FormatFloat(lat, 'f', 3, 64) + " ms"
		row.Time = t.Format("02/01/2006 15:04:05")
		records = append(records, row)
	}

	return &domain.PaginatedResult[domain.AttendanceHistoryRow]{
		Records: records, Total: total,
		Page: page, PerPage: perPage, TotalPages: totalPages,
	}, nil
}

func (r *FaceRPCASQLRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats := &domain.DashboardStats{}

	rows, err := r.DB.QueryContext(ctx,
		`SELECT id_status, COUNT(DISTINCT id_user) FROM tbl_log_attendance
		 WHERE DATE(presence_time)=CURRENT_DATE GROUP BY id_status`)
	if err != nil {
		return stats, nil
	}
	defer rows.Close()

	for rows.Next() {
		var sid, cnt int
		if rows.Scan(&sid, &cnt) != nil {
			continue
		}
		switch sid {
		case constant.StatusMenunggu:
			stats.Menunggu = cnt
		case constant.StatusHadir:
			stats.Hadir = cnt
		case constant.StatusTidakHadir:
			stats.TidakHadir = cnt
		case constant.StatusTelat:
			stats.Telat = cnt
		case constant.StatusKeluar:
			stats.Keluar = cnt
		case constant.StatusCuti:
			stats.Cuti = cnt
		}
	}
	
	// Calculate total unique users who are either Hadir or Telat
	err = r.DB.QueryRowContext(ctx,
		`SELECT COUNT(DISTINCT id_user) FROM tbl_log_attendance
		 WHERE DATE(presence_time)=CURRENT_DATE AND id_status IN ($1, $2)`,
		constant.StatusHadir, constant.StatusTelat).Scan(&stats.TotalHadir)
	if err != nil {
		return stats, nil
	}

	stats.UpdatedAt = time.Now().Format("15:04:05, 02 January 2006")
	return stats, nil
}

func (r *FaceRPCASQLRepository) GetHourlyFrequency(ctx context.Context) (*domain.HourlyFrequency, error) {
	freq := &domain.HourlyFrequency{
		Labels: make([]string, 18),
		Data:   make([]int, 18),
	}
	for i := 0; i < 18; i++ {
		freq.Labels[i] = fmt.Sprintf("%02d:00", i+5)
	}

	rows, err := r.DB.QueryContext(ctx,
		`SELECT EXTRACT(HOUR FROM presence_time)::int, COUNT(DISTINCT id_user)
		 FROM tbl_log_attendance
		 WHERE DATE(presence_time)=CURRENT_DATE AND id_status IN ($1, $2)
		   AND EXTRACT(HOUR FROM presence_time) BETWEEN 5 AND 22
		 GROUP BY 1 ORDER BY 1`,
		constant.StatusHadir, constant.StatusTelat)
	if err != nil {
		return freq, nil
	}
	defer rows.Close()

	for rows.Next() {
		var hr, cnt int
		if rows.Scan(&hr, &cnt) == nil {
			if idx := hr - 5; idx >= 0 && idx < 18 {
				freq.Data[idx] = cnt
			}
		}
	}
	return freq, nil
}

func (r *FaceRPCASQLRepository) GetTotalUserCount(ctx context.Context) (int, error) {
	var count int
	err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM ms_user").Scan(&count)
	return count, err
}

func (r *FaceRPCASQLRepository) GetTodayAttendeeCount(ctx context.Context) (int, error) {
	var count int
	// Unique users with status Hadir (1) or Telat (3) today
	query := `SELECT COUNT(DISTINCT id_user) FROM tbl_log_attendance 
	          WHERE DATE(presence_time) = CURRENT_DATE AND id_status IN (1, 3)`
	err := r.DB.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}
