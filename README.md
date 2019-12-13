Set.*Format系メソッドは複数の値を入れる事が出来ないバグがあるかも


``` go
import (
  "github.com/Go-Go-LAND/gol"
)

func MyOpen() (*gol.DB, error) {
  ///// database config
  databaseType := gol.DatabaseTypePostgresql
  // databaseType := gol.DatabaseTypeMysql
  host := "localhost"
  port := "5432"
  user := "username"
  pass := "password"
  database := "database"
  optionMap := map[string]string{}

  ///// open
  # postgresql
  db, err := gol.Open(gol.DatabaseTypePostgresql, host, port, user, pass, database, optionMap)
  if err != nil {
    return nil, err
  }

  db.SetModeLog(false)

  return db, nil
}
```


# sample struct
```
// User > user table
type User struct {
  Id        int           `column:"id" json:"id"`
  CreatedAt time.Time     `column:"created_at" json:"createdAt"`
  CreatedBy gol.NullInt64 `column:"created_by" json:"createdBy"`
  UpdatedAt time.Time     `column:"updated_at" json:"updatedAt"`
  UpdatedBy gol.NullInt64 `column:"updated_by" json:"updatedBy"`
  DeletedAt gol.NullTime  `column:"deleted_at" json:"deletedAt"`
  DeletedBy gol.NullInt64 `column:"deleted_by" json:"deletedBy"`
  Uid       string        `column:"uid" json:"uid"`
}

// UserDetail > user_detail table
type UserDetail struct {
  Id        int           `column:"id" json:"id"`
  CreatedAt time.Time     `column:"created_at" json:"createdAt"`
  CreatedBy gol.NullInt64 `column:"created_by" json:"createdBy"`
  UpdatedAt time.Time     `column:"updated_at" json:"updatedAt"`
  UpdatedBy gol.NullInt64 `column:"updated_by" json:"updatedBy"`
  DeletedAt gol.NullTime  `column:"deleted_at" json:"deletedAt"`
  DeletedBy gol.NullInt64 `column:"deleted_by" json:"deletedBy"`
  UserId    int           `column:"user_id" json:"userId"`
  Mail      string        `column:"mail" json:"mail"`
  Name      string        `column:"name" json:"name"`
}
```

# null
## If null is allowed, use the structure written in nullTypes.go file.
- NullBool
- NullInt32
- NullInt64
- NullFloat
- NullString
- NullTime


## The following methods are provided to change the value
- Get()
- Set(value)
- Delete()


# transaction
``` go
func Sample() error {
  db, err := MyOpen()
  if err != nil {
    return err
  }
  defer func() {
    if p := recover(); p != nil {
      _ = db.Close()
      panic(p)
    }
    _ = db.Close()
  }()

  err = func() error {
    tx, err := db.Begin()
    if err != nil {
      return err
    }
    defer func() {
      if p := recover(); p != nil {
        _ = tx.Rollback()
        panic(p)
      }
      _ = tx.Rollback()
    }()

    // query...

    return tx.Commit()
  }()
  if err != nil {
    return err
  }

  return nil
}
```

# select
``` go
var resultList []User{}

table := User{}
query := tx.Query()
query.SetTable(&table)
query.SetSelectAll(&table)
query.SetWhereIs(&table.Id, data.Id)
err = query.Select(&resultList)
if err != nil {
  return err
}
```

# select result map
``` go
var resultList []map[string]interface{}
table := User{}
query := tx.Query()
query.SetTable(&table)
query.SetSelectAll(&table)
query.SetWhereIs(&table.Id, data.Id)
err = query.Select(&resultList)
if err != nil {
  return err
}
```

# select join
``` go
var resultList []struct{
    User
    Mail `column:"mail" json:"mail"`
    Name `column:"name" json:"name"`
}
table := User{}
tableDetail := UserDetail{}
query := tx.Query()
query.SetTable(&table)
query.SetJoin(&tableDetail, &tableDetail.UserId, &table.Id)
query.SetSelectAll(&table)
query.SetSelect(
  &table.Mail,
  &table.Name,
)
query.SetWhereIs(&table.Id, data.Id)
err = query.Select(&resultList)
if err != nil {
  return err
}
```


# insert
``` go
userId := 1
now := time.Now()

data := User{}
data.CreatedAt = now
data.CreatedBy.Set(userId)
data.UpdatedAt = now
data.UpdatedBy.Set(userId)
data.DeletedAt.Delete()
data.DeletedBy.Delete()
data.Uid = "sample"


table := User{}
query := tx.Query()
query.SetTable(&table)
query.SetValuesColumn(
  &table.CreatedAt,
  &table.CreatedBy,
  &table.UpdatedAt,
  &table.UpdatedBy,
  &table.DeletedAt,
  &table.DeletedBy,
  &table.Uid,
)
query.SetValues(
  data.CreatedAt,
  data.CreatedBy,
  data.UpdatedAt,
  data.UpdatedBy,
  data.DeletedAt,
  data.DeletedBy,
  data.Uid,
)
query.SetWhereIs(&table.UserId, data.Id)
_, err = query.Insert()
if err != nil {
  return err
}
```

# update
``` go
userId := 1
now := time.Now()

data := User{}
data.UpdatedAt.Set(now)
data.UpdatedBy.Set(userId)
data.Uid = "sample"

table := User{}
query := tx.Query()
query.SetTable(&table)
query.SetSet(&table.UpdatedAt, data.UpdatedAt)
query.SetSet(&table.UpdatedBy, data.UpdatedBy)
query.SetSet(&table.Uid, data.Uid)
query.SetWhereIs(&table.UserId, data.Id)
_, err = query.Update()
if err != nil {
  return err
}
```

# delete
``` go
id = 1

table := User{}
query := tx.Query()
query.SetTable(&table)
query.SetWhereIs(&table.Id, id)
_, err = query.Delete()
if err != nil {
  return err
}
```


# queryType

# table

|method|sql|
|---|---|
|SetTable(tablePtr interface{})|FROM tablePtr|
|SetTableAs(tablePtr interface{}, tableAs string)|FROM tablePtr as tableAs|

# join
|method|sql|
|---|---|
|SetJoin(tablePtr interface{}, columnPtr interface{}, whereColumnPtr interface{})|JOIN tablePtr ON tablePtr = whereColumnPtr|
|SetJoinAs(tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{})|JOIN tablePtr ON columnPtr = whereColumnPtr|
|SetJoinLeft(tablePtr interface{}, columnPtr interface{}, whereColumnPtr interface{})|LEFT JOIN tablePtr ON columnPtr = whereColumnPtr|
|SetJoinLeftAs(tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{})|LEFT JOIN tablePtr as tableAs ON columnPtr = whereColumnPtr|
|SetJoinRight(tablePtr interface{}, columnPtr interface{}, whereColumnPtr interface{})|RIGHT JOIN tablePtr ON columnPtr = whereColumnPtr|
|SetJoinRightAs(tablePtr interface{}, tableAs string, columnPtr interface{}, whereColumnPtr interface{})|RIGHT JOIN tablePtr ON columnPtr = whereColumnPtr|


# join where
SetJoinWhere.+ は SetWhere.+ Join-On句の中に書かれるwhere句でほぼ同じなので省略


# select
|method|sql|
|---|---|
|SetSelectString(str string)|SELECT str|
|SetSelectStringAs(str string, as string)|SELECT str AS as|
|SetSelectFormat(format string, columnPtr interface{})|SELECT format AS as|
|SetSelectFormatAs(format string, columnPtr interface{}, as string)|SELECT format AS as|
|SetSelect(columnPtrList ...interface{})|SELECT columnPtrList...|
|SetSelectAs(columnPtr interface{}, as string)|SELECT columnPtr AS as|
|SetSelectAll(tablePtr interface{})|SELECT tablePtr.*|


# set
|method|sql|
|---|---|
|SetSet(columnPtr interface{}, value interface{})|SET columnPtr = value|


# insert into
|method|sql|
|---|---|
|SetValuesColumn(columnPtrList ...interface{})|INTO ? (columnPtrList...)|
|SetValues(valueList ...interface{})|VALUES (valueList...)|

SetValuesClear() is values clear


# where
|method|sql|
|---|---|
|SetWhereString(str string, valueList ...interface{})|WHERE [and] columnPtr < ?|
|SetWhereFormat(format string, columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr < ?|
|SetWhereIs(columnPtr interface{}, value interface{})|WHERE [and] columnPtr = ?|
|SetWhereIsNot(columnPtr interface{}, value interface{})|WHERE [and] columnPtr IS NOT ?|
|SetWhereIsNull(columnPtr interface{})|WHERE [and] columnPtr IS NULL|
|SetWhereIsNotNull(columnPtr interface{})|WHERE [and] columnPtr IS NOT NULL|
|SetWhereLike(columnPtr interface{}, value interface{})|WHERE [and] columnPtr LIKE ?|
|SetWhereLikeNot(columnPtr interface{}, value interface{})|WHERE [and] columnPtr NOT LIKE ?|
|SetWhereIn(columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr IN (?)|
|SetWhereInNot(columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr NOT IN (?)|
|SetWhereGt(columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr > ?|
|SetWhereGte(columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr >= ?|
|SetWhereLt(columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr < ?|
|SetWhereLte(columnPtr interface{}, valueList ...interface{})|WHERE [and] columnPtr <= ?|
|SetWhereOrString(str string, valueList ...interface{})|WHERE [or] str|
|SetWhereOrFormat(format string, columnPtr interface{}, valueList ...interface{})|WHERE [or] format|
|SetWhereOrIs(columnPtr interface{}, value interface{})|WHERE [or] columnPtr = ?|
|SetWhereOrIsNot(columnPtr interface{}, value interface{})|WHERE [or] columnPtr IS NOT ?|
|SetWhereOrIsNull(columnPtr interface{})|WHERE [or] columnPtr IS NULL|
|SetWhereOrIsNullNot(columnPtr interface{})|WHERE [or] columnPtr IS NOT NULL|
|SetWhereOrLike(columnPtr interface{}, value interface{})|WHERE [or] columnPtr LIKE ?|
|SetWhereOrLikeNot(columnPtr interface{}, value interface{})|WHERE [or] columnPtr NOT LIKE ?|
|SetWhereOrIn(columnPtr interface{}, valueList ...interface{})|WHERE [or] columnPtr IN (?)|
|SetWhereOrInNot(columnPtr interface{}, valueList ...interface{})|WHERE [or] columnPtr NOT IN (?)|
|SetWhereOrGt(columnPtr interface{}, valueList ...interface{})|WHERE [or] columnPtr > ?|
|SetWhereOrGte(columnPtr interface{}, valueList ...interface{})|WHERE [or] columnPtr >= ?|
|SetWhereOrLt(columnPtr interface{}, valueList ...interface{})|WHERE [or] columnPtr < ?|
|SetWhereOrLte(columnPtr interface{}, valueList ...interface{})|WHERE [or] columnPtr =< ?|
|SetWhereNest()|WHERE ? [and] (|
|SetWhereOrNest()|WHERE ? [or] (|
|SetWhereNestClose()|WHERE ? )|


# group by
|method|sql|
|---|---|
|SetGroupBy(columnPtr interface{})|GROUP BY columnPtr|
|SetGroupByString(str string)|GROUP BY str|
|SetGroupByFormat(format string, columnPtr interface{})|GROUP BY format|


# having
SetHavingはHavingでSetWhere系とほぼ同じようなメソッドと動作


# order by
|method|sql|
|---|---|
|SetOrderBy(columnPtr interface{})|ORDER BY columnPtr|
|SetOrderByAsc(columnPtr interface{})|ORDER BY columnPtr|
|SetOrderByAscString(str string)|ORDER BY str|
|SetOrderByAscFormat(format string, columnPtr interface{})|ORDER BY format|
|SetOrderByDesc(columnPtr interface{}|ORDER BY columnPtr DESC|
|SetOrderByDescString(str string)|ORDER BY str DESC|
|SetOrderByDescFormat(format string, columnPtr interface{})|ORDER BY format DESC|


# limit
|method|sql|
|---|---|
|SetLimit(num int)|LIMIT num|

# offset
|method|sql|
|---|---|
|SetOffset(num int)|OFFSET num|



