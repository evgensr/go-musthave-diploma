package pg

import (
	"errors"
	"log"
	"strings"
)

func (box *Box) Get(key string) (Line, error) {

	var line Line
	err := box.db.QueryRow("SELECT original_url, short_url, user_id, correlation_id, status FROM  short  WHERE  short_url = $1",
		key,
	).Scan(&line.URL, &line.Short, &line.User, &line.CorrelationID, &line.Status)

	log.Println("err: ", err)
	log.Println(line)
	return line, err

}

// GetByUser получить url по id юзера
func (box *Box) GetByUser(idUser string) (lines []Line) {
	var line []Line
	var bLine Line

	log.Println(idUser)
	rows, err := box.db.Query("SELECT original_url, short_url, user_id, correlation_id, status FROM  short  WHERE  user_id = $1",
		idUser,
	)
	// обязательно закрываем перед возвратом функции
	defer func() {
		errClose := rows.Close()
		if errClose != nil {
			log.Println(errClose)
		}
	}()

	if err != nil {
		log.Println("err ** ", err)
	}

	for rows.Next() {
		err = rows.Scan(&bLine.URL, &bLine.Short, &bLine.User, &bLine.CorrelationID, &bLine.Status)
		if err != nil {
			log.Println("Scan ", err)
		}
		log.Println("original_url ", bLine)
		line = append(line, bLine)

	}

	log.Println("err: ", err)
	log.Println("lin: ", bLine)
	return line

}

func (box *Box) Set(line Line) error {

	var id int64
	err := box.db.QueryRow("INSERT INTO short (original_url, short_url, user_id, correlation_id) VALUES ($1, $2, $3, $4) RETURNING id",
		line.URL,
		line.Short,
		line.User,
		line.CorrelationID,
	).Scan(&id)

	//log.Println("err: ", err)
	//log.Println(id)

	// log.Println(strings.Contains(err.Error(), "duplicate"))
	duplicate := false
	if err != nil {
		duplicate = strings.Contains(err.Error(), "duplicate")
	}

	if duplicate {
		return errors.New("duplicate")
	}

	return nil
}

func (box *Box) Delete(line []Line) error {
	log.Println(line)
	box.chTaskDeleteURL <- line
	return nil
}
