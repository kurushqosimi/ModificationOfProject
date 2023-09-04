package repositories

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"

	"main/pkg/models"
)

type Repository struct {
	DB *gorm.DB
}

func (db *Repository) CheckUser(user *models.User, hash []byte) (int, error) {
	var check models.User
	sqlQuery := `select * from users where username = ?;`
	log.Println(string(hash))
	err := db.DB.Raw(sqlQuery, user.Username).Scan(&check).Error
	if err != nil {
		log.Println("ошибка в 22 строке в файле репазиторий")
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(check.Password), []byte(user.Password))
	if err != nil {
		return 0, errors.New("неправильный пароль")
	}
	if check.Username == "" {
		//log.Println("нету такого аккаунта")
		return 0, errors.New("нету такого аккаунта")
	}
	log.Println(user)
	if user.Active == false {
		return 0, errors.New("нету такого аккаунта")
	}
	log.Println(user.ID)
	return check.ID, nil
}

func GetConnection(config *models.Config) (*Repository, error) {
	dbUri := " host = " + config.DBSetting.Host + " port = " + config.DBSetting.Port + " user = " + config.DBSetting.Username + " password = " + config.DBSetting.Password + " dbname = " + config.DBSetting.Database
	db, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{}) //installs connection with database
	if err != nil {
		return nil, err
	}
	//	sqlQuery := `create table if not exists notes(
	//    id bigserial primary key ,
	//    content text not null,
	//    date timestamp default current_timestamp,
	//    active bool default true,
	//    updated_at timestamp,
	//    deleted_at timestamp
	//);`
	//err = db.Exec(sqlQuery).Error
	//if err != nil {
	//	return nil, err
	//}
	//if err != nil {
	//	return nil, err
	//}
	var rep Repository
	rep.DB = db
	return &rep, nil
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (db *Repository) Create(note *models.Notes) error {
	sqlQuery := `insert into notes (content, user_id) values (?, ?)`

	err := db.DB.Exec(sqlQuery, note.Content, note.UserID).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *Repository) Read(id *int) (models.Notes, error) {
	var note models.Notes
	sqlQuery := `select * from notes where id = ? `
	err := db.DB.Raw(sqlQuery, id).Scan(&note).Error
	if err != nil {
		log.Println("Error is in db.go 44th line!")
		note = *new(models.Notes)
		return note, err
	}
	if note.Content == "" {
		errors.New("there is not such note")
		note = *new(models.Notes)
		return note, err
	}
	return note, nil
}

func (db *Repository) Update(note *models.Notes) {
	sqlQuery := `update notes set content = ?, updated_at = ? where id = ? and active = true`
	err := db.DB.Exec(sqlQuery, note.Content, time.Now(), note.ID).Error
	if err != nil {
		log.Println(err)
		return
	}
}

func (db *Repository) Delete(id *int) {
	sqlQuery := `update notes set active = false, deleted_at = ? where id = ? and active = true`
	err := db.DB.Exec(sqlQuery, time.Now(), id).Error
	if err != nil {
		log.Println(err)
		return
	}
}

func (db *Repository) ReadAll() ([]models.Notes, error) {
	sqlQuery := `select * from notes where active = true`
	var notes []models.Notes
	err := db.DB.Raw(sqlQuery).Scan(notes).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return notes, nil
}
func (db *Repository) UserRegistration(user *models.User, hash []byte) error {
	sqlQuery := `select username from users where username = ?`
	var name string
	err := db.DB.Raw(sqlQuery, user.Username).Scan(&name).Error
	if err != nil {
		log.Println(err)
		return err
	}
	if name != "" {
		return errors.New("this kind of account already exists")
	}
	sqlQuery = `insert into users (username, password) values (?, ?)`
	err = db.DB.Exec(sqlQuery, user.Username, hash).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
