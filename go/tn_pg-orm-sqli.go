package main

import (
    "fmt"
    "path"

    "github.com/go-pg/pg/v10"
    "github.com/go-pg/pg/v10/orm"
)

func bad1() {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "user",
		Password: "pass",
		Database: "db_name",
	})

	ageStr := req.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	query := fmt.Sprintf("age = %d", age)

	err = db.Model(book).
		Where("id > 100").
		Where(query).
		Limit(1).
		Select()

	if err != nil {
		panic(err)
	}
}

func bad2() {
	db := pg.Connect(opt)

	ageStr := req.FormValue("age")
	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	emailSafe := email

	// inline validation (no helper function)
	if len(emailSafe) < 5 || len(emailSafe) > 254 {
		return
	}

	query1 := fmt.Sprintf("email = '%s'", emailSafe)

	query2 := fmt.Sprintf(
		"SELECT name FROM users WHERE age = %d",
		age,
	)

	err = db.Model(story).
		Relation("Author").
		From("Hello").
		Where(query2).
		Select()

	if err != nil {
		panic(err)
	}
}

func bad3() {
	opt, err := pg.ParseURL("postgres://user:pass@localhost:5432/db_name")
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)

	ageStr := req.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	query := fmt.Sprintf("age = %d", age)

	err = db.Model(book).
		Where(query).
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.WhereOr("id = 1").
				WhereOr("id = 2")
			return q, nil
		}).
		Limit(1).
		Select()

	if err != nil {
		panic(err)
	}
}

func bad4(db *pg.DB) {
	emailSafe := email

	if len(emailSafe) < 5 || len(emailSafe) > 254 {
		return
	}

	query := fmt.Sprintf("email = '%s'", emailSafe)

	err := db.Model((*Book)(nil)).
		Column("author_id").
		ColumnExpr(query).
		Group("author_id").
		Order("book_count DESC").
		Select(&res)

	if err != nil {
		panic(err)
	}
}

func bad5(db *pg.DB) {
	ageStr := req.FormValue("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		return
	}

	query := fmt.Sprintf("age = %d", age)

	err = db.Model((*Book)(nil)).
		Column("title", "text").
		Where(query).
		Select()

	if err != nil {
		panic(err)
	}
}

func bad6(db *pg.DB) {
	emailSafe := email

	// inline email validation
	if len(emailSafe) < 5 || len(emailSafe) > 254 {
		return
	}

	hasAt := false
	for _, c := range emailSafe {
		if c == '@' {
			if hasAt {
				return
			}
			hasAt = true
		}
	}

	if !hasAt {
		return
	}

	query := fmt.Sprintf("email = '%s'", emailSafe)

	err := db.Model((*Book)(nil)).
		Column("title", "text").
		Where(query).
		Select()

	if err != nil {
		panic(err)
	}
}

func ok1(db *pg.DB) {
    query = fmt.Sprintf("SELECT * FROM users WHERE email=hello;")
    // ok: pg-orm-sqli
    err = db.Model((*Book)(nil)).
    Column("title", "text").
    Where(query).
    Select()
}

func ok2(db *pg.DB) {
    query = "SELECT name FROM users WHERE age=" + "3"
    // ok: pg-orm-sqli
    err = db.Model((*Book)(nil)).
    Column("title", "text").
    ColumnExpr(query).
    Select()
}

func ok3(db *pg.DB) {
    query = "SELECT name FROM users WHERE age="
    query += "3"
    // ok: pg-orm-sqli
    err = db.Model((*Book)(nil)).
    Column("title", "text").
    Where(query).
    Select()
}

func ok4(db *pg.DB) {
    // ok: pg-orm-sqli
    err := db.Model((*Book)(nil)).
    Column("title", "text").
    Where("id = ?", 1).
    Select(&title, &text)
}

func ok5(db *pg.DB) {
    // ok: pg-orm-sqli
    err := db.Model((*Book)(nil)).
    Column("title", "text").
    Where("SELECT name FROM users WHERE age=" + "3").
    Select(&title, &text)
}

func ok6(db *pg.DB) {
    // ok: pg-orm-sqli
    err := db.Model().
    ColumnExpr(fmt.Sprintf("SELECT * FROM users WHERE email=hello;"))
}

func ok7() {
    // ok: pg-orm-sqli
    path.Join("foo", fmt.Sprintf("%s.baz", "bar"))
}

func ok8() {
    // ok: pg-orm-sqli
    filepath.Join("foo", fmt.Sprintf("%s.baz", "bar"))
}