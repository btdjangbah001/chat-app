package models

type Group struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	OwnerID     uint   `json:"owner_id"`
	CreatedAt   string `json:"created_at"`
}

type CreateGroup struct {
	Name        string `json:"name"`
	OwnerID	 uint   `json:"owner_id"`
}

type UpdateGroupDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (group *Group) CreateGroup() (err error) {
	err = DB.Create(&group).Error
	if err != nil {
		return err
	}
	return nil
}

func (group *Group) GetGroup() (err error) {
	err = DB.Where("id = ?", group.ID).First(&group).Error
	if err != nil {
		return err
	}
	return nil
}

func getGroups(groupIDs *[]uint) (*[]Group, error) {
	var groups []Group
	err := DB.Where("id IN (?)", groupIDs).Find(&groups).Error
	if err != nil {
		return nil, err
	}
	return &groups, nil
}

func GetGroupsForTwoUsers(user1ID uint, user2ID uint) (*[]Group, error) {
	groupIDs, err := CheckCommonMembershipForTwoUsers(user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	groups, err := getGroups(groupIDs)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// func (group *Group) UpdateGroupDescription() (err error) {
// 	err = DB.Model(&group).Where("id = ?", group.ID).Updates(&group).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }