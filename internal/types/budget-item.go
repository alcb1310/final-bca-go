package types

import (
	"database/sql"

	"github.com/google/uuid"
)

type BudgetItem struct {
	Id         uuid.UUID      `json:"id"`
	Code       string         `json:"code"`
	Name       string         `json:"name"`
	Level      uint8          `json:"level"`
	Accumulate bool           `json:"accumulate"`
	ParentId   uuid.NullUUID  `json:"parentId"`
	ParentCode sql.NullString `json:"parentCode"`
	ParentName sql.NullString `json:"parentName"`
}

type CreateBudgetItem struct {
	Code       string        `json:"code"`
	Name       string        `json:"name"`
	Accumulate bool          `json:"accumulate"`
	ParentId   uuid.NullUUID `json:"parentId"`
}
