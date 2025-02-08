// internal/entity/user.go
package entity

import "time"

type User struct {
	ID                    int       `json:"id"`
	LINEUserID            string    `json:"lineUserId"`
	SelectedAreaID        string    `json:"selectedAreaId"`
	SelectedAreaOfficeID  string    `json:"selectedAreaOfficeId"`  // 追加：都道府県(=area_office)選択結果
	SelectedAreaClass15ID string    `json:"selectedAreaClass15Id"` // 追加：area_class15選択結果
	NotifyTime            time.Time `json:"notifyTime"`
	IsActive              bool      `json:"isActive"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
	Status                string    `json:"status"` // 例: "awaiting_prefecture", "awaiting_municipality", "awaiting_confirmation", "awaiting_area_class10_selection", "awaiting_area_class15_selection", "awaiting_area_class20_selection", "completed"
}
