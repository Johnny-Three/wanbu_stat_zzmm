package calculate

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
	. "wanbu_stat_zzmm/src/config"
	. "wanbu_stat_zzmm/src/logs"
	. "wanbu_stat_zzmm/src/utils"
)

var (
	db         *sql.DB
	Sql_ch     chan string
	cc_num     int //The Number Of Sql to Execute concurrent
	Refresh_ch chan *Refresh
)

func DbInit() {

	Sql_ch = make(chan string, 0)
	Refresh_ch = make(chan *Refresh, 0)

	var err error
	db, err = sql.Open("mysql", DBConfig)
	CheckError(err)
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	db.Ping()

	go run()

	Logger.Info(DBConfig)
}

func run() {

	for {
		select {
		case sql := <-Sql_ch:

			result, err := db.Exec(sql)
			//Logger.Debug(sql)
			if err != nil {
				Logger.Critical("Execute SQL Error : ", err, " [Sql] : ", sql)
			} else {
				ra, _ := result.RowsAffected()
				if ra == 0 {
					Logger.Critical("Failed Insert [SQL] ", sql)
				}
			}
		}

	}
	defer close(Sql_ch)
}

func GetActiveInfo(ai *ActiveInfo) (err error) {
	sql := "select a.starttime, a.endtime,a.prestarttime, a.preendtime, IFNULL(b.credit2distance, 0), systemflag from wanbu_club_online a, wanbu_rule_config b where a.activeid = b.activeid and a.activeid = ?"
	row := db.QueryRow(sql, ai.Activeid)

	err = row.Scan(&ai.StartTime, &ai.EndTime, &ai.PreStartTime, &ai.PreEndTime, &ai.Credit2Distance, &ai.Systemflag)
	CheckError(err)
	return err
}

func GetGroupId(activeid, userid int) ([]int, error) {
	sql := "SELECT wgu.groupid FROM wanbu_group_user wgu ,wanbu_club_online wco,wanbu_group_info gi WHERE wco.storeflag!=2 AND  wco.closetime>UNIX_TIMESTAMP(NOW()) AND gi.activeid = wco.activeid AND gi.groupid = wgu.groupid AND gi.activeid = ? AND wgu.userid = ? UNION SELECT wgq.groupid FROM wanbu_group_quit wgq , wanbu_club_online wco,wanbu_group_info gi WHERE wco.storeflag!=2 AND wco.closetime>UNIX_TIMESTAMP(NOW())  AND gi.activeid = wco.activeid AND  gi.groupid = wgq.groupid AND gi.activeid = ? AND wgq.userid = ? "
	rows, err := db.Query(sql, activeid, userid, activeid, userid)
	CheckError(err)
	defer rows.Close()
	groups := make([]int, 0)
	gid := 0
	for rows.Next() {
		err := rows.Scan(&gid)
		if err == nil {
			groups = append(groups, gid)
		}
	}
	return groups, err
}

func LoopRefresh() {
	sql_ := "select uploadid, userid, activeid, walkdate from wanbu_data_zmrefresh_queue_" + Trix +
		" limit 1000"
	for {
		time.Sleep(time.Second * 2)

		rows, err := db.Query(sql_)
		CheckError(err)
		defer rows.Close()
		for rows.Next() {
			data := Refresh{}
			uploadid := sql.NullInt64{}
			err := rows.Scan(&uploadid, &data.Userid, &data.Activeid, &data.Walkdate)
			if err == nil {
				data.Uploadid = uploadid.Int64
				Refresh_ch <- &data
			}
		}
	}
}

func DelCredit(userid int, walkdate int64, aid int) (err error) {

	sql_ := "delete from wanbu_member_credit where taskid = -99 and userid = ? and walkdate = ? and activeid = ?"
	result, err := db.Exec(sql_, userid, walkdate, aid)

	sql := fmt.Sprintf("delete from wanbu_member_credit where taskid = -99 and userid = %d and walkdate = %d and activeid = %d", userid, walkdate, aid)

	if err != nil {
		Logger.Debug("Execute SQL Error : ", err, " [Sql] : ", sql)
	} else {
		ra, _ := result.RowsAffected()
		if ra == 0 {
			Logger.Debug("Failed Delete [SQL] ", sql)
		}
	}
	//fmt.Println("delete old ", err)
	return err
}

func DelQueue(uploadid int64) {
	sql_ := "delete from wanbu_data_zmrefresh_queue_" + Trix + " where uploadid = ?"
	result, err := db.Exec(sql_, uploadid)
	if err != nil {
		Logger.Debug("Execute SQL Error : ", err, " [Sql] : ", sql_)
	} else {
		ra, _ := result.RowsAffected()
		if ra == 0 {
			Logger.Debug("Failed Delete [SQL] ", sql_)
		}
	}
	//fmt.Println(sql_)
}

func GetWalkHour(userid int, walkdate int64) (wh_ []int, sum int, err error) {
	wh := make([]string, 26)
	wh_ = make([]int, 26)
	sql_ := "SELECT  `hour0`, `hour1`, `hour2`, `hour3`, `hour4`, `hour5`, `hour6`, `hour7`, `hour8`, `hour9`, `hour10`, `hour11`, `hour12`, `hour13`, `hour14`, `hour15`, `hour16`, `hour17`, `hour18`, `hour19`, `hour20`, `hour21`, `hour22`, `hour23`,`hour24`,`hour25` FROM `wanbu`.`wanbu_data_walkhour` WHERE `userid` = ? AND `walkdate` = ? "
	rows := db.QueryRow(sql_, userid, walkdate)
	err = rows.Scan(&wh[0], &wh[1], &wh[2], &wh[3], &wh[4], &wh[5], &wh[6], &wh[7], &wh[8], &wh[9], &wh[10], &wh[11], &wh[12], &wh[13], &wh[14], &wh[15], &wh[16], &wh[17], &wh[18], &wh[19], &wh[20], &wh[21], &wh[22], &wh[23], &wh[24], &wh[25])
	sum = 0
	for index, h := range wh {
		tmpv := strings.Split(h, ",")
		if len(tmpv) == 4 {
			step1, err1 := strconv.Atoi(tmpv[0])
			step2, err2 := strconv.Atoi(tmpv[2])
			if err1 == nil && err2 == nil {
				wh_[index] = step1 + step2
				if index > 1 {
					sum += wh_[index]
				}
			} else {
				Logger.Debug("Error Transfer User Walk Hour Data", userid, walkdate)
			}
		} else if len(tmpv) == 6 {
			steps, err := strconv.Atoi(tmpv[0])
			if err == nil {
				wh_[index] = steps
				if index > 1 {
					sum += wh_[index]
				}
			} else {
				Logger.Debug("Error Transfer User Walk Hour Data", userid, walkdate)
			}
		}
	}
	//fmt.Println(sum, wh_, err)
	return wh_, sum, err
}
