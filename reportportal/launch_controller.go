package reportportal

import "time"

// Launch represents a Project Launch
type Launch struct {
	ID      string `json:"id"`
	MSG     string `json:"msg"`
	CONTENT string `json:"content"`
}

// LaunchFilter test
type LaunchFilter struct {
	Content []struct {
		Owner       string `json:"owner"`
		Share       bool   `json:"share"`
		Description string `json:"description"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		Number      int    `json:"number"`
		StartTime   int64  `json:"start_time"`
		Status      string `json:"status"`
		Statistics  struct {
			Executions struct {
				Total   string `json:"total"`
				Passed  string `json:"passed"`
				Failed  string `json:"failed"`
				Skipped string `json:"skipped"`
			} `json:"executions"`
			Defects struct {
				ProductBug struct {
					Total int `json:"total"`
					PB001 int `json:"PB001"`
				} `json:"product_bug"`
				AutomationBug struct {
					AB001 int `json:"AB001"`
					Total int `json:"total"`
				} `json:"automation_bug"`
				SystemIssue struct {
					Total int `json:"total"`
					SI001 int `json:"SI001"`
				} `json:"system_issue"`
				ToInvestigate struct {
					Total int `json:"total"`
					TI001 int `json:"TI001"`
				} `json:"to_investigate"`
				NoDefect struct {
					ND001 int `json:"ND001"`
					Total int `json:"total"`
				} `json:"no_defect"`
			} `json:"defects"`
		} `json:"statistics"`
		Tags                []string `json:"tags"`
		Mode                string   `json:"mode"`
		IsProcessing        bool     `json:"isProcessing"`
		ApproximateDuration float32  `json:"approximateDuration"`
	} `json:"content"`
	Page struct {
		Number        int `json:"number"`
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
	} `json:"page"`
}

// StartsLaunch request structure
type StartsLaunch struct {
	Description string   `json:"description"`
	Mode        string   `json:"mode"`
	Name        string   `json:"name"`
	StartTime   string   `json:"start_time"`
	Tags        []string `json:"tags"`
}

// FinishLaunch request structure
type FinishLaunch struct {
	Description string   `json:"description"`
	EndTime     string   `json:"end_time"`
	Status      string   `json:"status"`
	Tags        []string `json:"tags"`
}

// GetListOfProjectLaunches request structure
type GetListOfProjectLaunches struct {
}

// StartsLaunchForSpecifiedProject returns the existing project LaunchID
func (c *Client) StartsLaunchForSpecifiedProject(projectName string, launchName string, tags []string) (Launch, error) {
	url := projectName + "/launch"
	timeNow := time.Now().Format(time.RFC3339)

	jsonStr := StartsLaunch{
		Description: launchName,
		Mode:        "DEFAULT",
		Name:        launchName,
		StartTime:   timeNow,
		Tags:        tags,
	}

	returnLaunch := Launch{}
	err := c.sendRequest("POST", url, jsonStr, &returnLaunch)
	return returnLaunch, err

}

// FinishLaunchForSpecifiedProject finishes the Launch specified.
func (c *Client) FinishLaunchForSpecifiedProject(projectName string, launchID string, status string, tags []string) (Launch, error) {
	url := projectName + "/launch/" + launchID + "/finish"
	timeNow := time.Now().Format(time.RFC3339)

	jsonStr := FinishLaunch{
		Description: launchID,
		EndTime:     timeNow,
		Status:      status,
		Tags:        tags,
	}

	returnLaunch := Launch{}
	err := c.sendRequest("PUT", url, jsonStr, &returnLaunch)
	return returnLaunch, err

}

// GetListOfProjectLaunchesByFilter gets project specific lauches with the filter
func (c *Client) GetListOfProjectLaunchesByFilter(projectName string, filterQuery string) (LaunchFilter, error) {
	url := projectName + "/launch?" + filterQuery
	jsonStr := GetListOfProjectLaunches{}

	returnLaunchFilter := LaunchFilter{}
	err := c.sendRequest("GET", url, jsonStr, &returnLaunchFilter)
	return returnLaunchFilter, err
}
