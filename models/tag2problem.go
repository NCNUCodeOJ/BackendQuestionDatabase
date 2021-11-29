package models

import "gorm.io/gorm"

// Tag2Problem Database
type Tag2Problem struct {
	gorm.Model
	TagName   string `gorm:"NOT NULL"`
	ProblemID uint   `gorm:"NOT NULL"`
}

//AddTag2Problem - problem have one tag
func AddTag2Problem(problemID uint, tagName string) (err error) {
	tag, notFound := GetTagByName(tagName)
	tag2table := Tag2Problem{}

	if notFound != nil {
		tag.Name = tagName
		addTag(&tag)
	}

	tag2table.TagName = tagName
	tag2table.ProblemID = problemID
	err = DB.Create(&tag2table).Error
	return
}

//DeleteProblemAllTags 用 problem id 刪除 所有與這個problem id 有關的 row
func DeleteProblemAllTags(problemID uint) (err error) {
	err = DB.Where(&Tag2Problem{ProblemID: problemID}).Delete(Tag2Problem{}).Error
	return
}

// DeleteProblemTag delete problem one tag
func DeleteProblemTag(problemID uint, tagName string) (err error) {
	err = DB.Where(&Tag2Problem{ProblemID: problemID, TagName: tagName}).Delete(&Tag2Problem{}).Error
	return
}

// GetProblemAllTags 查詢 problem 所有 tag
func GetProblemAllTags(ProblemID uint) (tags []string, err error) {
	var tag2problems []Tag2Problem
	if err = DB.Where(&Tag2Problem{ProblemID: ProblemID}).Find(&tag2problems).Error; err != nil {
		return
	}
	for _, tag2problem := range tag2problems {
		tags = append(tags, tag2problem.TagName)
	}
	return
}

// GetProblemsByTag 查詢 該 tag 所有 problems
func GetProblemsByTag(TagName string) (problems []Problem, err error) {
	var tag2problems []Tag2Problem
	if err = DB.Where(&Tag2Problem{TagName: TagName}).Find(&tag2problems).Error; err != nil {
		return
	}
	for _, tag2problem := range tag2problems {
		var problem Problem

		if problem, err = GetProblemByID(tag2problem.ProblemID); err != nil {
			return
		}
		problems = append(problems, problem)
	}
	return
}
