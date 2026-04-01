package model

import "gorm.io/gorm"

// Form 动态表单大类模型
type Form struct {
	gorm.Model
	Name        string      `json:"name" gorm:"type:varchar(100);not null"`
	Description string      `json:"description" gorm:"type:varchar(255)"`
	Fields      []FormField `json:"fields" gorm:"foreignKey:FormID;constraint:OnDelete:CASCADE;"`
}

// FormField 动态表单属性项模型
type FormField struct {
	gorm.Model
	FormID      uint   `json:"form_id" gorm:"index"`
	Label       string `json:"label" gorm:"type:varchar(100);not null"` // 如：签约日期
	FieldKey    string `json:"field_key" gorm:"type:varchar(50);not null"` // 如：sign_date
	FieldType   string `json:"field_type" gorm:"type:varchar(50);not null"` // input, select, date, number
	IsRequired  bool   `json:"is_required" gorm:"default:false"`
	Options     string `json:"options" gorm:"type:text"` // JSON 数组，用于 select 选项
	Order        int    `json:"order" gorm:"default:0"`
	ShowInList   bool   `json:"show_in_list" gorm:"default:true"`
	ShowInFilter bool   `json:"show_in_filter" gorm:"default:false"`
	DefaultValue string `json:"default_value" gorm:"type:varchar(255)"`
}
