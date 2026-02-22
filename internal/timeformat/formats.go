package timeformat

const (
	Date      = "2006-01-02"
	DateTime  = "2006-01-02T15:04:05Z07:00"
	DateTimeZ = "2006-01-02T15:04:05Z0700"
	DateTimeNoTZ = "2006-01-02T15:04:05"
	DateTimeSpace = "2006-01-02 15:04:05Z07:00"
	DateTimeSpaceZ = "2006-01-02 15:04:05Z0700"
	DateTimeSpaceNoTZ = "2006-01-02 15:04:05"
)

var LenientLayouts = []string{
	Date,
	DateTime,
	DateTimeZ,
	DateTimeNoTZ,
	DateTimeSpace,
	DateTimeSpaceZ,
	DateTimeSpaceNoTZ,
}
