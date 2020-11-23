//+build ignore

package main

import (
	"fmt"
	"os"

	"github.com/AlekSi/pointer"
	"github.com/xanzy/go-gitlab"
)

const URL = "http://gitlab.com/api/v4"

var tokenAuth = os.Getenv("GITLAB-TOKEN")

func main() {
	client, err := gitlab.NewClient(tokenAuth, gitlab.WithBaseURL(URL))
	if err != nil {
		panic(err)
	}
	//projects, _, err := client.Projects.ListProjects(&gitlab.ListProjectsOptions{})
	//if err != nil {
	//	panic(err)
	//}

	issues, _, err := client.Issues.ListIssues(&gitlab.ListIssuesOptions{
		Scope: pointer.ToString("all"),
	})
	for _, iss := range issues {
		fmt.Println(iss.ID, iss.Title)
	}

	//for _, project := range projects {
	//	fmt.Println("Project " + project.Name)
	//
	//	panic(err)
	//}
}
