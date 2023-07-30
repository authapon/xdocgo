package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func epassport(username, password string) map[string]interface{} {
	cc := fiber.Post("https://elogin.rmutsv.ac.th/json")
	cc.InsecureSkipVerify()
	cc.JSON(fiber.Map{
		"username": username,
		"password": password,
	})
	_, body, err := cc.String()
	if err != nil {
		edat := make(map[string]interface{})
		edat["success"] = "false"
		return edat
	}
	rdat := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &rdat); err != nil {
		log.Printf("Error: %v\n", err)
		edat := make(map[string]interface{})
		edat["status"] = "none"
		return edat
	}
	return rdat
}

func extractBase64(data string) (string, string, string) {
	sp1 := strings.Split(data, ":")
	sp2 := strings.Split(sp1[1], ";")
	sp3 := strings.Split(sp2[1], ",")
	return sp2[0], sp3[0], sp3[1]
}

func ExtractTags(tag string) []string {
	results := make([]string, 0)
	t := strings.Split(tag, ":")
	for _, tt := range t {
		if tt != "" {
			results = append(results, tt)
		}
	}
	return results
}

func makeTypeQuery(typex []string, username string) string {
	query := ""

	for _, typexx := range typex {
		if query == "" {
			if typexx == "private" {
				query = "(type LIKE " + "'" + typexx + "' and username = '" + username + "')"
			} else {
				query = "type LIKE " + "'" + typexx + "'"
			}
		} else {
			if typexx == "private" {
				query = query + " or " + "(type LIKE " + "'" + typexx + "' and username = '" + username + "')"
			} else {
				query = query + " or " + "type LIKE " + "'" + typexx + "'"
			}
		}
	}
	if query == "" {
		return "( type LIKE 'public' )"
	}
	return " (" + query + ")"
}

func makeTagsQuery(tags []string) string {
	query := ""
	for _, tag := range tags {
		if query == "" {
			query = "tags LIKE " + "'%:" + tag + ":%'"
		} else {
			query = query + " and " + "tags LIKE " + "'%:" + tag + ":%'"
		}
	}
	return query
}

func MakeKeywordsQuery(keywords string) string {
	query := ""
	keys := strings.Split(keywords, " ")
	for i := range keys {
		keys[i] = strings.TrimSpace(keys[i])
	}

	for i, key := range keys {
		if i == 0 {
			query = "filename LIKE " + "'%" + key + "%'"
		} else {
			query = query + " and " + "filename LIKE " + "'%" + key + "%'"
		}
	}
	return query
}

func encodeTags(tags []string) string {
	t := ""
	for i, tt := range tags {
		if i == 0 {
			t = ":" + tt + ":"
		} else {
			t = t + tt + ":"
		}
	}
	return t
}

func getDateInt() int {
	t := time.Now()
	y, m, d := t.Date()
	return (y * 10000) + (int(m) * 100) + d
}

func getTimeInt() int {
	t := time.Now()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	return (h * 10000) + (m * 100) + (s)
}

func getUserData(user string) PermissionData {
	r := PermissionData{}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return r
	}
	defer db.Close()

	rows, err := db.Query("SELECT canedit,canadmin FROM users WHERE username=$1;", user)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return r
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		rows.Scan(&r.CanEdit, &r.CanAdmin)
	}
	if found {
		r.CanView = 1
	}
	return r
}

func chkprivateprotect(data FileData) bool {
	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT username,type,protect FROM doc WHERE id=$1;", data.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		username := ""
		typedoc := ""
		protect := ""
		rows.Scan(&username, &typedoc, &protect)
		if username == data.Username {
			return true
		}
		if protect == "unprotect" && typedoc != "private" {
			return true
		}
	}
	return false
}

func chkdownloadright(data FileData) bool {
	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT username,type FROM doc WHERE id=$1;", data.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		username := ""
		typedoc := ""
		rows.Scan(&username, &typedoc)
		if typedoc == "public" {
			return true
		}
		if username == data.Username && typedoc == "private" {
			return true
		}
		if username != "" && typedoc != "private" {
			return true
		}
	}
	return false
}
