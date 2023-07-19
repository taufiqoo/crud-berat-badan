package models

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/taufiqoo/technical-test-berat/config"
	"github.com/taufiqoo/technical-test-berat/entities"
)

type UserModel struct {
	conn *sql.DB
}

func NewUserModel() *UserModel {
	conn, err := config.DBConnection()
	if err != nil {
		panic(err)
	}

	return &UserModel{
		conn: conn,
	}
}

func (u *UserModel) FindAll() ([]entities.User, error) {

	rows, err := u.conn.Query("select * from user")
	if err != nil {
		return []entities.User{}, err
	}
	defer rows.Close()

	var dataUser []entities.User
	for rows.Next() {
		var user entities.User
		rows.Scan(&user.Id, &user.Tanggal, &user.Max, &user.Min, &user.Diff)

		dataUser = append(dataUser, user)
	}

	return dataUser, nil

}

func (u *UserModel) Create(user entities.User) bool {

	err := u.CalculateDifference(&user)
	if err != nil {
		fmt.Println(err)
		return false
	}

	result, err := u.conn.Exec("insert into user (tanggal, max, min, diff) values(?,?,?,?)",
		user.Tanggal, user.Max, user.Min, user.Diff)

	if err != nil {
		fmt.Println(err)
		return false
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId > 0

}

func (u *UserModel) Find(id int64, user *entities.User) error {

	// menghitung selisih untuk data User yang akan diupdate
	err := u.CalculateDifference(user)
	if err != nil {
		return err
	}

	row := u.conn.QueryRow("SELECT * FROM user WHERE id = ?", id)
	err = row.Scan(&user.Id, &user.Tanggal, &user.Max, &user.Min, &user.Diff)
	return err

}

func (u *UserModel) Update(user *entities.User) error {
	// Hitung selisih berat badan
	err := u.CalculateDifference(user)
	if err != nil {
		return err
	}

	_, err = u.conn.Exec(
		"UPDATE user SET tanggal = ?, max = ?, min = ?, diff = ? WHERE id = ?",
		user.Tanggal, user.Max, user.Min, user.Diff, user.Id,
	)

	return err
}

func (u *UserModel) Delete(id int64) error {
	_, err := u.conn.Exec("DELETE FROM user WHERE id = ?", id)
	return err
}

func (u *UserModel) CalculateDifference(user *entities.User) error {
	max, err := strconv.Atoi(user.Max)
	if err != nil {
		return err
	}

	min, err := strconv.Atoi(user.Min)
	if err != nil {
		return err
	}

	user.Diff = max - min
	return nil
}

func (u *UserModel) CalculateAverages() (float64, float64, float64) {
	var avgMax, avgMin, avgDiff float64

	// Menghitung rata-rata berat badan max
	row := u.conn.QueryRow("SELECT AVG(max) FROM user")
	row.Scan(&avgMax)

	// Menghitung rata-rata berat badan min
	row = u.conn.QueryRow("SELECT AVG(min) FROM user")
	row.Scan(&avgMin)

	// Menghitung rata-rata perbedaan berat badan
	row = u.conn.QueryRow("SELECT AVG(diff) FROM user")
	row.Scan(&avgDiff)

	return avgMax, avgMin, avgDiff
}

func (u *UserModel) FindById(id int64, user *entities.User) error {
	row := u.conn.QueryRow("SELECT * FROM user WHERE id = ?", id)
	err := row.Scan(&user.Id, &user.Tanggal, &user.Max, &user.Min, &user.Diff)
	return err
}
