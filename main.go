package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dix0nym/FMD2-Importer/v2/protos"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"google.golang.org/protobuf/proto"
	"log"
	"net/url"
	"os"
	"strings"
)

type Mapping struct {
	ID   string `json:"m.ID"`
	Name string `json:"m.Name"`
	URL  string `json:"m.RootURL"`
}

func getMapping() (map[string]Mapping, error) {
	in, err := os.ReadFile("mapping.json")
	if err != nil {
		return nil, err
	}
	var mappings []Mapping
	if err := json.Unmarshal(in, &mappings); err != nil {
		return nil, err
	}

	mappingMap := make(map[string]Mapping)
	for _, m := range mappings {
		u, err := url.Parse(m.URL)
		if err != nil {
			return nil, err
		}
		mappingMap[u.Host] = m
	}
	return mappingMap, nil
}

const (
	SQLCreateTable    = `CREATE TABLE IF NOT EXISTS "favorites" ("id" VARCHAR(3000) NOT NULL PRIMARY KEY,"order" INTEGER,"enabled" BOOLEAN,"moduleid" TEXT,"link" TEXT,"title" TEXT,"status" TEXT,"currentchapter" TEXT,"downloadedchapterlist" TEXT,"saveto" TEXT,"dateadded" DATETIME,"datelastchecked" DATETIME,"datelastupdated" DATETIME);`
	SQLGetLastOrderId = "SELECT MAX('order') from favorites;"
	SQLInsertFavorite = "INSERT INTO favorites VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, strftime('%s', 'now'), ?, ?);"
)

func start(savePath string, backupPath string) {
	in, err := os.ReadFile(backupPath)
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	backup := &protos.Backup{}
	if err := proto.Unmarshal(in, backup); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}

	db, err := sql.Open("sqlite3", "favorites.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(SQLCreateTable)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow(SQLGetLastOrderId)
	var lastOrder int
	err = row.Scan(&lastOrder)
	if err != nil {
		lastOrder = 0
	}

	sources := backup.GetBackupSources()
	if len(sources) == 0 {
		log.Printf("sources empty: len: %d\n", len(sources))
	}
	for _, i := range sources {
		log.Printf("%s - %d\n", i.Name, i.SourceId)
	}

	mappings, err := getMapping()
	if err != nil {
		log.Fatal(err)
	}

	mangas := backup.BackupManga
	for _, i := range mangas {
		log.Printf("%s\n", i.String())

		chapters := i.GetChapters()
		if len(chapters) == 0 {
			log.Printf("no chapters found for %s\n", i.Title)
			continue
		}

		baseUrl := i.Chapters[0].Scanlator
		if baseUrl == nil {
			u, err := url.Parse(*i.ThumbnailUrl)
			if err != nil {
				log.Printf("failed to find baseUrl for %s\n", i.Title)
				continue
			}
			baseUrl = &u.Host
		}

		mapping, ok := mappings[*baseUrl]
		if !ok {
			log.Printf("failed to find mapping for %s\n", *baseUrl)
			continue
		}

		moduleId := mapping.ID
		path := i.Url
		if strings.Contains(i.Url, "http") {
			u, err := url.Parse(i.Url)
			if err != nil {
				log.Printf("failed to parse url for %s: '%s'\n", i.Title, i.Url)
				continue
			}
			path = u.Path
		}
		id := fmt.Sprintf("%s%s", moduleId, path)
		saveTo := fmt.Sprintf("%s%s", savePath, i.Title)
		lastOrder += 1
		res, err := db.Exec(SQLInsertFavorite,
			id, lastOrder, 1, moduleId, path, i.Title, i.Status, len(i.Chapters), "", saveTo, nil, nil)
		if err != nil {
			log.Fatal(err)
		}

		var lastId int64
		if lastId, err = res.LastInsertId(); err != nil {
			log.Fatal(err)
		}
		log.Printf("Inserted %s at idx %d\n", i.Title, int(lastId))
	}
}

func main() {
	app := &cli.App{
		Name: "MD-Backup-Converter",
	}
}
