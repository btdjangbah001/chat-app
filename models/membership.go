package models

type Membership struct {
	UserID  uint `json:"user_id"`
	GroupID uint `json:"group_id"`
}

func (membership *Membership) CreateMembership() (err error) {
	err = DB.Create(&membership).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckMembership(userID uint, groupID uint) (bool, error) {
	var membership Membership
	err := DB.Where("user_id = ? AND group_id = ?", userID, groupID).First(&membership).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (membership *Membership) DeleteMembership() (err error) {
	err = DB.Where("user_id = ? AND group_id = ?", membership.UserID, membership.GroupID).Delete(&membership).Error
	if err != nil {
		return err
	}
	return nil
}

func CheckCommonMembershipForTwoUsers(user1ID uint, user2ID uint) (*[]uint, error){
	var memberships []Membership
	var groupIDs []uint
	err := DB.Where("user_id = ?", user1ID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	
	var memberships2 []Membership
	err = DB.Where("user_id = ?", user2ID).Find(&memberships2).Error
	if err != nil {
		return nil, err
	}
	
	for _, membership := range memberships {
		for _, membership2 := range memberships2 {
			if membership.GroupID == membership2.GroupID {
				groupIDs = append(groupIDs, membership.GroupID)
			}
		}
	}
	return &groupIDs, nil
}

func GetMembershipsForUser(userID uint) (*[]uint, error) {
	var memberships []Membership
	err := DB.Where("user_id = ?", userID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	var groupIDs []uint
	for _, membership := range memberships {
		groupIDs = append(groupIDs, membership.GroupID)
	}
	return &groupIDs, nil
}

func GetGroupParticipants(groupID uint) (*[]uint, error) {
	var memberships []Membership
	err := DB.Where("group_id = ?", groupID).Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	var userIDs []uint
	for _, membership := range memberships {
		userIDs = append(userIDs, membership.UserID)
	}
	return &userIDs, nil
}