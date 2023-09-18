package canvas

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/aidanaden/canvas-sync/internal/pkg/nodes"
	"github.com/aidanaden/canvas-sync/internal/pkg/utils"
)

const apiPath = "/api/v1"

type CanvasClient struct {
	client      *http.Client
	apiPath     *url.URL
	accessToken string
}

func NewClient(client *http.Client, rawUrl string, accessToken string) *CanvasClient {
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
	apiPath := url.URL{
		Scheme: "https",
		Host:   canvasHost,
		Path:   apiPath,
	}
	return &CanvasClient{
		client:      client,
		accessToken: accessToken,
		apiPath:     &apiPath,
	}
}

func (c *CanvasClient) ExtractDomainBrowserCookies() {
	baseUrl := url.URL{Scheme: c.apiPath.Scheme, Host: c.apiPath.Host}
	cookieJar := utils.ExtractCanvasBrowserCookies(baseUrl.String())
	c.client.Jar = cookieJar
}

func (c *CanvasClient) GetActiveEnrolledCourses() []nodes.CourseNode {
	coursesUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + "/users/self/courses",
		RawQuery: url.Values{
			"enrollment_state": {"active"},
		}.Encode(),
	}

	// courses request
	req, err := http.NewRequest("GET", coursesUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		panic(err)
	}

	resp, err := c.client.Do(req)
	courseJson := utils.ExtractResponseToString(resp)
	var courses []nodes.CourseNode
	json.Unmarshal([]byte(courseJson), &courses)
	if strings.Contains(courseJson, "user authorisation required") {
		panic(fmt.Errorf("\nerror querying '/api/v1/users/self/courses?enrollment_state=active' endpoint: %v\nTRY LAUNCHING https://canvas.nus.edu.sg in a chrome/safari/edge browser and try again!", courseJson))
	}

	return courses
}

func (c *CanvasClient) getCourseUrl(id int) url.URL {
	courseUrl := url.URL{
		Scheme: c.apiPath.Scheme,
		Host:   c.apiPath.Host,
		Path:   c.apiPath.Path + fmt.Sprintf("/courses/%d/folders/root", id),
	}
	return courseUrl
}

func (c *CanvasClient) GetCourseRootFolder(courseId int) nodes.DirectoryNode {
	courseUrl := c.getCourseUrl(courseId)
	req, err := http.NewRequest("GET", courseUrl.String(), nil)
	utils.SetQueryAccessToken(req, c.accessToken)
	if err != nil {
		panic(err)
	}
	rootRes, err := c.client.Do(req)
	if err != nil {
		log.Fatalf("failed to query course %d root directory: %s", courseId, err.Error())
	}
	rootJson := utils.ExtractResponseToString(rootRes)
	var rootNode nodes.DirectoryNode
	if err := json.Unmarshal([]byte(rootJson), &rootNode); err != nil {
		log.Fatalf("failed to extract course %d root directory: %s", courseId, err.Error())
	}
	return rootNode
}

func (c *CanvasClient) RecurseDirectoryNode(node *nodes.DirectoryNode, parent *nodes.DirectoryNode) {
	dir := ""
	if parent != nil {
		dir += parent.Directory
	}
	if node == nil {
		panic("node is nil!")
	}
	dir += fmt.Sprintf("%s/", node.Name)
	node.Directory = dir

	if node.FilesCount > 0 {
		fileReq, err := http.NewRequest("GET", node.FilesUrl, nil)
		if err != nil {
			panic(err)
		}
		utils.SetQueryAccessToken(fileReq, c.accessToken)
		filesRes, err := c.client.Do(fileReq)
		if err != nil {
			panic(err)
		}
		filesJson := utils.ExtractResponseToString(filesRes)
		var files []*nodes.FileNode
		json.Unmarshal([]byte(filesJson), &files)
		for f := range files {
			files[f].Directory = fmt.Sprintf("%s%s", dir, files[f].Display_name)
		}
		node.FileNodes = files
	}

	if node.FoldersCount > 0 {
		folderReq, err := http.NewRequest("GET", node.FoldersUrl, nil)
		if err != nil {
			panic(err)
		}
		utils.SetQueryAccessToken(folderReq, c.accessToken)
		foldersRes, err := c.client.Do(folderReq)
		if err != nil {
			panic(err)
		}
		foldersJson := utils.ExtractResponseToString(foldersRes)
		var folders []*nodes.DirectoryNode
		json.Unmarshal([]byte(foldersJson), &folders)
		for fi := range folders {
			c.RecurseDirectoryNode(folders[fi], node)
		}
		node.FolderNodes = folders
	}
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

func (c *CanvasClient) RecursiveCreateNode(node *nodes.DirectoryNode, updateNumDownloads func(numDownloads int)) {
	if node == nil {
		return
	}
	if err := os.MkdirAll(node.Directory, 755); err != nil {
		panic(err)
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
			c.downloadFileNode(node.FileNodes[i])
		}(j)
	}
	updateNumDownloads(numDownloads)
	for d := range node.FolderNodes {
		c.RecursiveCreateNode(node.FolderNodes[d], updateNumDownloads)
	}
	wg.Wait()
}

func (c *CanvasClient) RecursiveUpdateNode(node *nodes.DirectoryNode, updateNumDownloads func(numDownloads int)) {
	if node == nil {
		return
	}
	// create directory if doesnt exist
	if _, err := os.Stat(node.Directory); os.IsNotExist(err) {
		if err := os.MkdirAll(node.Directory, 755); err != nil {
			panic(err)
		}
	}
	numDownloads := 0
	for j := range node.FileNodes {
		if node.FileNodes[j] == nil {
			continue
		}
		file, err := os.Stat(node.FileNodes[j].Directory)
		if err != nil {
			numDownloads += 1
			c.downloadFileNode(node.FileNodes[j])
		} else {
			if file.ModTime().Unix() < node.FileNodes[j].UpdatedAt.Unix() {
				numDownloads += 1
				c.downloadFileNode(node.FileNodes[j])
			}
		}
	}
	updateNumDownloads(numDownloads)
	for d := range node.FolderNodes {
		c.RecursiveUpdateNode(node.FolderNodes[d], updateNumDownloads)
	}
}
