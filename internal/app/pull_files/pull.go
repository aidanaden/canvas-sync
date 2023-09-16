package pull_files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/aidanaden/canvas-sync/internal/pkg"
	"github.com/chelnak/ysmrr"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_ "github.com/zellyn/kooky/browser/all" // register cookie store finders!
)

func recurseDirectoryNode(client *http.Client, node *DirectoryNode, parent *DirectoryNode, accessToken string) {
	dir := ""
	if parent != nil {
		dir += parent.Directory
	}
	if node == nil {
		panic("node is nil!")
	}
	dir += fmt.Sprintf("%s/", node.Name)
	node.Directory = dir

	// // increment methods won't trigger complete event because bar was constructed with total = 0
	// bar.IncrBy(rand.Intn(1024) + 1)
	// // following call is not required, it's called to show some progress instead of an empty bar
	// bar.SetTotal(bar.Current()+2048, false)
	// // time.Sleep(time.Duration(rand.Intn(10)+1) * maxSleep / 10)

	if node.FilesCount > 0 {
		fileReq, err := http.NewRequest("GET", node.FilesUrl, nil)
		if err != nil {
			panic(err)
		}
		setQueryAccessToken(fileReq, accessToken)
		filesRes, err := client.Do(fileReq)
		if err != nil {
			panic(err)
		}
		filesJson := pkg.ExtractResponseToString(filesRes)
		var files []*FileNode
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
		setQueryAccessToken(folderReq, accessToken)
		foldersRes, err := client.Do(folderReq)
		if err != nil {
			panic(err)
		}
		foldersJson := pkg.ExtractResponseToString(foldersRes)
		var folders []*DirectoryNode
		json.Unmarshal([]byte(foldersJson), &folders)
		for fi := range folders {
			recurseDirectoryNode(client, folders[fi], node, accessToken)
		}
		node.FolderNodes = folders
	}
}

func downloadFileNode(client *http.Client, node *FileNode) error {
	if client == nil || node == nil {
		return errors.New("cannot download file without http client or file node")
	}
	file, err := os.Create(node.Directory)
	if err != nil {
		return err
	}
	defer file.Close()
	res, err := client.Get(node.Url)
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

func recursiveCreateNode(client *http.Client, node *DirectoryNode) {
	if node == nil {
		return
	}
	if err := os.MkdirAll(node.Directory, os.ModePerm); err != nil {
		panic(err)
	}
	for j := range node.FileNodes {
		if node.FileNodes[j] == nil {
			continue
		}
		downloadFileNode(client, node.FileNodes[j])
	}
	for d := range node.FolderNodes {
		recursiveCreateNode(client, node.FolderNodes[d])
	}
}

func setQueryAccessToken(req *http.Request, accessToken string) {
	if accessToken == "" {
		return
	}
	q := req.URL.Query()
	q.Set("access_token", accessToken)
	req.URL.RawQuery = q.Encode()
}

func getCourseCodesFromArgs(rawArgs []string) []string {
	args := make([]string, 0)
	for _, raw := range rawArgs {
		if strings.Contains(raw, ",") {
			splits := strings.Split(raw, ",")
			for _, split := range splits {
				args = append(args, strings.ToLower(split))
			}
		} else {
			args = append(args, strings.ToLower(raw))
		}
	}
	return args
}

func Run(cmd *cobra.Command, args []string) {
	targetDir := fmt.Sprintf("%v/files", viper.Get("data_dir"))
	accessToken := fmt.Sprintf("%v", viper.Get("access_token"))

	fmt.Printf("\nfiles will be downloaded to data_dir: %s", targetDir)

	client := http.Client{}
	if accessToken == "" {
		fmt.Printf("\nno cookies found, getting auth cookies from browser...")

		var rawCookies []*http.Cookie
		rawCookies = pkg.ExtractCookies("canvas.nus.edu.sg")
		if rawCookies == nil || len(rawCookies) == 0 {
			if err := browser.OpenURL("https://canvas.nus.edu.sg"); err != nil {
				panic(err)
			}
			for {
				rawCookies = pkg.ExtractCookies("canvas.nus.edu.sg")
				if rawCookies != nil && len(rawCookies) > 0 {
					break
				}
			}
		}
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalf("Got error while creating cookie jar %s", err.Error())
		}
		client.Jar = jar
		url, err := url.Parse("https://canvas.nus.edu.sg")
		if err != nil {
			panic(err)
		}
		client.Jar.SetCookies(url, rawCookies)
	}

	fmt.Printf("\naccess token found: %s", accessToken)

	providedCodes := getCourseCodesFromArgs(args)
	// courses request
	req, err := http.NewRequest("GET", "https://canvas.nus.edu.sg/api/v1/courses", nil)
	setQueryAccessToken(req, accessToken)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	courseJson := pkg.ExtractResponseToString(resp)
	var rawCourses []CourseNode
	json.Unmarshal([]byte(courseJson), &rawCourses)

	if strings.Contains(courseJson, "user authorisation required") {
		panic(fmt.Errorf("\nerror querying '/api/v1/courses' endpoint: %v\nTRY LAUNCHING https://canvas.nus.edu.sg in a chrome/safari/edge browser and try again!", courseJson))
	}
	courses := make([]CourseNode, 0)
	for _, raw := range rawCourses {
		if raw.CourseCode == "" {
			continue
		}
		// add any course if no code provided
		if len(providedCodes) == 0 {
			courses = append(courses, raw)
			continue
		}
		// add course if code matches provided codes
		for _, provided := range providedCodes {
			if strings.ToLower(raw.CourseCode) == provided {
				courses = append(courses, raw)
			}
		}
	}

	var wg sync.WaitGroup
	sm := ysmrr.NewSpinnerManager()

	for ci := range courses {
		wg.Add(1)
		sp := sm.AddSpinner(fmt.Sprintf("Starting files download for %s...", courses[ci].CourseCode))
		go func(i int, sp *ysmrr.Spinner) {
			defer wg.Done()
			code := courses[i].CourseCode

			req, err := http.NewRequest("GET", fmt.Sprintf("https://canvas.nus.edu.sg/api/v1/courses/%d/folders/root", courses[i].ID), nil)
			setQueryAccessToken(req, accessToken)
			if err != nil {
				panic(err)
			}
			rootRes, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error occured. Error is: %s", err.Error())
			}
			rootJson := pkg.ExtractResponseToString(rootRes)
			var rootNode DirectoryNode
			json.Unmarshal([]byte(rootJson), &rootNode)
			rootNode.Name = fmt.Sprintf("%s/%s", targetDir, code)

			sp.UpdateMessagef("Pulling files info for %s", code)
			recurseDirectoryNode(&client, &rootNode, nil, accessToken)
			sp.UpdateMessagef("Downloading files for %s", code)
			recursiveCreateNode(&client, &rootNode)
			sp.UpdateMessagef("Downloaded files for %s", code)
			sp.Complete()
		}(ci, sp)
	}
	sm.Start()
	wg.Wait()
	sm.Stop()
	fmt.Printf("\ndownloaded files - view here: %s", targetDir)
}
