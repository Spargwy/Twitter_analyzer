package storage

import (
	_ "github.com/lib/pq" //must for work psql
)

//CreateTables table in database
func (database *PostgreSQL) CreateTables() error {

	_, err := database.db.Query("CREATE TABLE IF NOT EXISTS ACCOUNTS_DATA(FOLLOWERS INTEGER, FOLLOWING INTEGER, TWEETS INTEGER, USERNAME VARCHAR(40))")
	if err != nil {
		return err
	}
	return nil
}

func (database *PostgreSQL) InsertData(users []Account) error {
	query := "INSERT INTO ACCOUNTS_DATA(followers, following, tweets, username) VALUES($1, $2, $3, $4)"
	for i := range users {
		_, err := database.db.Exec(query, users[i].Followers, users[i].Following, users[i].Tweets, users[i].Data.Username)
		if err != nil {
			return err
		}
	}

	return nil
}
