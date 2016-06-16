package config

import (
	"flag"
	"fmt"
	config "github.com/msbranco/goconfig"
	"strconv"
	"strings"
)

var (
	DBConfig     string
	XlsxPath     string
	Trix         string
	Consumerip   string
	Consumerport string
	Activeids    = []int{}
	HOST         string
	config_path  = flag.String("f", "../conf/config.conf", "config path")
	err          error
)

func LoadConfig() {
	flag.Parse()
	cf, _ := config.ReadConfigFile(*config_path)
	rdip1, err1 := cf.GetString("DBCONN1", "IP")
	if err1 != nil {
		fmt.Println("rdip1", err1)
	}
	rdusr1, err2 := cf.GetString("DBCONN1", "USERID")
	if err2 != nil {
		fmt.Println("rdusr1", err2)
	}
	rdpwd1, err3 := cf.GetString("DBCONN1", "USERPWD")
	if err3 != nil {
		fmt.Println("rdpwd1", err3)
	}
	rdname1, err4 := cf.GetString("DBCONN1", "DBNAME")
	if err4 != nil {
		fmt.Println("rdname1", err4)
	}

	Consumerip, err = cf.GetString("CONSUMER", "IP")
	if err != nil {
		fmt.Println(err)
	}
	Consumerport, err = cf.GetString("CONSUMER", "PORT")
	if err != nil {
		fmt.Println(err)
	}

	Trix, err = cf.GetString("PATH", "TRIX")
	if err != nil {
		fmt.Println(err)
	}

	DBConfig = rdusr1 + ":" + rdpwd1 + "@tcp(" + rdip1 + ")/" + rdname1 + "?charset=utf8"
	XlsxPath, _ = cf.GetString("PATH", "XlSX")
	ids, err0 := cf.GetString("INFO", "ACTIVE")
	if err0 != nil {
		fmt.Println(err0)
	}
	getActives(ids)
}

func getActives(ids string) {
	aids := strings.Split(ids, ",")
	for _, aid := range aids {
		id, _ := strconv.Atoi(aid)
		Activeids = append(Activeids, id)
	}
}
