package gol

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	structFieldTagNameColumn = "column"

	resultKeyModeNone = iota
	resultKeyModeCamelCase
	resultKeyModeSnakeCase

	joinModeInner = iota
	joinModeLeft
	joinModeRight

	queryModeOne = iota
	queryModeAll
	queryModeString
	queryModeFormat
	queryModeIs
	queryModeIsNot
	queryModeIsNull
	queryModeIsNullNot
	queryModeLike
	queryModeLikeNot
	queryModeIn
	queryModeInNot
	queryModeGt
	queryModeGte
	queryModeLt
	queryModeLte
	queryModeNest
	queryModeNestClose

	queryPrefixNone = iota
	queryPrefixAnd
	queryPrefixOr

	Asc = iota
	Desc
)

type tableType struct {
	Str      string
	TablePtr interface{}
	TableAs  string
}

type joinType struct {
	Mode           int
	TablePtr       interface{}
	TableAs        string
	ColumnPtr      interface{}
	WhereColumnPtr interface{}
}

type joinWhereType struct {
	Mode      int
	Prefix    int
	TablePtr  interface{}
	Str       string
	ColumnPtr interface{}
	ValueList []interface{}
}

type selectType struct {
	Mode      int
	Str       string
	ColumnPtr interface{}
	ColumnAs  interface{}
}

type setType struct {
	ColumnPtr interface{}
	Value     interface{}
}

type valuesColumnType struct {
	ColumnPtr interface{}
}

type valuesType struct {
	Value interface{}
}

type whereType struct {
	Mode      int
	Prefix    int
	Str       string
	ColumnPtr interface{}
	ValueList []interface{}
}

type groupByType struct {
	Mode      int
	Str       string
	ColumnPtr interface{}
}

type havingType struct {
	Mode      int
	Prefix    int
	Str       string
	ColumnPtr interface{}
	ValueList []interface{}
}

type orderByType struct {
	Mode      int
	Order     int
	Str       string
	ColumnPtr interface{}
}

type buildType struct {
	Table          string
	TableForSelect string
	Join           string
	Select         string
	Set            string
	ValuesColumn   string
	Values         string
	Where          string
	WhereForSelect string
	GroupBy        string
	Having         string
	Order          string
	Limit          string
	Offset         string
	ValueList      []interface{}
}

type metaType struct {
	TableBase     string
	TableAsBase   string
	Table         string
	TableAs       string
	TableColumn   string
	TableAsColumn string
	ColumnBase    string
	Column        string
}

type QueryType struct {
	DB                *sql.DB
	TX                *sql.Tx
	modeDatabaseType  string
	modeLog           bool
	modeResultKey     int
	modeResetAuto     bool
	getTableName      func(string, string) (string, string)
	getColumnName     func(string, string, string) (string, string, string)
	getPlaceholder    func() string
	Table             *tableType
	JoinList          []*joinType
	JoinWhereList     []*joinWhereType
	SelectList        []*selectType
	SetList           []*setType
	ValuesColumnList  []*valuesColumnType
	ValuesColumnCount int
	ValuesList        [][]*valuesType
	WhereList         []*whereType
	Limit             int
	Offset            int
	GroupByList       []*groupByType
	HavingList        []*havingType
	OrderByList       []*orderByType
	Data              *buildType
	MetaMap           map[string]*metaType
}

func (rec *QueryType) Init(db *sql.DB, tx *sql.Tx, databaseType string) {
	rec.DB = db
	rec.TX = tx

	rec.modeResetAuto = true

	switch databaseType {
	case DatabaseTypeMysql:
		rec.modeDatabaseType = databaseType
		rec.getTableName = rec.getTableNameMysql
		rec.getColumnName = rec.getColumnNameMysql
		rec.getPlaceholder = rec.getPlaceHolderForMysql
	case DatabaseTypePostgresql:
		rec.modeDatabaseType = databaseType
		rec.getTableName = rec.getTableNamePostgresql
		rec.getColumnName = rec.getColumnNamePostgresql
		rec.getPlaceholder = rec.getPlaceHolderForPostgresql
	default:
	}

	rec.Reset()
}

func (rec *QueryType) Reset() {
	rec.ValuesColumnCount = 0
}

func (rec *QueryType) SetModeLog(mode bool) {
	rec.modeLog = mode
}

func (rec *QueryType) SetModeReset(auto bool) {
	rec.modeResetAuto = auto
}

func (rec *QueryType) SetModeResultKey() {
	rec.modeResultKey = resultKeyModeNone
}

func (rec *QueryType) SetModeResultKeyCamelCase() {
	rec.modeResultKey = resultKeyModeCamelCase
}

func (rec *QueryType) SetModeResultKeySnakeCase() {
	rec.modeResultKey = resultKeyModeSnakeCase
}

func (rec *QueryType) getTableNameMysql(tableBase string, tableAsBase string) (string, string) {
	table := fmt.Sprintf("%v", tableBase)
	tableAs := table
	if tableAsBase != "" {
		tableAs = fmt.Sprintf("%v", tableBase)
	}
	return table, tableAs
}

func (rec *QueryType) getColumnNameMysql(table string, tableAs string, columnBase string) (string, string, string) {
	tableColumn := fmt.Sprintf("%v.%v", table, columnBase)
	tableAsColumn := fmt.Sprintf("%v.%v", tableAs, columnBase)
	column := fmt.Sprintf("%v", columnBase)
	return tableColumn, tableAsColumn, column
}

func (rec *QueryType) getPlaceHolderForMysql() string {
	return "?"
}

func (rec *QueryType) getTableNamePostgresql(tableBase string, tableAsBase string) (string, string) {
	table := fmt.Sprintf("\"%v\"", tableBase)
	tableAs := table
	if tableAsBase != "" {
		tableAs = fmt.Sprintf("\"%v\"", tableBase)
	}
	return table, tableAs
}

func (rec *QueryType) getColumnNamePostgresql(table string, tableAs string, columnBase string) (string, string, string) {
	tableColumn := fmt.Sprintf("\"%v\".\"%v\"", table, columnBase)
	tableAsColumn := fmt.Sprintf("%v.\"%v\"", tableAs, columnBase)
	column := fmt.Sprintf("\"%v\"", columnBase)
	return tableColumn, tableAsColumn, column
}

func (rec *QueryType) getPlaceHolderForPostgresql() string {
	str := ""
	rec.ValuesColumnCount++
	str = fmt.Sprintf("$%v", rec.ValuesColumnCount)

	return str
}

func (rec *QueryType) setTable(str string, tablePtr interface{}, tableAs string) {
	tableData := &tableType{
		Str:      "",
		TablePtr: tablePtr,
		TableAs:  tableAs,
	}

	rec.Table = tableData
	rec.Data = nil
}

func (rec *QueryType) SetTable(tablePtr interface{}) {
	rec.setTable("", tablePtr, "")
}

func (rec *QueryType) SetTableAs(tablePtr interface{}, tableAs string) {
	rec.setTable("", tablePtr, tableAs)
}

func (rec *QueryType) setJoin(mode int, tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{}) {
	joinData := &joinType{
		Mode:           mode,
		TablePtr:       tablePtr,
		TableAs:        tableAs,
		ColumnPtr:      columnPtr,
		WhereColumnPtr: whereColumnPtr,
	}

	rec.JoinList = append(rec.JoinList, joinData)
	rec.Data = nil
}

func (rec *QueryType) SetJoin(tablePtr interface{}, columnPtr interface{}, whereColumnPtr interface{}) {
	rec.setJoin(joinModeInner, tablePtr, "", columnPtr, whereColumnPtr)
}

func (rec *QueryType) SetJoinAs(tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{}) {
	rec.setJoin(joinModeInner, tablePtr, tableAs, columnPtr, whereColumnPtr)
}

func (rec *QueryType) SetJoinLeft(tablePtr interface{}, columnPtr interface{}, whereColumnPtr interface{}) {
	rec.setJoin(joinModeLeft, tablePtr, "", columnPtr, whereColumnPtr)
}

func (rec *QueryType) SetJoinLeftAs(tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{}) {
	rec.setJoin(joinModeLeft, tablePtr, tableAs, columnPtr, whereColumnPtr)
}

func (rec *QueryType) SetJoinRight(tablePtr interface{}, columnPtr interface{}, whereColumnPtr interface{}) {
	rec.setJoin(joinModeRight, tablePtr, "", columnPtr, whereColumnPtr)
}

func (rec *QueryType) SetJoinRightAs(tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{}) {
	rec.setJoin(joinModeRight, tablePtr, tableAs, columnPtr, whereColumnPtr)
}

func (rec *QueryType) setJoinWhere(mode int, prefix int, tablePtr interface{}, str string, columnPtr interface{}, valueList ...interface{}) {
	joinWhereData := &joinWhereType{
		Mode:      mode,
		Prefix:    prefix,
		TablePtr:  tablePtr,
		Str:       str,
		ColumnPtr: columnPtr,
		ValueList: valueList,
	}

	rec.JoinWhereList = append(rec.JoinWhereList, joinWhereData)
	rec.Data = nil
}

func (rec *QueryType) SetJoinWhereString(tablePtr interface{}, str string, valueList ...interface{}) {
	rec.setJoinWhere(queryModeString, queryPrefixAnd, tablePtr, str, nil, valueList...)
}

func (rec *QueryType) SetJoinWhereFormat(tablePtr interface{}, format string, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeFormat, queryPrefixAnd, tablePtr, format, columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereIs(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeIs, queryPrefixAnd, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereIsNot(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeIsNot, queryPrefixAnd, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereIsNull(tablePtr interface{}, columnPtr interface{}) {
	rec.setJoinWhere(queryModeIsNull, queryPrefixAnd, tablePtr, "", columnPtr)
}

func (rec *QueryType) SetJoinWhereIsNotNull(tablePtr interface{}, columnPtr interface{}) {
	rec.setJoinWhere(queryModeIsNullNot, queryPrefixAnd, tablePtr, "", columnPtr)
}

func (rec *QueryType) SetJoinWhereLike(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeLike, queryPrefixAnd, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereLikeNot(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeLikeNot, queryPrefixAnd, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereIn(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeIn, queryPrefixAnd, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereInNot(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeInNot, queryPrefixAnd, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereGt(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeGt, queryPrefixAnd, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereGte(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeGte, queryPrefixAnd, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereLt(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeLt, queryPrefixAnd, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereLte(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeLte, queryPrefixAnd, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrString(tablePtr interface{}, str string, valueList ...interface{}) {
	rec.setJoinWhere(queryModeString, queryPrefixOr, tablePtr, str, nil, valueList...)
}

func (rec *QueryType) SetJoinWhereOrFormat(tablePtr interface{}, format string, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeFormat, queryPrefixOr, tablePtr, format, columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrIs(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeIs, queryPrefixOr, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereOrIsNot(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeIsNot, queryPrefixOr, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereOrIsNull(tablePtr interface{}, columnPtr interface{}) {
	rec.setJoinWhere(queryModeIsNull, queryPrefixOr, tablePtr, "", columnPtr)
}

func (rec *QueryType) SetJoinWhereOrIsNullNot(tablePtr interface{}, columnPtr interface{}) {
	rec.setJoinWhere(queryModeIsNullNot, queryPrefixOr, tablePtr, "", columnPtr)
}

func (rec *QueryType) SetJoinWhereOrLike(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeLike, queryPrefixOr, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereOrLikeNot(tablePtr interface{}, columnPtr interface{}, value interface{}) {
	rec.setJoinWhere(queryModeLikeNot, queryPrefixOr, tablePtr, "", columnPtr, value)
}

func (rec *QueryType) SetJoinWhereOrIn(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeIn, queryPrefixOr, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrInNot(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeInNot, queryPrefixOr, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrGt(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeGt, queryPrefixOr, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrGte(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeGte, queryPrefixOr, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrLt(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeLt, queryPrefixOr, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereOrLte(tablePtr interface{}, columnPtr interface{}, valueList ...interface{}) {
	rec.setJoinWhere(queryModeLte, queryPrefixOr, tablePtr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetJoinWhereNest(tablePtr interface{}) {
	rec.setJoinWhere(queryModeNest, queryPrefixAnd, tablePtr, "", nil)
}

func (rec *QueryType) SetJoinWhereOrNest(tablePtr interface{}) {
	rec.setJoinWhere(queryModeNest, queryPrefixOr, tablePtr, "", nil)
}

func (rec *QueryType) SetJoinWhereNestClose(tablePtr interface{}) {
	rec.setJoinWhere(queryModeNestClose, queryPrefixNone, tablePtr, "", nil)
}

func (rec *QueryType) setSelect(mode int, str string, columnPtr interface{}, columnAs string) {
	selectData := &selectType{
		Mode:      mode,
		Str:       str,
		ColumnPtr: columnPtr,
		ColumnAs:  columnAs,
	}

	rec.SelectList = append(rec.SelectList, selectData)
	rec.Data = nil
}

func (rec *QueryType) SetSelectString(str string) {
	rec.setSelect(queryModeString, str, nil, "")
}

func (rec *QueryType) SetSelectStringAs(str string, as string) {
	rec.setSelect(queryModeString, str, nil, as)
}

func (rec *QueryType) SetSelectFormat(format string, columnPtr interface{}) {
	rec.setSelect(queryModeFormat, format, columnPtr, "")
}

func (rec *QueryType) SetSelectFormatAs(format string, columnPtr interface{}, as string) {
	rec.setSelect(queryModeFormat, format, columnPtr, as)
}

func (rec *QueryType) SetSelect(columnPtrList ...interface{}) {
	for _, columnPtr := range columnPtrList {
		rec.setSelect(queryModeOne, "", columnPtr, "")
	}
}

func (rec *QueryType) SetSelectAs(columnPtr interface{}, as string) {
	rec.setSelect(queryModeOne, "", columnPtr, as)
}

func (rec *QueryType) SetSelectAll(tablePtr interface{}) {
	rec.setSelect(queryModeAll, "", tablePtr, "")
}

func (rec *QueryType) SetSet(columnPtr interface{}, value interface{}) {
	setData := &setType{
		ColumnPtr: columnPtr,
		Value:     value,
	}

	rec.SetList = append(rec.SetList, setData)
	rec.Data = nil
}

func (rec *QueryType) SetValuesColumn(columnPtrList ...interface{}) {
	for _, columnPtr := range columnPtrList {
		valuesColumnData := &valuesColumnType{
			ColumnPtr: columnPtr,
		}

		rec.ValuesColumnList = append(rec.ValuesColumnList, valuesColumnData)
		rec.Data = nil
	}
}

func (rec *QueryType) SetValues(valueList ...interface{}) {
	if len(valueList) > 0 {
		var valList []*valuesType

		for _, value := range valueList {
			valuesData := &valuesType{
				Value: value,
			}

			valList = append(valList, valuesData)
		}

		rec.ValuesList = append(rec.ValuesList, valList)
		rec.Data = nil
	}
}

func (rec *QueryType) SetValuesClear() {
	rec.ValuesList = make([][]*valuesType, 0)
	rec.Data = nil
}

func (rec *QueryType) setWhere(mode int, prefix int, str string, columnPtr interface{}, valueList ...interface{}) {
	whereData := &whereType{
		Mode:      mode,
		Prefix:    prefix,
		Str:       str,
		ColumnPtr: columnPtr,
		ValueList: valueList,
	}

	rec.WhereList = append(rec.WhereList, whereData)
	rec.Data = nil
}

func (rec *QueryType) SetWhereString(str string, valueList ...interface{}) {
	rec.setWhere(queryModeString, queryPrefixAnd, str, nil, valueList...)
}

func (rec *QueryType) SetWhereFormat(format string, columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeFormat, queryPrefixAnd, format, columnPtr, valueList...)
}

func (rec *QueryType) SetWhereIs(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeIs, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetWhereIsNot(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeIsNot, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetWhereIsNull(columnPtr interface{}) {
	rec.setWhere(queryModeIsNull, queryPrefixAnd, "", columnPtr)
}

func (rec *QueryType) SetWhereIsNotNull(columnPtr interface{}) {
	rec.setWhere(queryModeIsNullNot, queryPrefixAnd, "", columnPtr)
}

func (rec *QueryType) SetWhereLike(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeLike, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetWhereLikeNot(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeLikeNot, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetWhereIn(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeIn, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereInNot(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeInNot, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereGt(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeGt, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereGte(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeGte, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereLt(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeLt, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereLte(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeLte, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrString(str string, valueList ...interface{}) {
	rec.setWhere(queryModeString, queryPrefixOr, str, nil, valueList...)
}

func (rec *QueryType) SetWhereOrFormat(format string, columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeFormat, queryPrefixOr, format, columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrIs(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeIs, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetWhereOrIsNot(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeIsNot, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetWhereOrIsNull(columnPtr interface{}) {
	rec.setWhere(queryModeIsNull, queryPrefixOr, "", columnPtr)
}

func (rec *QueryType) SetWhereOrIsNullNot(columnPtr interface{}) {
	rec.setWhere(queryModeIsNullNot, queryPrefixOr, "", columnPtr)
}

func (rec *QueryType) SetWhereOrLike(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeLike, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetWhereOrLikeNot(columnPtr interface{}, value interface{}) {
	rec.setWhere(queryModeLikeNot, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetWhereOrIn(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeIn, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrInNot(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeInNot, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrGt(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeGt, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrGte(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeGte, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrLt(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeLt, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereOrLte(columnPtr interface{}, valueList ...interface{}) {
	rec.setWhere(queryModeLte, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetWhereNest() {
	rec.setWhere(queryModeNest, queryPrefixAnd, "", nil)
}

func (rec *QueryType) SetWhereOrNest() {
	rec.setWhere(queryModeNest, queryPrefixOr, "", nil)
}

func (rec *QueryType) SetWhereNestClose() {
	rec.setWhere(queryModeNestClose, queryPrefixNone, "", nil)
}

func (rec *QueryType) setGroupBy(mode int, str string, columnPtr interface{}) {
	groupByData := &groupByType{
		Mode:      mode,
		Str:       str,
		ColumnPtr: columnPtr,
	}

	rec.GroupByList = append(rec.GroupByList, groupByData)
	rec.Data = nil
}

func (rec *QueryType) SetGroupBy(columnPtr interface{}) {
	rec.setGroupBy(queryModeOne, "", columnPtr)
}

func (rec *QueryType) SetGroupByString(str string) {
	rec.setGroupBy(queryModeString, str, nil)
}

func (rec *QueryType) SetGroupByFormat(format string, columnPtr interface{}) {
	rec.setGroupBy(queryModeFormat, format, columnPtr)
}

func (rec *QueryType) setHaving(mode int, prefix int, str string, columnPtr interface{}, valueList ...interface{}) {
	havingData := &havingType{
		Mode:      mode,
		Prefix:    prefix,
		Str:       str,
		ColumnPtr: columnPtr,
		ValueList: valueList,
	}

	rec.HavingList = append(rec.HavingList, havingData)
	rec.Data = nil
}

func (rec *QueryType) SetHavingString(str string, valueList ...interface{}) {
	rec.setHaving(queryModeString, queryPrefixAnd, str, nil, valueList...)
}

func (rec *QueryType) SetHavingFormat(format string, columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeFormat, queryPrefixAnd, format, columnPtr, valueList...)
}

func (rec *QueryType) SetHavingIs(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeIs, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetHavingIsNot(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeIsNot, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetHavingIsNull(columnPtr interface{}) {
	rec.setHaving(queryModeIsNull, queryPrefixAnd, "", columnPtr)
}

func (rec *QueryType) SetHavingIsNotNull(columnPtr interface{}) {
	rec.setHaving(queryModeIsNullNot, queryPrefixAnd, "", columnPtr)
}

func (rec *QueryType) SetHavingLike(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeLike, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetHavingLikeNot(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeLikeNot, queryPrefixAnd, "", columnPtr, value)
}

func (rec *QueryType) SetHavingIn(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeIn, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingInNot(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeInNot, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingGt(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeGt, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingGte(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeGte, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingLt(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeLt, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingLte(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeLte, queryPrefixAnd, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrString(str string, valueList ...interface{}) {
	rec.setHaving(queryModeString, queryPrefixOr, str, nil, valueList...)
}

func (rec *QueryType) SetHavingOrFormat(format string, columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeFormat, queryPrefixOr, format, columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrIs(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeIs, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetHavingOrIsNot(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeIsNot, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetHavingOrIsNull(columnPtr interface{}) {
	rec.setHaving(queryModeIsNull, queryPrefixOr, "", columnPtr)
}

func (rec *QueryType) SetHavingOrIsNullNot(columnPtr interface{}) {
	rec.setHaving(queryModeIsNullNot, queryPrefixOr, "", columnPtr)
}

func (rec *QueryType) SetHavingOrLike(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeLike, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetHavingOrLikeNot(columnPtr interface{}, value interface{}) {
	rec.setHaving(queryModeLikeNot, queryPrefixOr, "", columnPtr, value)
}

func (rec *QueryType) SetHavingOrIn(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeIn, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrInNot(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeInNot, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrGt(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeGt, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrGte(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeGte, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrLt(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeLt, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingOrLte(columnPtr interface{}, valueList ...interface{}) {
	rec.setHaving(queryModeLte, queryPrefixOr, "", columnPtr, valueList...)
}

func (rec *QueryType) SetHavingNest() {
	rec.setHaving(queryModeNest, queryPrefixAnd, "", nil)
}

func (rec *QueryType) SetHavingOrNest() {
	rec.setHaving(queryModeNest, queryPrefixOr, "", nil)
}

func (rec *QueryType) SetHavingNestClose() {
	rec.setHaving(queryModeNestClose, queryPrefixNone, "", nil)
}

func (rec *QueryType) setOrderBy(mode int, order int, str string, columnPtr interface{}) {
	orderByData := &orderByType{
		Mode:      mode,
		Order:     order,
		Str:       str,
		ColumnPtr: columnPtr,
	}

	rec.OrderByList = append(rec.OrderByList, orderByData)
	rec.Data = nil
}

func (rec *QueryType) SetOrderBy(columnPtr interface{}) {
	rec.SetOrderByAsc(columnPtr)
}

func (rec *QueryType) SetOrderByAsc(columnPtr interface{}) {
	rec.setOrderBy(queryModeOne, Asc, "", columnPtr)
}

func (rec *QueryType) SetOrderByAscString(str string) {
	rec.setOrderBy(queryModeString, Asc, str, nil)
}

func (rec *QueryType) SetOrderByAscFormat(format string, columnPtr interface{}) {
	rec.setOrderBy(queryModeFormat, Asc, format, columnPtr)
}

func (rec *QueryType) SetOrderByDesc(columnPtr interface{}) {
	rec.setOrderBy(queryModeOne, Desc, "", columnPtr)
}

func (rec *QueryType) SetOrderByDescString(str string) {
	rec.setOrderBy(queryModeString, Desc, str, nil)
}

func (rec *QueryType) SetOrderByDescFormat(format string, columnPtr interface{}) {
	rec.setOrderBy(queryModeFormat, Desc, format, columnPtr)
}

func (rec *QueryType) SetLimit(num int) {
	rec.Limit = num
	rec.Data = nil
}

func (rec *QueryType) SetOffset(num int) {
	rec.Offset = num
	rec.Data = nil
}

func (rec *QueryType) buildMeta() error {
	rec.MetaMap = make(map[string]*metaType, 0)
	rec.Data = &buildType{}
	if rec.modeResetAuto {
		rec.Reset()
	}

	_setMeta := func(tablePtr interface{}, tableAs string) error {
		tableType := reflect.TypeOf(tablePtr).Elem()
		tableVal := reflect.ValueOf(tablePtr).Elem()

		tableBase := toSnakeCase(tableType.Name())
		tableAsBase := tableAs

		numField := tableType.NumField()
		if numField < 1 {
			return errors.New("tablePtr none field")
		}

		for i := 0; i < tableType.NumField(); i++ {
			fieldType := tableType.Field(i)
			fieldVal := tableVal.FieldByName(fieldType.Name)

			columnBase := fieldType.Tag.Get(structFieldTagNameColumn)
			if columnBase == "" {
				continue
			}

			table, tableAs := rec.getTableName(tableBase, tableAsBase)
			tableColumn, tableAsColumn, column := rec.getColumnName(table, tableAs, columnBase)

			metaData := &metaType{
				TableBase:     tableBase,
				TableAsBase:   tableAsBase,
				Table:         table,
				TableAs:       tableAs,
				TableColumn:   tableColumn,
				TableAsColumn: tableAsColumn,
				ColumnBase:    columnBase,
				Column:        column,
			}

			key, err := getAddr(fieldVal)
			if err != nil {
				return err
			}

			rec.MetaMap[key] = metaData
		}

		return nil
	}

	if rec.Table != nil {
		err := _setMeta(rec.Table.TablePtr, rec.Table.TableAs)
		if err != nil {
			return err
		}
	}

	for _, table := range rec.JoinList {
		err := _setMeta(table.TablePtr, table.TableAs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rec *QueryType) buildValue(value interface{}) ([]string, error) {
	var strList []string

	val := reflect.ValueOf(value)
	kind := val.Kind()
	for kind == reflect.Interface || kind == reflect.Ptr {
		if kind == reflect.Ptr {
			elem := val.Elem()
			if !elem.IsValid() {
				break
			}
		}

		switch kind {
		case reflect.Ptr:
			val = val.Elem()
			kind = val.Kind()
		default:
			val = val.Elem()
			kind = val.Kind()
		}
	}

	var valList []interface{}
	if kind == reflect.Array || kind == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			valList = append(valList, val.Index(i).Interface())
		}
	} else if kind == reflect.Invalid {
		valList = append(valList, nil)
	} else {
		valList = append(valList, val.Interface())
	}

	for _, val := range valList {
		valValue := reflect.ValueOf(val)
		kind = valValue.Kind()
		if kind == reflect.Invalid {
			strList = append(strList, "NULL")
			continue
		} else if kind == reflect.Ptr {
			elem := valValue.Elem()
			if !elem.IsValid() {
				strList = append(strList, "NULL")
				continue
			}
		}

		rec.Data.ValueList = append(rec.Data.ValueList, val)
		strList = append(strList, rec.getPlaceholder())
	}

	return strList, nil
}

func (rec *QueryType) buildTable() error {
	var str string

	if rec.Table == nil {
		return errors.New("table not exist")
	}

	if rec.Table.Str != "" {
		rec.Data.Table = rec.Table.Str
		str = fmt.Sprintf("FROM %s", rec.Table.Str)
		rec.Data.TableForSelect = str
		return nil
	}

	addr, err := getAddrFromInterface(rec.Table.TablePtr)
	if err != nil {
		return err
	}

	meta, ok := rec.MetaMap[addr]
	if !ok {
		return errors.New("table meta not exist")
	}

	table := meta.Table

	rec.Data.Table = fmt.Sprintf("%s", table)

	if meta.TableAsBase != "" {
		table = fmt.Sprintf("%s as %s", str, meta.TableAs)
	}

	rec.Data.TableForSelect = fmt.Sprintf("FROM %s", table)

	return nil
}

func (rec *QueryType) buildJoin() error {
	var strList []string

	var joinWhereMap map[string][]*joinWhereType
	joinWhereMap = make(map[string][]*joinWhereType)
	for _, joinWhereData := range rec.JoinWhereList {
		addr, err := getAddrFromInterface(joinWhereData.TablePtr)
		if err != nil {
			return err
		}

		if _, ok := joinWhereMap[addr]; !ok {
			joinWhereMap[addr] = make([]*joinWhereType, 0)
		}

		joinWhereMap[addr] = append(joinWhereMap[addr], joinWhereData)
	}

	for _, joinData := range rec.JoinList {
		var joinWhereList []string

		var metaTable *metaType
		var addrTable string
		{
			addr, err := getAddrFromInterface(joinData.TablePtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("build join column1 meta not exist")
			}

			metaTable = meta
			addrTable = addr
		}

		var metaColumn *metaType
		{
			addr, err := getAddrFromInterface(joinData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("build join column1 meta not exist")
			}

			metaColumn = meta
		}

		var metaWhere *metaType
		{
			addr, err := getAddrFromInterface(joinData.WhereColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("build join column2 meta not exist")
			}

			metaWhere = meta
		}

		prefix := ""
		switch joinData.Mode {
		case joinModeInner:
			prefix = "INNER"
		case joinModeLeft:
			prefix = "LEFT"
		case joinModeRight:
			prefix = "RIGHT"
		default:
			return errors.New("join unknown mode")
		}

		type dataType struct {
			Meta  *metaType
			Base  string
			Value string
		}

		if _, ok := joinWhereMap[addrTable]; ok {
			prefixFlag := true
			for _, joinWhereData := range joinWhereMap[addrTable] {
				var strList []string

				data := &dataType{}

				if !isNil(joinWhereData.ColumnPtr) {
					addr, err := getAddrFromInterface(joinWhereData.ColumnPtr)
					if err != nil {
						return err
					}

					meta, ok := rec.MetaMap[addr]
					if !ok {
						return errors.New("joinWhere meta not exist")
					}

					data.Meta = meta
				}

				if len(joinWhereData.ValueList) > 0 {
					var strList []string

					for _, val := range joinWhereData.ValueList {
						valList, err := rec.buildValue(val)
						if err != nil {
							return err
						}

						strList = append(strList, valList...)
					}

					if len(strList) > 0 {
						data.Value = strings.Join(strList, ", ")
					}
				}

				if prefixFlag {
					switch joinWhereData.Prefix {
					case queryPrefixAnd:
						str := "AND"
						strList = append(strList, str)
					case queryPrefixOr:
						str := "OR"
						strList = append(strList, str)
					}
				} else {
					prefixFlag = true
				}

				switch joinWhereData.Mode {
				case queryModeString:
					data.Base = joinWhereData.Str
				case queryModeFormat:
					data.Base = joinWhereData.Str
				case queryModeIs:
					data.Base = "%s = %s"
				case queryModeIsNot:
					data.Base = "%s != %s"
				case queryModeIsNull:
					data.Base = "%s IS NULL"
				case queryModeIsNullNot:
					data.Base = "%s IS NOT NULL"
				case queryModeLike:
					data.Base = "%s LIKE %s"
				case queryModeLikeNot:
					data.Base = "%s NOT LIKE %s"
				case queryModeIn:
					data.Base = "%s IN (%s)"
				case queryModeInNot:
					data.Base = "%s NOT IN (%s)"
				case queryModeGt:
					data.Base = "%s > %s"
				case queryModeGte:
					data.Base = "%s >= %s"
				case queryModeLt:
					data.Base = "%s < %s"
				case queryModeLte:
					data.Base = "%s <= %s"
				case queryModeNest:
					data.Base = "("
					prefixFlag = false
				case queryModeNestClose:
					data.Base = ")"
				default:
					return errors.New("joinWhere mode not exist")
				}

				if data.Meta != nil {
					if data.Value != "" {
						strList = append(strList, fmt.Sprintf(data.Base, data.Meta.TableAsColumn, data.Value))
					} else {
						strList = append(strList, fmt.Sprintf(data.Base, data.Meta.TableAsColumn))
					}
				} else {
					if data.Value != "" {
						strList = append(strList, fmt.Sprintf(data.Base, data.Value))
					} else {
						strList = append(strList, data.Base)
					}
				}

				joinWhereList = append(joinWhereList, strings.Join(strList, " "))
			}
		}

		where := metaWhere.TableAsColumn
		if len(joinWhereList) > 0 {
			str := strings.Join(joinWhereList, " ")
			where = fmt.Sprintf("%s %s", where, str)
		}

		table := metaTable.Table
		if metaTable.Table != metaTable.TableAs {
			table = fmt.Sprintf("%s as %s", table, metaTable.TableAs)
		}

		str := fmt.Sprintf("%s JOIN %s ON %s = %s", prefix, table, metaColumn.TableAsColumn, where)

		strList = append(strList, str)
	}

	if len(strList) > 0 {
		rec.Data.Join = strings.Join(strList, " ")
	}

	return nil
}

func (rec *QueryType) buildSelect() error {
	var strList []string

	type dataType struct {
		Meta *metaType
	}

	for _, selectData := range rec.SelectList {
		data := dataType{}

		if !isNil(selectData.ColumnPtr) {
			addr, err := getAddrFromInterface(selectData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("select column meta not exist")
			}

			data.Meta = meta
		}

		switch selectData.Mode {
		case queryModeOne:
			if data.Meta == nil {
				return errors.New("select column meta not exist")
			}
			str := data.Meta.TableAsColumn
			if selectData.ColumnAs != "" {
				str = fmt.Sprintf("%s as \"%s\"", str, selectData.ColumnAs)
			}
			strList = append(strList, str)
		case queryModeAll:
			if data.Meta == nil {
				return errors.New("select column table not exist")
			}
			str := fmt.Sprintf("%s.*", data.Meta.TableAs)
			strList = append(strList, str)
		case queryModeString:
			strList = append(strList, selectData.Str)
		case queryModeFormat:
			str := selectData.Str
			if data.Meta != nil {
				str = fmt.Sprintf(str, data.Meta.TableAsColumn)
			}
			if selectData.ColumnAs != "" {
				str = fmt.Sprintf("%s as \"%s\"", str, selectData.ColumnAs)
			}
			strList = append(strList, str)
		default:
			return errors.New("select mode not exist")
		}
	}

	if len(strList) > 0 {
		rec.Data.Select = fmt.Sprintf("SELECT %s", strings.Join(strList, ", "))
	}

	return nil
}

func (rec *QueryType) buildValuesColumnAndValues() error {
	if len(rec.ValuesColumnList) < 1 {
		return nil
	}

	var valuesColumnCount int

	{
		var valuesColumnList []string

		for _, valuesColumnData := range rec.ValuesColumnList {
			addr, err := getAddrFromInterface(valuesColumnData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("valuesColumn meta not exist")
			}

			valuesColumnList = append(valuesColumnList, meta.Column)
		}

		valuesColumnCount = len(valuesColumnList)
		if len(valuesColumnList) > 0 {
			str := strings.Join(valuesColumnList, ", ")
			rec.Data.ValuesColumn = fmt.Sprintf("(%s)", str)
		}
	}

	{
		var valuesList []string

		for _, valuesDataList := range rec.ValuesList {
			if len(valuesDataList) != valuesColumnCount {
				return errors.New("values is different from the number of valuescolumn")
			}

			var strList []string

			for _, valuesData := range valuesDataList {
				valList, err := rec.buildValue(valuesData.Value)
				if err != nil {
					return err
				}

				strList = append(strList, valList...)
			}

			if len(strList) > 0 {
				str := strings.Join(strList, ", ")
				valuesList = append(valuesList, fmt.Sprintf("(%s)", str))
			}
		}

		if len(valuesList) > 0 {
			str := strings.Join(valuesList, ", ")
			rec.Data.Values = fmt.Sprintf("VALUES %s", str)
		}
	}

	return nil
}

func (rec *QueryType) buildSet() error {
	var setList []string

	for _, setData := range rec.SetList {
		addr, err := getAddrFromInterface(setData.ColumnPtr)
		if err != nil {
			return err
		}

		meta, ok := rec.MetaMap[addr]
		if !ok {
			return errors.New("set meta not exist")
		}

		valList, err := rec.buildValue(setData.Value)
		if err != nil {
			return err
		}

		if len(valList) != 1 {
			return errors.New("set value length is not 1")
		}

		setList = append(setList, fmt.Sprintf("%s = %v", meta.Column, valList[0]))
	}

	if len(setList) > 0 {
		rec.Data.Set = fmt.Sprintf("SET %s", strings.Join(setList, ", "))
	}

	return nil
}

func (rec *QueryType) buildWhere() error {
	var whereList []string
	var whereForSelectList []string

	type dataType struct {
		Meta  *metaType
		Base  string
		Value string
	}
	prefixFlag := false

	for _, whereData := range rec.WhereList {
		var strList []string
		var strForSelectList []string

		data := &dataType{}

		if !isNil(whereData.ColumnPtr) {
			addr, err := getAddrFromInterface(whereData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("where meta not exist")
			}

			data.Meta = meta
		}

		if len(whereData.ValueList) > 0 {
			var strList []string

			for _, val := range whereData.ValueList {
				valList, err := rec.buildValue(val)
				if err != nil {
					return err
				}

				strList = append(strList, valList...)
			}

			if len(strList) > 0 {
				data.Value = strings.Join(strList, ", ")
			}
		}

		if prefixFlag {
			switch whereData.Prefix {
			case queryPrefixAnd:
				str := "AND"
				strList = append(strList, str)
				strForSelectList = append(strForSelectList, str)
			case queryPrefixOr:
				str := "OR"
				strList = append(strList, str)
				strForSelectList = append(strForSelectList, str)
			}
		} else {
			prefixFlag = true
		}

		switch whereData.Mode {
		case queryModeString:
			data.Base = whereData.Str
		case queryModeFormat:
			data.Base = whereData.Str
		case queryModeIs:
			data.Base = "%s = %s"
		case queryModeIsNot:
			data.Base = "%s != %s"
		case queryModeIsNull:
			data.Base = "%s IS NULL"
		case queryModeIsNullNot:
			data.Base = "%s IS NOT NULL"
		case queryModeLike:
			data.Base = "%s LIKE %s"
		case queryModeLikeNot:
			data.Base = "%s NOT LIKE %s"
		case queryModeIn:
			data.Base = "%s IN (%s)"
		case queryModeInNot:
			data.Base = "%s NOT IN (%s)"
		case queryModeGt:
			data.Base = "%s > %s"
		case queryModeGte:
			data.Base = "%s >= %s"
		case queryModeLt:
			data.Base = "%s < %s"
		case queryModeLte:
			data.Base = "%s <= %s"
		case queryModeNest:
			data.Base = "("
			prefixFlag = false
		case queryModeNestClose:
			data.Base = ")"
		default:
			return errors.New("where type not exist")
		}

		if data.Meta != nil {
			if data.Value != "" {
				strList = append(strList, fmt.Sprintf(data.Base, data.Meta.Column, data.Value))
				strForSelectList = append(strForSelectList, fmt.Sprintf(data.Base, data.Meta.TableAsColumn, data.Value))
			} else {
				strList = append(strList, fmt.Sprintf(data.Base, data.Meta.Column))
				strForSelectList = append(strForSelectList, fmt.Sprintf(data.Base, data.Meta.TableAsColumn))
			}
		} else {
			if data.Value != "" {
				strList = append(strList, fmt.Sprintf(data.Base, data.Value))
				strForSelectList = append(strForSelectList, fmt.Sprintf(data.Base, data.Value))
			} else {
				strList = append(strList, data.Base)
				strForSelectList = append(strForSelectList, data.Base)
			}
		}

		whereList = append(whereList, strings.Join(strList, " "))
		whereForSelectList = append(whereForSelectList, strings.Join(strForSelectList, " "))
	}

	if len(whereList) > 0 {
		rec.Data.Where = fmt.Sprintf("WHERE %s", strings.Join(whereList, " "))
	}

	if len(whereForSelectList) > 0 {
		rec.Data.WhereForSelect = fmt.Sprintf("WHERE %s", strings.Join(whereForSelectList, " "))
	}

	return nil
}

func (rec *QueryType) buildGroupBy() error {
	var groupByList []string

	type dataType struct {
		Meta *metaType
		Base string
	}

	for _, groupByData := range rec.GroupByList {
		data := &dataType{}

		if !isNil(groupByData) {
			addr, err := getAddrFromInterface(groupByData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("group by build meta not exist")
			}

			data.Meta = meta
		}

		switch groupByData.Mode {
		case queryModeOne:
			data.Base = "%s"
		case queryModeString:
			data.Base = groupByData.Str
		case queryModeFormat:
			data.Base = groupByData.Str
		default:
			return errors.New("groupBy mode not exist")
		}

		str := data.Base
		if data.Meta != nil {
			str = fmt.Sprintf(str, data.Meta.Column)
		}

		groupByList = append(groupByList, str)
	}

	if len(groupByList) > 0 {
		rec.Data.GroupBy = fmt.Sprintf("GROUP BY %s", strings.Join(groupByList, ", "))
	}

	return nil
}

func (rec *QueryType) buildHaving() error {
	var havingList []string

	type dataType struct {
		Meta  *metaType
		Base  string
		Value string
	}
	prefixFlag := false

	for _, havingData := range rec.HavingList {
		var strList []string

		data := dataType{}

		if !isNil(havingData.ColumnPtr) {
			addr, err := getAddrFromInterface(havingData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("having meta not exist")
			}

			data.Meta = meta
		}

		if len(havingData.ValueList) > 0 {
			var strList []string

			for _, val := range havingData.ValueList {
				valList, err := rec.buildValue(val)
				if err != nil {
					return err
				}

				strList = append(strList, valList...)
			}

			if len(strList) > 0 {
				data.Value = strings.Join(strList, ", ")
			}
		}

		if prefixFlag {
			switch havingData.Prefix {
			case queryPrefixAnd:
				str := "AND"
				strList = append(strList, str)
			case queryPrefixOr:
				str := "OR"
				strList = append(strList, str)
			}
		} else {
			prefixFlag = true
		}

		{
			switch havingData.Mode {
			case queryModeString:
				data.Base = havingData.Str
			case queryModeFormat:
				data.Base = havingData.Str
			case queryModeIs:
				data.Base = "%s = %s"
			case queryModeIsNot:
				data.Base = "%s != %s"
			case queryModeIsNull:
				data.Base = "%s IS NULL"
			case queryModeIsNullNot:
				data.Base = "%s IS NOT NULL"
			case queryModeLike:
				data.Base = "%s LIKE %s"
			case queryModeLikeNot:
				data.Base = "%s NOT LIKE %s"
			case queryModeIn:
				data.Base = "%s IN (%s)"
			case queryModeInNot:
				data.Base = "%s NOT IN (%s)"
			case queryModeGt:
				data.Base = "%s > %s"
			case queryModeGte:
				data.Base = "%s >= %s"
			case queryModeLt:
				data.Base = "%s < %s"
			case queryModeLte:
				data.Base = "%s <= %s"
			case queryModeNest:
				data.Base = "("
				prefixFlag = false
			case queryModeNestClose:
				data.Base = ")"
			default:
				return errors.New("having type not exist")
			}
		}

		if data.Meta != nil {
			if data.Value != "" {
				strList = append(strList, fmt.Sprintf(data.Base, data.Meta.TableAsColumn, data.Value))
			} else {
				strList = append(strList, fmt.Sprintf(data.Base, data.Meta.TableAsColumn))
			}
		} else {
			if data.Value != "" {
				strList = append(strList, fmt.Sprintf(data.Base, data.Value))
			} else {
				strList = append(strList, data.Base)
			}
		}

		havingList = append(havingList, strings.Join(strList, " "))
	}

	if len(havingList) > 0 {
		rec.Data.Having = fmt.Sprintf("HAVING %s", strings.Join(havingList, " "))
	}

	return nil
}

func (rec *QueryType) buildOrderBy() error {
	var orderByList []string

	type dataType struct {
		Meta *metaType
		Base string
	}

	for _, orderByData := range rec.OrderByList {
		data := &dataType{}

		if !isNil(orderByData.ColumnPtr) {
			addr, err := getAddrFromInterface(orderByData.ColumnPtr)
			if err != nil {
				return err
			}

			meta, ok := rec.MetaMap[addr]
			if !ok {
				return errors.New("build order meta not exist")
			}

			data.Meta = meta
		}

		switch orderByData.Mode {
		case queryModeOne:
			data.Base = "%s"
		case queryModeString:
			data.Base = orderByData.Str
		case queryModeFormat:
			data.Base = orderByData.Str
		default:
			return errors.New("orderBy mode not exist")
		}

		str := data.Base
		if data.Meta != nil {
			str = fmt.Sprintf(str, data.Meta.TableAsColumn)
		}
		if orderByData.Order == Desc {
			str = fmt.Sprintf("%s DESC", str)
		}

		orderByList = append(orderByList, str)
	}

	if len(orderByList) > 0 {
		rec.Data.Order = fmt.Sprintf("ORDER BY %s", strings.Join(orderByList, ", "))
	}

	return nil
}

func (rec *QueryType) buildLimit() error {
	if rec.Limit > 0 {
		rec.Data.Limit = fmt.Sprintf("LIMIT %v", rec.Limit)
	}

	return nil
}

func (rec *QueryType) buildOffset() error {
	if rec.Offset > 0 {
		rec.Data.Offset = fmt.Sprintf("OFFSET %v", rec.Offset)
	}

	return nil
}

func (rec *QueryType) GetSelectQuery() (string, []interface{}, error) {
	var err error
	query := ""

	err = rec.buildMeta()
	if err != nil {
		return "", nil, err
	}

	{
		err = rec.buildSelect()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Select
		if str == "" {
			return "", nil, errors.New("select not exist")
		}
		query = rec.Data.Select
	}

	{
		err = rec.buildTable()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.TableForSelect
		if str == "" {
			return "", nil, errors.New("select table not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	{
		err = rec.buildJoin()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Join
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildWhere()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.WhereForSelect
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildGroupBy()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.GroupBy
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildHaving()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Having
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildOrderBy()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Order
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildLimit()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Limit
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildOffset()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Offset
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	valueList := rec.Data.ValueList

	return query, valueList, nil
}

func (rec *QueryType) GetSelectCountQuery() (string, []interface{}, error) {
	var err error
	query := ""

	err = rec.buildMeta()
	if err != nil {
		return "", nil, err
	}

	query = "SELECT count(*)"

	{
		err = rec.buildTable()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.TableForSelect
		if str == "" {
			return "", nil, errors.New("select table not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	{
		err = rec.buildJoin()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Join
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildWhere()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.WhereForSelect
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildGroupBy()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.GroupBy
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	{
		err = rec.buildHaving()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Having
		if str != "" {
			query = fmt.Sprintf("%s %s", query, str)
		}
	}

	valueList := rec.Data.ValueList

	return query, valueList, nil
}

func (rec *QueryType) GetInsertQuery() (string, []interface{}, error) {
	var err error
	query := ""

	err = rec.buildMeta()
	if err != nil {
		return "", nil, err
	}

	query = "INSERT INTO"

	{
		err = rec.buildTable()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Table
		if str == "" {
			return "", nil, errors.New("table not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	{
		err = rec.buildValuesColumnAndValues()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.ValuesColumn
		if str == "" {
			return "", nil, errors.New("valuesColumn not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)

		str = rec.Data.Values
		if str == "" {
			return "", nil, errors.New("values not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	valueList := rec.Data.ValueList

	return query, valueList, nil
}

func (rec *QueryType) GetUpdateQuery() (string, []interface{}, error) {
	var err error
	query := ""

	err = rec.buildMeta()
	if err != nil {
		return "", nil, err
	}

	query = "UPDATE"

	{
		err = rec.buildTable()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Table
		if str == "" {
			return "", nil, errors.New("table not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	{
		err = rec.buildSet()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Set
		if str == "" {
			return "", nil, errors.New("set not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	{
		err = rec.buildWhere()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Where
		if str == "" {
			return "", nil, errors.New("where not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	valueList := rec.Data.ValueList

	return query, valueList, nil
}

func (rec *QueryType) GetDeleteQuery() (string, []interface{}, error) {
	var err error
	query := ""

	err = rec.buildMeta()
	if err != nil {
		return "", nil, err
	}

	query = "DELETE FROM"

	{
		err = rec.buildTable()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Table
		if str == "" {
			return "", nil, errors.New("table not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	{
		err = rec.buildWhere()
		if err != nil {
			return "", nil, err
		}

		str := rec.Data.Where
		if str == "" {
			return "", nil, errors.New("where not exist")
		}
		query = fmt.Sprintf("%s %s", query, str)
	}

	valueList := rec.Data.ValueList

	return query, valueList, nil
}

func (rec *QueryType) Exec(query string, valueList ...interface{}) (sql.Result, error) {
	var err error
	var result sql.Result

	if rec.DB == nil && rec.TX == nil {
		return nil, errors.New("database is null")
	}

	if rec.modeLog {
		fmt.Printf("query: %v\n", query)
		fmt.Printf("value: %v\n", valueList)
	}

	if rec.TX != nil {
		result, err = rec.TX.Exec(query, valueList...)
	} else {
		result, err = rec.DB.Exec(query, valueList...)
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rec *QueryType) ExecQuery(dest interface{}, query string, valueList ...interface{}) error {
	var err error
	var rows *sql.Rows

	if rec.DB == nil && rec.TX == nil {
		return errors.New("database is null")
	}

	if rec.modeLog {
		fmt.Printf("query: %v\n", query)
		fmt.Printf("value: %v\n", valueList)
	}

	if rec.TX != nil {
		rows, err = rec.TX.Query(query, valueList...)
	} else {
		rows, err = rec.DB.Query(query, valueList...)
	}
	if err != nil {
		return err
	}
	defer func() {
		_ = rows.Close()
	}()

	columnList, err := rows.Columns()
	if err != nil {
		return err
	}

	destValue := reflect.ValueOf(dest)
	destValueType := destValue.Type().String()
	{
		elem := destValue
		if elem.Kind() != reflect.Ptr {
			return errors.New("dest is not pointer. Should be type *[]struct")
		}
		elem = destValue.Elem()
		if elem.Kind() != reflect.Slice {
			return errors.New("dest is not pointer slice. Should be type *[]struct or *[]map[string]interface {}")
		}
	}
	destDirect := reflect.Indirect(destValue)

	destType := destValue.Type()
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}

	base := destType.Elem()

	switch base.Kind() {
	case reflect.Struct:
		var tagIndexMap map[string][]int
		{
			tagIndexMap, err = makeTagIndexMap(base, structFieldTagNameColumn)
			if err != nil {
				return err
			}

			if len(tagIndexMap) != len(columnList) {
				tagIndexMapFlag := true
				for _, column := range columnList {
					 _, ok := tagIndexMap[column]
					 if !ok {
						  tagIndexMapFlag = true
						  break
					 }
					 tagIndexMapFlag = false
				}
				if tagIndexMapFlag {
					 baseValue := reflect.New(base)
					 val := reflect.Indirect(baseValue)
					 if val.NumField() != len(columnList) {
						  return errors.New("length does not match")
					 }
					 for key, column := range columnList {
						  tagIndexMap[column] = []int{key}
					 }
				}
			}
		}

		scanList := make([]interface{}, len(columnList))
		for rows.Next() {
			baseValue := reflect.New(base)
			val := reflect.Indirect(baseValue)

			for key, column := range columnList {
				indexList, ok := tagIndexMap[column]
				if !ok {
					return errors.New("column not exist")
				}

				field := val
				for _, index := range indexList {
					field = field.Field(index)
				}

				scanList[key] = field.Addr().Interface()
			}

			err = rows.Scan(scanList...)
			if err != nil {
				return err
			}

			destDirect.Set(reflect.Append(destDirect, val))
		}
	default:
		if destValueType == "*[]map[string]interface {}" {
			columnChangeMap := make(map[string]string)
			if rec.modeResultKey == resultKeyModeCamelCase {
				for _, val := range columnList {
					columnChangeMap[val] = toCamelCase(val)
				}
			} else if rec.modeResultKey == resultKeyModeSnakeCase {
				for _, val := range columnList {
					columnChangeMap[val] = toSnakeCase(val)
				}
			} else {
				for _, val := range columnList {
					columnChangeMap[val] = val
				}
			}

			var valList = make([]interface{}, len(columnList))
			var scanList = make([]interface{}, len(columnList))
			for key, _ := range columnList {
				scanList[key] = &valList[key]
			}

			for rows.Next() {
				err = rows.Scan(scanList...)
				if err != nil {
					return err
				}

				scanMap := make(map[string]interface{})
				for key, column := range columnList {
					name := columnChangeMap[column]
					scanMap[name] = valList[key]
				}

				scanMapType := reflect.ValueOf(scanMap)
				destDirect.Set(reflect.Append(destDirect, scanMapType))
			}
		} else {
			return errors.New("type *[]struct or *[]map[string]interface{}")
		}
	}

	return nil
}

func (rec *QueryType) Select(dest interface{}) error {
	query, valueList, err := rec.GetSelectQuery()
	if err != nil {
		return err
	}

	err = rec.ExecQuery(dest, query, valueList...)
	if err != nil {
		return err
	}

	return nil
}

func (rec *QueryType) SelectCount(dest interface{}) error {
	query, valueList, err := rec.GetSelectCountQuery()
	if err != nil {
		return err
	}

	err = rec.ExecQuery(dest, query, valueList...)
	if err != nil {
		return err
	}

	return nil
}

func (rec *QueryType) Insert() (sql.Result, error) {
	query, valueList, err := rec.GetInsertQuery()
	if err != nil {
		return nil, err
	}

	result, err := rec.Exec(query, valueList...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rec *QueryType) Update() (sql.Result, error) {
	query, valueList, err := rec.GetUpdateQuery()
	if err != nil {
		return nil, err
	}

	result, err := rec.Exec(query, valueList...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rec *QueryType) Delete() (sql.Result, error) {
	query, valueList, err := rec.GetDeleteQuery()
	if err != nil {
		return nil, err
	}

	result, err := rec.Exec(query, valueList...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
