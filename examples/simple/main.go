package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/miku/clam"
	"github.com/miku/lynd"
)

type Download struct {
	Link string
}

func (task Download) Requires() []lynd.Task { return nil }
func (task Download) Run() error {
	output, err := clam.RunOutput(`curl "{{ link }}" > {{ output }}`, clam.Map{"link": task.Link})
	if err != nil {
		return err
	}
	return os.Rename(output, task.Output().(lynd.LocalTarget).Path)
}
func (task Download) Output() lynd.Target {
	enc := base64.StdEncoding.EncodeToString([]byte(task.Link))
	return lynd.LocalTarget{Path: fmt.Sprintf("/tmp/lynd-Download-%s", enc)}
}

type Concat struct {
	Links []string
}

func (task Concat) Requires() []lynd.Task {
	var tasks []lynd.Task
	for _, link := range task.Links {
		tasks = append(tasks, Download{Link: link})
	}
	return tasks
}

func (task Concat) Run() error {
	output, err := clam.RunOutput("cat {{ files }} > {{ output }}",
		clam.Map{"files": strings.Join(lynd.LocalInputPaths(task), " ")})
	if err != nil {
		return err
	}
	return os.Rename(output, lynd.LocalOutputPath(task))
}

func (task Concat) Output() lynd.Target {
	return lynd.LocalTarget{Path: "/tmp/lynd-Download-2"}
}

func main() {
	err := lynd.Build(Concat{Links: []string{"http://www.heise.de", "http://www.faz.net", "http://www.bild.de"}})
	if err != nil {
		log.Fatal(err)
	}
}
