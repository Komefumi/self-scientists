package data

import (
	"database/sql"
	"fmt"
)

func getMapListFromSQLRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	dataMapList := []map[string]interface{}{}
	columnNames, errGettingColumns := rows.Columns()
	if errGettingColumns != nil {
		fmt.Println(errGettingColumns)
		return dataMapList, errGettingColumns
	}

	for rows.Next() {
		columnDataList := make([]interface{}, len(columnNames))
		columnDataPointerList := make([]interface{}, len(columnNames))

		for i, _ := range columnNames {
			columnDataPointerList[i] = &columnDataList[i]
		}

		if err := rows.Scan(columnDataPointerList...); err != nil {
			fmt.Println(err)
			return nil, err
		}

		dataMap := make(map[string]interface{})

		for i, colName := range columnNames {
			val := columnDataPointerList[i].(*interface{})
			dataMap[colName] = val
		}

		dataMapList = append(dataMapList, dataMap)
	}

	return dataMapList, nil
}
