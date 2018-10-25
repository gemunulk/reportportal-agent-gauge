package reportportal

import "time"

// TestItem represents a Project TestItem
type TestItem struct {
	ID string `json:"id"`
}

// TestFilter represents Project Test filters
type TestFilter struct {
	Content []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"parameters"`
		Tags       []string `json:"tags"`
		Type       string   `json:"type"`
		StartTime  int64    `json:"start_time"`
		Status     string   `json:"status"`
		Statistics struct {
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
		Parent    string `json:"parent"`
		PathNames struct {
			FiveAe46Ee136D1A000013E2Fe9 string `json:"5ae46ee136d1a000013e2fe9"`
		} `json:"path_names"`
		HasChilds bool   `json:"has_childs"`
		LaunchID  string `json:"launchId"`
	} `json:"content"`
	Page struct {
		Number        int `json:"number"`
		Size          int `json:"size"`
		TotalElements int `json:"totalElements"`
		TotalPages    int `json:"totalPages"`
	} `json:"page"`
}

//RootAndChildTestItem request structure
type RootAndChildTestItem struct {
	Description string `json:"description"`
	LaunchID    string `json:"launch_id"`
	Name        string `json:"name"`
	Parameters  []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"parameters"`
	StartTime string   `json:"start_time"`
	Tags      []string `json:"tags"`
	Type      string   `json:"type"`
	UniqueID  string   `json:"uniqueId"`
}

// Parameters param structure
type Parameters []struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Issue param structure
type Issue struct {
	Comment              string `json:"comment"`
	ExternalSystemIssues []struct {
		SubmitDate int    `json:"submitDate"`
		Submitter  string `json:"submitter"`
		SystemID   string `json:"systemId"`
		TicketID   string `json:"ticketId"`
		URL        string `json:"url"`
	} `json:"externalSystemIssues"`
	IssueType string `json:"issue_type"`
}

// ExternalSystemIssues param structure
type ExternalSystemIssues []struct {
	SubmitDate int    `json:"submitDate"`
	Submitter  string `json:"submitter"`
	SystemID   string `json:"systemId"`
	TicketID   string `json:"ticketId"`
	URL        string `json:"url"`
}

// FinishTest param structure
type FinishTest struct {
	Description string `json:"description"`
	EndTime     string `json:"end_time"`
	Issue       struct {
		Comment              string `json:"comment"`
		ExternalSystemIssues []struct {
			SubmitDate int    `json:"submitDate"`
			Submitter  string `json:"submitter"`
			SystemID   string `json:"systemId"`
			TicketID   string `json:"ticketId"`
			URL        string `json:"url"`
		} `json:"externalSystemIssues"`
		IssueType string `json:"issue_type"`
	} `json:"issue"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
}

// FindTestItems request structure
type FindTestItems struct {
}

// StartARootTestItem starts specified Root test item.
func (c *Client) StartARootTestItem(projectName string, launchID string, testItemName string, tags []string, testType string) (TestItem, error) {
	url := projectName + "/item"
	timeNow := time.Now().Format(time.RFC3339)

	jsonStr := RootAndChildTestItem{
		Description: testItemName,
		LaunchID:    launchID,
		Name:        testItemName,
		Parameters: Parameters{
			{Key: "string",
				Value: "string"},
		},
		StartTime: timeNow,
		Tags:      tags,
		Type:      testType,
		UniqueID:  "string",
	}

	returnTestItem := TestItem{}
	err := c.sendRequest("POST", url, jsonStr, &returnTestItem)
	return returnTestItem, err

}

// FinishTestItem finishes the Launch specified.
func (c *Client) FinishTestItem(projectName string, testItemID string, testItemName string, status string, tags []string) (TestItem, error) {
	url := projectName + "/item/" + testItemID
	timeNow := time.Now().Format(time.RFC3339)

	jsonStr := FinishTest{
		Description: testItemName,
		EndTime:     timeNow,
		Issue: Issue{
			Comment:   testItemName,
			IssueType: "TO_INVESTIGATE",
		},
		Status: status,
		Tags:   tags,
	}

	returnTestItem := TestItem{}
	err := c.sendRequest("PUT", url, jsonStr, &returnTestItem)
	return returnTestItem, err

}

// StartAchildTestItem starts specified Child test item.
func (c *Client) StartAchildTestItem(projectName string, launchID string, parentItem string, testItemName string, tags []string, testType string) (TestItem, error) {
	url := projectName + "/item/" + parentItem
	timeNow := time.Now().Format(time.RFC3339)

	jsonStr := RootAndChildTestItem{
		Description: testItemName,
		LaunchID:    launchID,
		Name:        testItemName,
		Parameters: Parameters{
			{Key: "string",
				Value: "string"},
		},
		StartTime: timeNow,
		Tags:      tags,
		Type:      testType,
		UniqueID:  "string",
	}
	returnTestItem := TestItem{}
	err := c.sendRequest("POST", url, jsonStr, &returnTestItem)
	return returnTestItem, err

}

// FindTestItemsBySpecifiedFilter specific tests by filter
func (c *Client) FindTestItemsBySpecifiedFilter(projectName string, filterQuery string) (TestFilter, error) {
	url := projectName + "/item?" + filterQuery
	jsonStr := FindTestItems{}

	returnTestFilter := TestFilter{}
	err := c.sendRequest("GET", url, jsonStr, &returnTestFilter)
	return returnTestFilter, err
}
