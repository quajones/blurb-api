package models

import "time"

type Blurb struct {
	ID               string    `uri:"id" json:"id,omitempty"`
	UserID           string    `json:"userID,omitempty"`
	Topic            string    `json:"topic,omitempty`
	Title            string    `json:"title,omitempty"`
	Content          string    `json:"content,omitempty"`
	CreatedDate      time.Time `json:"createdDate,omitempty"`
	PublishedDate    time.Time `json:"publishedDate,omitempty"`
	LastModifiedDate time.Time `json:"lastModifiedDate,omitempty"`
}

var (
	getBlurbsForFollowingQuery = `select a.blurb_id, a.user_id, a.topic, a.title, a.content, a.created_date, a.last_modified 
									from blurbs a join user_following b on  a.user_id = b.following_id where b.user_id = $1
									order by a.last_modified desc`
)

// create a blurb given a user and topic
func (db DB) CreateBlurb(userId string, topic int, title, content string) error {
	if userId == "" || title == "" || content == "" {
		return errInvalidInputData
	}
	_, err := db.Exec(`insert into blurbs (user_id, topic, title, content, created_date, last_modified) values ($1, $2, $3, $4, $5, $6)`,
		userId, topic, title, content, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
func (db DB) GetAllBlurbsForUser(userId string) ([]Blurb, error) {
	if userId == "" {
		return nil, errInvalidInputData
	}
	rows, err := db.Query(`select a.blurb_id, a.user_id, b.topic, a.title, a.content, a.created_date, a.last_modified from blurbs a join topics b on a.topic = b.topic_id where a.user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	var blurbs []Blurb
	for rows.Next() {
		var blurb Blurb
		if err := rows.Scan(&blurb.ID, &blurb.UserID, &blurb.Topic, &blurb.Title, &blurb.Content, &blurb.CreatedDate, &blurb.LastModifiedDate); err != nil {
			return nil, err
		}
		blurbs = append(blurbs, blurb)
	}
	return blurbs, nil
}

func (db DB) GetAllBlurbsForFollowing(userId string) ([]Blurb, error) {
	rows, err := db.Query(getBlurbsForFollowingQuery, userId)
	if err != nil {
		return nil, err
	}
	var blurbs []Blurb
	for rows.Next() {
		var blurb Blurb
		if err := rows.Scan(&blurb.ID, &blurb.UserID, &blurb.Topic, &blurb.Title, &blurb.Content, &blurb.CreatedDate, &blurb.LastModifiedDate); err != nil {
			return nil, err
		}
		blurbs = append(blurbs, blurb)
	}
	return blurbs, nil
}

func (db DB) DeleteBlurb(blurbId string) error {
	if blurbId == "" {
		return errInvalidInputData
	}
	_, err := db.Exec(`delete from blurbs where blurb_id = $1`, blurbId)
	if err != nil {
		return err
	}
	return nil
}
