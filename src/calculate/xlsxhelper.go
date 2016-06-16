package calculate

import (
	. "wanbu_stat_zzmm/src/config"
	. "wanbu_stat_zzmm/src/utils"

	"github.com/tealeg/xlsx"
)

/*
func PharseFile(filePath string, rules map[int](map[int]*ZM)) error {
	xlFile, err := xlsx.OpenFile("nice.xlsx")
	CheckError(err)
	for _, sheet := range xlFile.Sheets {
		activeid, err := strconv.Atoi(strings.Trim(sheet.Name, " "))
		CheckError(err)
		for irw, row := range sheet.Rows {
			if irw != 0 { // The First Line is Name of Columns
				if len(row.Cells) >= 6 {
					id, _ := row.Cells[0].Int()
					if _, ok := rules[id]; !ok {
						rules[id] = make(map[int]*ZM)
					}

					if _, ok := rules[id][activeid]; !ok {
						rules[id][activeid] = &ZM{}
					}

					// type of rule
					t, _ := row.Cells[1].Int()
					start, _ := row.Cells[2].Int()
					end, _ := row.Cells[3].Int()
					condition, _ := row.Cells[4].Int()
					credit, _ := row.Cells[5].Int()

					switch t {
					case 1:
						rules[id][activeid].Mor = ZMRule{Start: start, End: end, Condition: condition, Credit: credit}
					case 2:
						rules[id][activeid].Mid = ZMRule{Start: start, End: end, Condition: condition, Credit: credit}
					case 3:
						rules[id][activeid].Moon = ZMRule{Start: start, End: end, Condition: condition, Credit: credit}
					default:
						fmt.Println("error type")

					}
				}

			}

		}
	}
	return err
}*/

/*
sheet0 	: Document Readme
sheet1 	: One Active Config
sheet2 	: All UserIds That No zzmm Rule
sheet3..: Rule Details(sheet3, sheet4, ...)
*/
func PharseFile(filePath string, rules map[int]*ZM, ainfo *ActiveInfo, uids map[int]bool, user2rule map[int]int) error {

	xlFile, err := xlsx.OpenFile(XlsxPath)
	CheckError(err)

	for sheetNo, sheet := range xlFile.Sheets {
		if sheetNo == 0 {

			// 使用说明无需处理
			//fmt.Println("使用说明无需处理")

		} else if sheetNo == 1 {
			// 活动规则

			activeid, err1 := sheet.Rows[0].Cells[1].Int()
			zmCondition, err2 := sheet.Rows[1].Cells[1].Int()
			zmUpline, err3 := sheet.Rows[2].Cells[1].Int()

			if err1 == nil && err2 == nil && err3 == nil {
				ainfo.Activeid = activeid
				ainfo.ZmCondition = zmCondition
				ainfo.ZmUpline = zmUpline
			}

		} else if sheetNo == 2 {
			//没有设置规则的用户

			for index, row := range sheet.Rows {
				if index != 0 {
					uid, err := row.Cells[0].Int()
					if err == nil {
						uids[uid] = true
					}
				}
			}

		} else {
			// 方案设置

			zm := &ZM{}
			//fmt.Println("方案设置")
			for i := 2; i < 5; i++ {
				start, err1 := sheet.Rows[i].Cells[3].Int()
				end, err2 := sheet.Rows[i].Cells[4].Int()
				steps, err3 := sheet.Rows[i].Cells[5].Int()
				prize, err4 := sheet.Rows[i].Cells[6].Int()
				if err1 == nil && err2 == nil && err3 == nil && err4 == nil {
					if i == 2 {
						zm.Mor.Start = start
						zm.Mor.End = end
						zm.Mor.Condition = steps
						zm.Mor.Credit = prize
					} else if i == 3 {
						zm.Mid.Start = start
						zm.Mid.End = end
						zm.Mid.Condition = steps
						zm.Mid.Credit = prize
					} else if i == 4 {
						zm.Moon.Start = start
						zm.Moon.End = end
						zm.Moon.Condition = steps
						zm.Moon.Credit = prize
					}
				}
			}
			rules[sheetNo] = zm
			for index, row := range sheet.Rows {
				if index != 0 {
					uid, err := row.Cells[0].Int()
					if err == nil {
						user2rule[uid] = sheetNo
					}
				}
			}
		}
	}
	return err
}
