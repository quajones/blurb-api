package models

import (
	"time"

	"github.com/lib/pq"
)

type Blurb struct {
	ID          string    `uri:"id" json:"id,omitempty"`
	UserID      string    `json:"userID,omitempty"`
	Topic       string    `json:"topic,omitempty`
	Title       string    `json:"title,omitempty"`
	Content     string    `json:"content,omitempty"`
	Tags        []string  `json:"tags,omitEmpty"`
	CreatedDate time.Time `json:"createdDate,omitempty"`
	//PublishedDate    time.Time `json:"publishedDate,omitempty"`
	LastModifiedDate time.Time `json:"lastModifiedDate,omitempty"`
}

var (
	getBlurbsForFollowingQuery = `select a.blurb_id, a.user_id, c.topic, a.title, a.content, a.created_date, a.last_modified,
										coalesce(d.tags, '{}') as tags 
									from blurbs a join user_following b on a.user_id = b.following_id 
									left join topics c on 
										a.topic = c.topic_id 
									left join (
											select a.blurb_id, array_agg(b.tag) as tags
											from blurbs_tags a 
											join tags b on a.tag_id = b.tag_id
											group by a.blurb_id
										) d on d.blurb_id = a.blurb_id
									where b.user_id = $1
									order by a.last_modified desc`

	getBlurbsForUserQuery = `select a.blurb_id, a.user_id, b.topic, a.title, a.content, a.created_date, a.last_modified, 
									coalesce(c.tags, '{}') as tags
								from blurbs a 
								left join topics b on 
									a.topic = b.topic_id 
								left join (
										select a.blurb_id, array_agg(b.tag) as tags
										from blurbs_tags a 
										join tags b on a.tag_id = b.tag_id
										group by a.blurb_id
									) c on c.blurb_id = a.blurb_id
								where a.user_id = $1`
)

// create a blurb given a user and topic
func (db DB) CreateBlurb(userId string, topic int, title, content string) (string, error) {
	if userId == "" || title == "" || content == "" {
		return "", errInvalidInputData
	}
	row := db.QueryRow(`insert into blurbs (user_id, topic, title, content, created_date, last_modified) values ($1, $2, $3, $4, $5, $6) returning blurb_id`,
		userId, topic, title, content, time.Now(), time.Now())
	var blurb Blurb
	if err := row.Scan(&blurb.ID); err != nil {
		return "", err
	}
	return blurb.ID, nil
}
func (db DB) GetAllBlurbsForUser(userId string) ([]Blurb, error) {
	if userId == "" {
		return nil, errInvalidInputData
	}
	rows, err := db.Query(getBlurbsForUserQuery, userId)
	if err != nil {
		return nil, err
	}
	var blurbs []Blurb
	for rows.Next() {
		var blurb Blurb
		if err := rows.Scan(&blurb.ID, &blurb.UserID, &blurb.Topic, &blurb.Title,
			&blurb.Content, &blurb.CreatedDate, &blurb.LastModifiedDate, pq.Array(&blurb.Tags)); err != nil {
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
		if err := rows.Scan(&blurb.ID, &blurb.UserID, &blurb.Topic, &blurb.Title, &blurb.Content,
			&blurb.CreatedDate, &blurb.LastModifiedDate, pq.Array(&blurb.Tags)); err != nil {
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
