// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUserUserGroupRelation = "user_user_group_relation"

// UserUserGroupRelation mapped from table <user_user_group_relation>
type UserUserGroupRelation struct {
	UserID      uint64 `gorm:"column:user_id;primaryKey;comment:auto increment id" json:"userId,string"`
	UserGroupID uint64 `gorm:"column:user_group_id;primaryKey;comment:auto increment id" json:"userGroupId,string"`
}

// TableName UserUserGroupRelation's table name
func (*UserUserGroupRelation) TableName() string {
	return TableNameUserUserGroupRelation
}
