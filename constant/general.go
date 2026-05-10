package constant

const (
	TimeLocation     = "Asia/Jakarta"
	MethodNotAllowed = "Method Not Allowed"

	// Attendance Status IDs
	StatusMenunggu   = 0 // Menunggu Kehadiran
	StatusHadir      = 1 // Hadir (on-time, before 09:00)
	StatusTidakHadir = 2 // Tidak Hadir (absent)
	StatusTelat      = 3 // Telat Hadir (late, after 09:00)
	StatusKeluar     = 4 // Keluar (clocked out / end of day)
	StatusCuti       = 5 // Cuti (paid leave, set by admin)
)

type ResultError struct {
	Code InternalError
	Err  error
}
