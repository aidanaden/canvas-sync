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
