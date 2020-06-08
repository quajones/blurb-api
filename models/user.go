package models

import (
	"errors"
	"time"

	"github.com/alexedwards/argon2id"

	"github.com/google/uuid"
)

var (
	errInvalidInputData      = errors.New("invalid input data")
	errNoRowsAffected        = errors.New("no rows affected")
	getFollowersForUserQuery = `select b.user_id, b.username, b.email 
							from user_following a 
							left join users b on 
							a.user_id = b.user_id where a.following_id = $1`
	getFollowingForUserQuery = `select b.user_id, b.username, b.email 
							from user_following a 
							left join users b on 
							a.following_id = b.user_id where a.user_id = $1`
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Email    string    `json:"email,omitempty"`
	Username string    `json:"username,omitempty"`
	HashedPW string    `json:"password,omitempty"`
}

func (db DB) GetAllUsers() ([]User, error) {
	var users []User
	rows, err := db.Query(`select user_id, username, email from users`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (db DB) GetFollowingForUser(userId string) ([]User, error) {
	var following []User
	rows, err := db.Query(getFollowingForUserQuery, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var follow User
		if err := rows.Scan(&follow.ID, &follow.Username, &follow.Email); err != nil {
			return nil, err
		}
		following = append(following, follow)
	}
	return following, nil
}

func (db DB) GetFollowersForUser(userId string) ([]User, error) {
	var followers []User
	rows, err := db.Query(getFollowersForUserQuery, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var follower User
		if err := rows.Scan(&follower.ID, &follower.Username, &follower.Email); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return followers, nil
}

func (db DB) UnfollowUser(userId, unfollowingId string) error {
	if userId == "" || unfollowingId == "" {
		return errInvalidInputData
	}
	_, err := db.Exec(`delete from user_following where user_id=$1 and following_id=$2`, userId, unfollowingId)
	if err != nil {
		return err
	}
	return nil
}

func (db DB) FollowUser(userId, followingId string) error {
	_, err := db.Exec(`insert into user_following (user_id, following_id) values ($1, $2)`, userId, followingId)
	if err != nil {
		return err
	}
	return nil
}

func (db DB) LogIn(username, password string) (User, error) {
	var user User
	row := db.QueryRow(`select user_id, username, hashedpassword, email from users where username=$1`, username)
	if err := row.Scan(&user.ID, &user.Username, &user.HashedPW, &user.Email); err != nil {
		return User{}, err
	}
	_, err := db.Exec(`update users set last_login=$1 where username=$2 `, time.Now(), user.Username)
	if err != nil {
		return User{}, err
	}
	match, err := argon2id.ComparePasswordAndHash(password, user.HashedPW)
	if err != nil {
		return User{}, err
	}
	if !match {
		return User{}, errors.New("invalid password")
	}
	return user, nil
}

func (db DB) CreateUser(username, password, email string) (*User, error) {
	var user User
	if username == "" || password == "" || email == "" {
		return nil, errInvalidInputData
	}
	hashedPw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(`insert into users (username, email, hashedpassword) values ($1, $2, $3) RETURNING user_id, username, email`,
		username, email, hashedPw)

	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (db DB) GetUser(id string) (User, error) {
	var user User
	row := db.QueryRow(`select user_id, username, email from users where user_id=$1`, id)
	if err := row.Scan(&user.ID, &user.Email, &user.Username); err != nil {
		return User{}, err
	}
	return user, nil
}

func (db DB) DeleteUser(id string) error {
	res, err := db.Exec(`delete from users where user_id=$1`, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errNoRowsAffected
	}
	return nil
}
