package canvas

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
	"github.com/pterm/pterm"
)

const apiPath = "/api/v1"

type CanvasClient struct {
	client      *http.Client
	canvasPath  *url.URL
	apiPath     *url.URL
	accessToken string
}

func NewClient(rawUrl string, accessToken string) *CanvasClient {
	schemas := []string{"http://", "https://"}
	canvasHost := ""
	for _, schema := range schemas {
		if strings.Contains(rawUrl, schema) {
			splits := strings.Split(rawUrl, schema)
			canvasHost = splits[1]
			break
		}
	}
	if canvasHost == "" {
		canvasHost = rawUrl
	}
	canvasPath := url.URL{
		Scheme: "https",
		Host:   canvasHost,
	}
	apiPath := url.URL{
		Scheme: canvasPath.Scheme,
		Host:   canvasPath.Host,
		Path:   apiPath,
	}
	httpClient := http.Client{}
	return &CanvasClient{
		client:      &httpClient,
		accessToken: accessToken,
		canvasPath:  &canvasPath,
		apiPath:     &apiPath,
	}
}

func (c *CanvasClient) GetActiveEnrolledCoursesURL() url.URL {
	return url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + "/users/self/courses",
		RawQuery: url.Values{
			"enrollment_state": {"active"},
		}.Encode(),
	}
}

func (c *CanvasClient) GetActiveEnrolledCourses() ([]nodes.CourseNode, error) {
	coursesUrl := c.GetActiveEnrolledCoursesURL()
	// courses request
	req, err := http.NewRequest("GET", coursesUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	courseJson := utils.ExtractResponseToString(resp)
	var courses []nodes.CourseNode
	json.Unmarshal([]byte(courseJson), &courses)
	if strings.Contains(courseJson, "user authorisation required") {
		return nil, fmt.Errorf("existing config invalid, please run 'canvas-sync init'")
	}

	return courses, nil
}

func (c *CanvasClient) getCourseUrl(id int) url.URL {
	courseUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + fmt.Sprintf("/courses/%d/folders/root", id),
	}
	return courseUrl
}

func (c *CanvasClient) GetCourseRootFolder(courseId int) (*nodes.DirectoryNode, error) {
	courseUrl := c.getCourseUrl(courseId)
	req, err := http.NewRequest("GET", courseUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		return nil, err
	}
	rootRes, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	rootJson := utils.ExtractResponseToString(rootRes)
	var rootNode *nodes.DirectoryNode
	if err := json.Unmarshal([]byte(rootJson), &rootNode); err != nil {
		return nil, err
	}
	return rootNode, nil
}

var PER_PAGE = 100

func (c *CanvasClient) RecurseDirectoryNode(node *nodes.DirectoryNode, parent *nodes.DirectoryNode) error {
	dir := ""
	if parent != nil {
		dir = filepath.Join(parent.Directory)
	}
	if node == nil {
		return errors.New("cannot recurse nil directory node")
	}
	dir = filepath.Join(dir, node.Name)
	node.Directory = dir

	if node.FilesCount > 0 {
		parsedFileUrl, err := url.Parse(node.FilesUrl)
		if err != nil {
			return err
		}

		var allFiles []*nodes.FileNode
		page := 0
		for {
			page += 1
			nolimitFileUrl := url.URL{
				Scheme: parsedFileUrl.Scheme,
				Host:   parsedFileUrl.Host,
				Path:   parsedFileUrl.Path,
				RawQuery: url.Values{
					"page":     {strconv.Itoa(page)},
					"per_page": {strconv.Itoa(PER_PAGE)},
				}.Encode(),
			}
			fileReq, err := http.NewRequest("GET", nolimitFileUrl.String(), nil)
			if err != nil {
				return err
			}
			utils.SetQueryAccessToken(fileReq, c.accessToken)
			filesRes, err := c.client.Do(fileReq)
			if err != nil {
				return err
			}
			filesJson := utils.ExtractResponseToString(filesRes)
			var files []*nodes.FileNode
			json.Unmarshal([]byte(filesJson), &files)
			allFiles = append(allFiles, files...)
			// break if less than 100 files, otherwise query next page
			if len(files) < PER_PAGE {
				break
			}
		}
		for f := range allFiles {
			allFiles[f].Directory = filepath.Join(dir, allFiles[f].Display_name)
		}
		node.FileNodes = allFiles
	}

	if node.FoldersCount > 0 {
		parsedFolderUrl, err := url.Parse(node.FoldersUrl)
		if err != nil {
			return err
		}

		var allFolders []*nodes.DirectoryNode
		page := 0
		for {
			page += 1
			nolimitFolderUrl := url.URL{
				Scheme: parsedFolderUrl.Scheme,
				Host:   parsedFolderUrl.Host,
				Path:   parsedFolderUrl.Path,
				RawQuery: url.Values{
					"page":     {strconv.Itoa(page)},
					"per_page": {strconv.Itoa(PER_PAGE)},
				}.Encode(),
			}
			folderReq, err := http.NewRequest("GET", nolimitFolderUrl.String(), nil)
			if err != nil {
				return err
			}
			utils.SetQueryAccessToken(folderReq, c.accessToken)
			foldersRes, err := c.client.Do(folderReq)
			if err != nil {
				return err
			}
			foldersJson := utils.ExtractResponseToString(foldersRes)
			var folders []*nodes.DirectoryNode
			json.Unmarshal([]byte(foldersJson), &folders)
			allFolders = append(allFolders, folders...)
			// break if less than 100 folders, otherwise query next page
			if len(folders) < PER_PAGE {
				break
			}
		}
		for fi := range allFolders {
			if err := c.RecurseDirectoryNode(allFolders[fi], node); err != nil {
				return err
			}
		}
		node.FolderNodes = allFolders
	}

	return nil
}

func (c *CanvasClient) downloadFileNode(node *nodes.FileNode) error {
	if node == nil {
		return errors.New("cannot download file without file node")
	}
	file, err := os.Create(node.Directory)
	if err != nil {
		return err
	}
	defer file.Close()
	res, err := c.client.Get(node.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (c *CanvasClient) RecursiveCreateNode(node *nodes.DirectoryNode, updateNumDownloads func(numDownloads int)) error {
	if node == nil {
		return errors.New("cannot recurse nil directory node")
	}
	if err := os.MkdirAll(node.Directory, 0755); err != nil {
		return err
	}
	var wg sync.WaitGroup
	numDownloads := 0
	for j := range node.FileNodes {
		if node.FileNodes[j] == nil {
			continue
		}
		wg.Add(1)
		numDownloads += 1
		go func(i int) {
			defer wg.Done()
			var err error
			err = c.downloadFileNode(node.FileNodes[i])
			for err != nil {
				pterm.Error.Printfln("Error downloading file %s: %s", node.FileNodes[i].Display_name, err.Error())
				err = c.downloadFileNode(node.FileNodes[i])
			}
		}(j)
	}
	updateNumDownloads(numDownloads)
	for d := range node.FolderNodes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := c.RecursiveCreateNode(node.FolderNodes[i], updateNumDownloads); err != nil {
				pterm.Error.Printfln("Error downloading folder %s: %s", node.FileNodes[i].Display_name, err.Error())
			}
		}(d)
	}
	wg.Wait()
	return nil
}

func (c *CanvasClient) RecursiveUpdateNode(node *nodes.DirectoryNode, updateStaleFiles bool, updateNumDownloads func(numDownloads int)) error {
	if node == nil {
		return errors.New("cannot recurse nil directory node")
	}
	// create directory if doesnt exist
	if _, err := os.Stat(node.Directory); os.IsNotExist(err) {
		if err := os.MkdirAll(node.Directory, 0755); err != nil {
			return err
		}
	}
	var wg sync.WaitGroup
	numDownloads := 0
	for j := range node.FileNodes {
		if node.FileNodes[j] == nil {
			continue
		}
		file, err := os.Stat(node.FileNodes[j].Directory)
		if err != nil {
			wg.Add(1)
			numDownloads += 1
			go func(i int) {
				defer wg.Done()
				var err error
				err = c.downloadFileNode(node.FileNodes[i])
				for err != nil {
					pterm.Error.Printfln("Error downloading file %s: %s", node.FileNodes[i].Display_name, err.Error())
					err = c.downloadFileNode(node.FileNodes[i])
				}
			}(j)
		} else {
			if updateStaleFiles && file.ModTime().Unix() < node.FileNodes[j].UpdatedAt.Unix() {
				wg.Add(1)
				numDownloads += 1
				go func(i int) {
					defer wg.Done()
					if err := c.downloadFileNode(node.FileNodes[i]); err != nil {
						pterm.Error.Printfln("Error downloading file %s: %s", node.FileNodes[i].Display_name, err.Error())
					}
				}(j)
			}
		}
	}
	updateNumDownloads(numDownloads)
	for d := range node.FolderNodes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := c.RecursiveUpdateNode(node.FolderNodes[i], updateStaleFiles, updateNumDownloads); err != nil {
				pterm.Error.Printfln("Error updating folder %s: %s", node.FileNodes[i].Display_name, err.Error())
			}
		}(d)
	}
	wg.Wait()
	return nil
}

func extractEventFromString(rawJson string) ([]nodes.EventNode, error) {
	var events []nodes.EventNode
	json.Unmarshal([]byte(rawJson), &events)
	if strings.Contains(rawJson, "user authorisation required") {
		return nil, errors.New("invalid auth cookies/access token, request unauthorized")
	}
	return events, nil
}

func (c *CanvasClient) GetRecentCalendarEvents() ([]nodes.EventNode, error) {
	now := utils.TimestampToJavaScriptISO(time.Now())
	eventsUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + "/planner/items",
		RawQuery: url.Values{
			"end_date": {now},
			"order":    {"asc"},
			"per_page": {"50"},
		}.Encode(),
	}

	// events request
	req, err := http.NewRequest("GET", eventsUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	eventsJson := utils.ExtractResponseToString(resp)
	events, err := extractEventFromString(eventsJson)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (c *CanvasClient) GetIncomingCalendarEvents() ([]nodes.EventNode, error) {
	now := utils.TimestampToJavaScriptISO(time.Now())
	eventsUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + "/planner/items",
		RawQuery: url.Values{
			"start_date": {now},
			"per_page":   {"50"},
		}.Encode(),
	}

	// events request
	req, err := http.NewRequest("GET", eventsUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	eventsJson := utils.ExtractResponseToString(resp)
	events, err := extractEventFromString(eventsJson)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func extractPeopleFromString(rawJson string) ([]nodes.PersonNode, error) {
	var people []nodes.PersonNode
	json.Unmarshal([]byte(rawJson), &people)
	if strings.Contains(rawJson, "user authorisation required") {
		return nil, errors.New("invalid auth cookies/access token, request unauthorized")
	}
	return people, nil
}

func (c *CanvasClient) GetCoursePeople(code string) ([]nodes.PersonNode, error) {
	rawCourses, err := c.GetActiveEnrolledCourses()
	if err != nil {
		return nil, err
	}
	var courseId int
	for _, raw := range rawCourses {
		if strings.EqualFold(code, raw.CourseCode) {
			courseId = raw.ID
		}
	}
	if courseId == 0 {
		return nil, errors.New("course not found")
	}

	peopleUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + "/courses/" + strconv.Itoa(courseId) + "/users",
		RawQuery: url.Values{
			"include[]": {"avatar_url", "observed_users"},
			"per_page":  {"100"},
			"page":      {"1"},
		}.Encode(),
	}

	req, err := http.NewRequest("GET", peopleUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	peopleJson := utils.ExtractResponseToString(resp)
	people, err := extractPeopleFromString(peopleJson)
	if err != nil {
		return nil, err
	}
	return people, nil
}

func extractAnnouncementsFromString(rawJson string) ([]nodes.AnnouncementNode, error) {
	var announcements []nodes.AnnouncementNode
	json.Unmarshal([]byte(rawJson), &announcements)
	if strings.Contains(rawJson, "user authorisation required") {
		return nil, errors.New("invalid auth cookies/access token, request unauthorized")
	}
	return announcements, nil
}

func (c *CanvasClient) GetCourseAnnouncements(code string) ([]nodes.AnnouncementNode, error) {
	rawCourses, err := c.GetActiveEnrolledCourses()
	if err != nil {
		return nil, err
	}
	var courseId int
	for _, raw := range rawCourses {
		if strings.EqualFold(code, raw.CourseCode) {
			courseId = raw.ID
		}
	}
	if courseId == 0 {
		return nil, errors.New("course not found")
	}

	peopleUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + "/courses/" + strconv.Itoa(courseId) + "/discussion_topics",
		RawQuery: url.Values{
			"only_announcements": {"true"},
			"per_page":           {"50"},
			"page":               {"1"},
		}.Encode(),
	}

	req, err := http.NewRequest("GET", peopleUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	announcementsJson := utils.ExtractResponseToString(resp)
	announcements, err := extractAnnouncementsFromString(announcementsJson)
	if err != nil {
		return nil, err
	}
	return announcements, nil
}

func (c *CanvasClient) GetCourseGrades(code string) error {
	// rawCourses, err := c.GetActiveEnrolledCourses()
	// if err != nil {
	// 	return err
	// }
	// var courseId int
	// for _, raw := range rawCourses {
	// 	if code == raw.CourseCode {
	// 		courseId = raw.ID
	// 	}
	// }
	return nil
}
