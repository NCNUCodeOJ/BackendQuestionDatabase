package models

import "gorm.io/gorm"

// Tag2Problem Database
type Tag2Problem struct {
	gorm.Model
	TagID     uint `gorm:"NOT NULL"`
	ProblemID uint `gorm:"NOT NULL"`
}

//AddTag2Problem - problem have one tag
func AddTag2Problem(tagName string, problemID uint) (err error) {
	tag, notFound := GetTagByName(tagName)
	tag2table := Tag2Problem{}

	if notFound != nil {
		tag.Name = tagName
		addTag(&tag)
	}

	tag2table.TagID = tag.ID
	tag2table.ProblemID = problemID
	err = DB.Create(&tag2table).Error
	return
}

//DeleteProblemAllTags 用 problem id 刪除 所有與這個problem id 有關的 row
func DeleteProblemAllTags(problemID uint) (err error) {
	err = DB.Where(&Tag2Problem{ProblemID: problemID}).Delete(Tag2Problem{}).Error
	return
}

// GetProblemAllTags 查詢 problem 所有 tag
func GetProblemAllTags(ProblemID uint) (tags []Tag, err error) {
	var tag2problems []Tag2Problem
	if err = DB.Where(&Tag2Problem{ProblemID: ProblemID}).Find(&tag2problems).Error; err != nil {
		return
	}
	for _, tag2problem := range tag2problems {
		var tag Tag
		if tag, err = getTagByID(tag2problem.TagID); err != nil {
			return
		}
		tags = append(tags, tag)
	}
	return
}
