package database

import "context"

func (db *Database) CheckUserExist(mail string) (bool, error) {
	var exist bool
	query := "select exists(select mail from user_t where mail = $1);"
	err := db.pool.QueryRow(context.Background(), query, mail).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (db *Database) InsertUser(mail, phone string) error {
	query := "insert into user_t (mail, phone) values ($1, $2);"
	_, err := db.pool.Exec(context.Background(), query, mail, phone)
	if err != nil {
		return err
	}
	return nil
}
