package models

import (
	"errors"
	"fmt"
)

type Tag struct {
	TagID string `json:"tagID,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

var (
	checkTagQuery  = `select tag_id, tag from tags where tag = $1`
	insertTagQuery = `insert into tags (tag) values ($1) returning tag_id, tag`
)

func (db DB) createTag(blurbId, tag string) error {
	var tagObj Tag
	row := db.QueryRow(insertTagQuery, tag)
	if err := row.Scan(&tagObj.TagID, &tagObj.Tag); err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	fmt.Println(blurbId, tagObj.TagID)
	_, err := db.Exec(`insert into blurbs_tags (blurb_id, tag_id) values ($1, $2)`, blurbId, tagObj.TagID)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	return nil
}

func (db DB) CheckTag(blurbId string, tag string) error {
	if tag == "" {
		return errors.New("no tag specified")
	}
	row := db.QueryRow(checkTagQuery, tag)
	var tagObj Tag
	if err := row.Scan(&tagObj.TagID, &tagObj.Tag); err != nil {
		return db.createTag(blurbId, tag)
	}
	if tagObj.TagID != "" {
		return nil
	}
	_, err := db.Exec(`insert into blurbs_tags (blurb_id, tag_id) values ($1, $2)`, blurbId, tagObj.TagID)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	return nil
}
