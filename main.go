package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"gopkg.in/yaml.v3"
)

type (
	Conf struct {
		Db     string `yaml:"db"`
		Listen string `yaml:"listen"`
	}
	LoginData struct {
		Epassport string `json:"epassport"`
		Epassword string `json:"epassword"`
	}
	PermissionData struct {
		CanView  int `json:"canview"`
		CanEdit  int `json:"canedit"`
		CanAdmin int `json:"canadmin"`
	}
	FileData struct {
		ID       int      `json:"id"`
		Filename string   `json:"filename"`
		Tags     []string `json:"tags"`
		Data     string   `json:"data"`
		Username string   `json:"username"`
		Type     string   `json:"type"`
		Protect  string   `json:"protect"`
	}
	SearchInput struct {
		Username string   `json:"username"`
		Keywords string   `json:"keywords"`
		Tags     []string `json:"tags"`
		Type     []string `json:"type"`
	}
	TagAPI struct {
		Tag string `json:"tag"`
	}
	UserData struct {
		Username string `json:"username"`
		Edit     int    `json:"edit"`
		Admin    int    `json:"admin"`
	}
	UserAPI struct {
		Username   string   `json:"username"`
		Permission []string `json:"permission"`
	}
	UserAPI2 struct {
		Username string `json:"username"`
		Edit     bool   `json:"edit"`
		Admin    bool   `json:"admin"`
	}
)

var (
	conf Conf
)

func main() {
	processConfig()

	app := fiber.New(fiber.Config{
		Views:       html.New("./template", ".html"),
		ProxyHeader: fiber.HeaderXForwardedFor,
	})

	app.Static("/", "./vue/dist")
	app.Post("/api/login", login)
	app.Get("/api/tags", gettags)
	app.Get("/api/users", getusers)
	app.Post("/api/adduser", adduser)
	app.Post("/api/deluser", deluser)
	app.Post("/api/updateuser", updateuser)
	app.Post("/api/fileupload", fileupload)
	app.Post("/api/search", searchxdoc)
	app.Post("/api/addnewtag", addnewtag)
	app.Post("/api/deletetag", deletetag)
	app.Post("/api/deletefile", deletefile)
	app.Get("/file/:fnum", fileload)
	app.Get("/qrcode/:fnum", qrcodegen)
	app.Post("/api/download", downloadfile)
	app.Post("/api/updatemetafile", fileupdatemeta)
	app.Post("/api/updatefileupload", updatefileupload)
	log.Fatal(app.Listen(conf.Listen))
}

func processConfig() {
	confdata := readFile("/etc/xdoc.yml")
	yaml.Unmarshal(confdata, &conf)

}

func readFile(f string) []byte {
	content, err := os.ReadFile(f)
	if err != nil {
		return []byte("")
	}
	return content
}

func deletetag(c *fiber.Ctx) error {
	deltag := TagAPI{}
	err := json.Unmarshal(c.Body(), &deltag)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM tags WHERE tagname=$1;", deltag.Tag)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func deluser(c *fiber.Ctx) error {
	deluser := UserAPI{}
	err := json.Unmarshal(c.Body(), &deluser)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE username=$1;", deluser.Username)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func deletefile(c *fiber.Ctx) error {
	delfile := FileData{}
	err := json.Unmarshal(c.Body(), &delfile)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	if !chkprivateprotect(delfile) {
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM doc WHERE id=$1;", delfile.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func addnewtag(c *fiber.Ctx) error {
	newtag := TagAPI{}
	err := json.Unmarshal(c.Body(), &newtag)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO tags (tagname) VALUES ($1);", newtag.Tag)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func adduser(c *fiber.Ctx) error {
	newuser := UserAPI{}
	err := json.Unmarshal(c.Body(), &newuser)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	canedit := 0
	canadmin := 0
	username := strings.TrimSpace(strings.ToLower(newuser.Username))
	for _, v := range newuser.Permission {
		switch v {
		case "edit":
			canedit = 1
		case "admin":
			canadmin = 1
		}
	}

	_, err = db.Exec("INSERT INTO users (username,canedit,canadmin) VALUES ($1,$2,$3);", username, canedit, canadmin)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func updateuser(c *fiber.Ctx) error {
	userx := UserAPI2{}
	err := json.Unmarshal(c.Body(), &userx)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	canedit := 0
	canadmin := 0
	if userx.Edit {
		canedit = 1
	}
	if userx.Admin {
		canadmin = 1
	}

	_, err = db.Exec("UPDATE users SET canedit=$1,canadmin=$2 WHERE username=$3;", canedit, canadmin, userx.Username)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func qrcodegen(c *fiber.Ctx) error {
	fid, err := c.ParamsInt("fnum")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error1")
	}

	return c.Render("qrcode", fiber.Map{
		"ID": fmt.Sprintf("%s/file/%d", c.BaseURL(), fid),
	})

}

func fileload(c *fiber.Ctx) error {
	fid, err := c.ParamsInt("fnum")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error1")
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error2")
	}
	defer db.Close()

	qq := "SELECT data,type FROM doc WHERE id=$1"
	rows, err := db.Query(qq, fid)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error3")
	}
	defer rows.Close()

	data := ""
	typex := ""
	for rows.Next() {
		rows.Scan(&data, &typex)
		if typex != "public" {
			return c.SendString("Error5")
		}
		break
	}
	if data == "" {
		return c.SendString("Error6")
	}
	mime, _, datax := extractBase64(data)
	c.Append("Content-Type", mime)
	c.Append("Content-Disposition", "attachment")
	dataBin, err := base64.StdEncoding.DecodeString(datax)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error4")
	}
	c.Append("Content-Length", fmt.Sprintf("%d", len(dataBin)))
	return c.Send(dataBin)
}

func downloadfile(c *fiber.Ctx) error {
	dataload := FileData{}
	err := json.Unmarshal(c.Body(), &dataload)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error1")
	}

	if !chkdownloadright(dataload) {
		return c.SendString("Error5")
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error2")
	}
	defer db.Close()

	qq := "SELECT data FROM doc WHERE id=$1"
	rows, err := db.Query(qq, dataload.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.SendString("Error3")
	}
	defer rows.Close()

	data := ""
	for rows.Next() {
		rows.Scan(&data)
		return c.SendString(data)
	}
	return c.SendString("Error4")
}

func searchxdoc(c *fiber.Ctx) error {
	searchinput := SearchInput{}
	err := json.Unmarshal(c.Body(), &searchinput)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"results": []FileData{},
		})
	}

	keywordsquery := MakeKeywordsQuery(searchinput.Keywords)
	tagsquery := makeTagsQuery(searchinput.Tags)
	typequery := makeTypeQuery(searchinput.Type, searchinput.Username)

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"results": []FileData{},
		})
	}
	defer db.Close()

	query := ""
	if keywordsquery != "" {
		query = "WHERE " + keywordsquery
	}

	if tagsquery != "" {
		if query != "" {
			query = query + " and " + tagsquery
		} else {
			query = "WHERE " + tagsquery
		}
	}

	if typequery != "" {
		if query != "" {
			query = query + " and " + typequery
		} else {
			query = "WHERE " + typequery
		}
	}

	qq := "SELECT id,filename,username,tags,type,protect FROM doc " + query
	rows, err := db.Query(qq)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"results": []FileData{},
		})
	}
	defer rows.Close()

	searchresult := make([]FileData, 0)

	for rows.Next() {
		sresult := FileData{}
		tag := ""
		rows.Scan(&sresult.ID, &sresult.Filename, &sresult.Username, &tag, &sresult.Type, &sresult.Protect)
		sresult.Tags = ExtractTags(tag)
		searchresult = append(searchresult, sresult)
	}
	return c.JSON(fiber.Map{
		"results": searchresult,
	})
}

func fileupload(c *fiber.Ctx) error {
	fupload := FileData{}
	err := json.Unmarshal(c.Body(), &fupload)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	tags := encodeTags(fupload.Tags)
	datex := getDateInt()
	timex := getTimeInt()

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO doc (filename,username,tags,data,date,time,type,protect) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);", fupload.Filename, fupload.Username, tags, fupload.Data, datex, timex, fupload.Type, fupload.Protect)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func updatefileupload(c *fiber.Ctx) error {
	fupload := FileData{}
	err := json.Unmarshal(c.Body(), &fupload)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	if !chkprivateprotect(fupload) {
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("UPDATE doc SET data=$1 WHERE id=$2;", fupload.Data, fupload.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func fileupdatemeta(c *fiber.Ctx) error {
	meta := FileData{}
	err := json.Unmarshal(c.Body(), &meta)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	if !chkprivateprotect(meta) {
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	tags := encodeTags(meta.Tags)
	datex := getDateInt()
	timex := getTimeInt()

	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}
	defer db.Close()

	_, err = db.Exec("UPDATE doc SET filename=$1,username=$2,tags=$3,date=$4,time=$5,type=$6,protect=$7 WHERE id=$8;", meta.Filename, meta.Username, tags, datex, timex, meta.Type, meta.Protect, meta.ID)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"status": "false",
		})
	}

	return c.JSON(fiber.Map{
		"status": "true",
	})
}

func login(c *fiber.Ctx) error {
	loginDat := LoginData{}
	err := json.Unmarshal(c.Body(), &loginDat)
	if err != nil {
		return c.JSON(fiber.Map{
			"status": "invalid request",
		})
	}
	rr := epassport(loginDat.Epassport, loginDat.Epassword)
	if rr["status"] != "ok" {
		return c.JSON(fiber.Map{
			"status": "invalid login",
		})
	}
	permission := getUserData(loginDat.Epassport)
	return c.JSON(fiber.Map{
		"status":   "ok",
		"canview":  permission.CanView,
		"canedit":  permission.CanEdit,
		"canadmin": permission.CanAdmin,
	})
}

func gettags(c *fiber.Ctx) error {
	tags := make([]string, 0)
	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"tags": tags,
		})
	}
	defer db.Close()

	rows, err := db.Query("SELECT tagname FROM tags;")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"tags": tags,
		})
	}
	defer rows.Close()

	for rows.Next() {
		tag := ""
		rows.Scan(&tag)
		tags = append(tags, tag)
	}
	return c.JSON(fiber.Map{
		"tags": tags,
	})
}

func getusers(c *fiber.Ctx) error {
	users := make([]UserData, 0)
	db, err := sql.Open("postgres", conf.Db)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"users": users,
		})
	}
	defer db.Close()

	rows, err := db.Query("SELECT username,canedit,canadmin FROM users;")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return c.JSON(fiber.Map{
			"users": users,
		})
	}
	defer rows.Close()

	for rows.Next() {
		user := UserData{}
		rows.Scan(&user.Username, &user.Edit, &user.Admin)
		users = append(users, user)
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}
