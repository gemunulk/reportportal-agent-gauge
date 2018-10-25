package reportportal

import "time"

// Log represents a Project Logs
type Log struct {
	ID string `json:"id"`
}

// File structue
type File struct {
	Name string `json:"name"`
}

// CreateLogRequest structure
type CreateLogRequest struct {
	File struct {
		Name string `json:"name"`
	} `json:"file"`
	ItemID  string `json:"item_id"`
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

// CreateLog creates a log in Report Portal under the item specified.
func (c *Client) CreateLog(projectName string, fileName string, itemID string, logLevel string, message string) (Log, error) {
	url := projectName + "/log"
	timeNow := time.Now().Format(time.RFC3339)

	jsonStr := CreateLogRequest{
		File: File{
			Name: fileName,
		},
		ItemID:  itemID,
		Level:   logLevel,
		Message: message,
		Time:    timeNow,
	}
	returnLog := Log{}
	err := c.sendRequest("POST", url, jsonStr, &returnLog)
	return returnLog, err

}
