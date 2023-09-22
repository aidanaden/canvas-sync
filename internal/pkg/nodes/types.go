package nodes

import "time"

type CourseNode struct {
	ID                               int    `json:"id"`
	Name                             string `json:"name"`
	CourseCode                       string `json:"course_code"`
	RestrictEnrollmentsToCourseDates bool   `json:"restrict_enrollments_to_course_dates"`
	RootDirectory                    *DirectoryNode
}

type FileNode struct {
	Directory    string
	Display_name string      `json:"display_name"`
	UpdatedAt    time.Time   `json:"updated_at"`
	ContentType  string      `json:"content-type"`
	Url          string      `json:"url"`
	Locked       bool        `json:"locked"`
	ModifiedAt   time.Ticker `json:"modified_at"`
}

type DirectoryNode struct {
	Directory     string
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Updated_at    time.Time `json:"updated_at"`
	Locked        bool      `json:"locked"`
	FoldersUrl    string    `json:"folders_url"`
	FoldersCount  int       `json:"folders_count"`
	FilesUrl      string    `json:"files_url"`
	FilesCount    int       `json:"files_count"`
	LockedForUser bool      `json:"locked_for_user"`
	HiddenForUser bool      `json:"hidden_for_user"`
	FolderNodes   []*DirectoryNode
	FileNodes     []*FileNode
}

type BasePlannableNode struct {
	Title string `json:"title"`
}

type AnnouncementPlannableNode struct {
	ReadState   string `json:"read_state"`
	UnreadCount int    `json:"unread_count"`
}

type AssignmentPlannableNode struct {
	DueAt          time.Time `json:"due_at"`
	PointsPossible float64   `json:"points_possible"`
}

type LivePlannableNode struct {
	LocationAddress string `json:"location_address"`
}

type ZoomPlannableNode struct {
	OnlineMeetingUrl string `json:"online_meeting_url"`
}

type EventPlannableNode struct {
	AllDay       bool      `json:"all_day"`
	LocationName string    `json:"location_name"`
	StartAt      time.Time `json:"start_at"`
	Description  string    `json:"description"`
	*LivePlannableNode
	*ZoomPlannableNode
}

type PlannableNode struct {
	BasePlannableNode
	*AnnouncementPlannableNode
	*AssignmentPlannableNode
	*EventPlannableNode
}

type EventNode struct {
	ContextName   string        `json:"context_name"`
	ContextType   string        `json:"context_type"`
	CourseId      int           `json:"course_id"`
	HtmlUrl       string        `json:"html_url"`
	NewActivity   bool          `json:"new_activity"`
	Plannable     PlannableNode `json:"plannable"`
	PlannableDate time.Time     `json:"plannable_date"`
	PlannableType string        `json:"plannable_type"`
}

type PersonNode struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
	AvatarUrl string `json:"avatar_url"`
}
