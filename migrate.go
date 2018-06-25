package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Basename []string

func (p Basename) Len() int      { return len(p) }
func (p Basename) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Basename) Less(i, j int) bool {
	return filepath.Base(p[i]) < filepath.Base(p[j])
}

func main() {
	var scanDirName, mysqlHost, mysqlPort, mysqlUser, mysqlPassword string
	flag.StringVar(&scanDirName, "dir", "migrations", "Dir with sql files")
	flag.StringVar(&mysqlHost, "host", "localhost", "MySQL host")
	flag.StringVar(&mysqlPort, "port", "3306", "MySQL port")
	flag.StringVar(&mysqlUser, "user", "root", "MySQL user")
	flag.StringVar(&mysqlPassword, "password", "", "MySQL password")
	flag.Parse()

	fileList, err := scanFiles(scanDirName, ".sql")
	check(err)

	if len(fileList) == 0 {
		log.Fatal(".sql files not found")
	}
	sort.Sort(Basename(fileList))
	var strBuilder strings.Builder
	for _, filePath := range fileList {
		dat, err := ioutil.ReadFile(filePath)
		check(err)
		strBuilder.Write(dat)
	}

	//mysql
	mysqlConf := mysql.NewConfig()
	mysqlConf.Net = "tcp"
	mysqlConf.Addr = mysqlHost + ":" + mysqlPort
	mysqlConf.User = mysqlUser
	mysqlConf.Passwd = mysqlPassword
	mysqlConf.MultiStatements = true

	db, err := sql.Open("mysql", mysqlConf.FormatDSN())
	//nolint
	defer db.Close()
	check(err)
	_, err = db.Exec(strBuilder.String())
	check(err)
}

func scanFiles(dir, ext string) ([]string, error) {
	fileList := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList, err
}
