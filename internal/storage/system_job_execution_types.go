package storage

import "time"

const (
	DefaultSystemJobHistoryRetentionDays        = 30
	DefaultRoutineSystemJobRetentionHours       = 24
	MinSystemJobScheduleIntervalSeconds   int32 = 15
)

type SystemJobScheduleDefinition struct {
	ID                    string
	Name                  string
	Category              string
	Description           string
	Kind                  string
	Queue                 string
	IntervalSeconds       int32
	IntervalConfigurable  bool
	HistoryPolicy         string
	Automatic             bool
	ManualActionAvailable bool
}

type SystemJobSchedule struct {
	ID                    string
	Name                  string
	Category              string
	Description           string
	Kind                  string
	Queue                 string
	IntervalSeconds       int32
	IntervalConfigurable  bool
	HistoryPolicy         string
	Automatic             bool
	ManualActionAvailable bool
	Enabled               bool
	Paused                bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
	NextRunAt             *time.Time
	ActiveRiverJobID      *int64
	ActiveStatus          string
	ActiveProgressPercent *int32
	ActiveProgressLabel   string
	ActiveInfoMessage     string
	LastRiverJobID        *int64
	LastStatus            string
	LastCreatedAt         *time.Time
	LastFinalizedAt       *time.Time
}

type SystemJobExecution struct {
	RiverJobID      int64
	ScheduleID      string
	Classification  string
	HistoryPolicy   string
	Status          string
	Kind            string
	Queue           string
	Attempt         int32
	MaxAttempts     int32
	Priority        int32
	ProgressPercent *int32
	ProgressLabel   string
	Args            string
	Metadata        string
	Errors          string
	InfoMessage     string
	ScheduledAt     time.Time
	CreatedAt       time.Time
	AttemptedAt     *time.Time
	FinalizedAt     *time.Time
	UpdatedAt       time.Time
}

type SystemJobExecutionInput struct {
	RiverJobID     int64
	ScheduleID     string
	Classification string
	Status         string
	Kind           string
	Queue          string
	Attempt        int32
	MaxAttempts    int32
	Priority       int32
	Args           []byte
	Metadata       []byte
	Errors         []byte
	InfoMessage    string
	ScheduledAt    time.Time
	CreatedAt      time.Time
	AttemptedAt    *time.Time
	FinalizedAt    *time.Time
}

type SystemJobExecutionFilters struct {
	States         []string
	ScheduleID     string
	Kind           string
	Queue          string
	Query          string
	IncludeRoutine bool
	Before         *time.Time
	Limit          int32
}

type SystemJobExecutionLog struct {
	ID         int64
	RiverJobID int64
	Severity   string
	Message    string
	Data       map[string]any
	CreatedAt  time.Time
}

type SystemJobHistorySettings struct {
	RetentionDays         int32
	RoutineRetentionHours int32
}
