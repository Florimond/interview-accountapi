package contracts

import "time"

// RecordInfo holds the basic info about any record returned by the API
type RecordInfo struct {
	ID             string    `json:"id"`
	CreatedOn      time.Time `json:"created_on"`
	ModifiedOn     time.Time `json:"modified_on"`
	OrganisationID string    `json:"organiation_id"`
	Type           string    `json:"type"`
	Version        uint      `json:"version"`
}
