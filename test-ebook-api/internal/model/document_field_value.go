package model

// DocumentFieldValue 文档动态属性值
type DocumentFieldValue struct {
    ID         uint   `gorm:"primarykey" json:"id"`
    DocumentID uint   `json:"document_id" gorm:"index;not null"`
    FieldID    uint   `json:"field_id" gorm:"index:idx_field_value;not null"`
    Value      string `json:"value" gorm:"index:idx_field_value,length:20;type:varchar(255)"`
    // 关联
    Field FormField `json:"field" gorm:"foreignKey:FieldID"`
}
