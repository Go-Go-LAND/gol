package gol

import (
	"fmt"
	"testing"
	"time"
)

type TestItem struct {
	Id        int       `column:"id" json:"id"`
	CreatedAt time.Time `column:"created_at" json:"createdAt"`
	CreatedBy NullInt64 `column:"created_by" json:"createdBy"`
	UpdatedAt time.Time `column:"updated_at" json:"updatedAt"`
	UpdatedBy NullInt64 `column:"updated_by" json:"updatedBy"`
	DeletedAt NullTime  `column:"deleted_at" json:"deletedAt"`
	DeletedBy NullInt64 `column:"deleted_by" json:"deletedBy"`
	Name      int       `column:"name" json:"name"`
	UserId    int       `column:"user_id" json:"userId"`
}

type TestUser struct {
	Id        int       `column:"id" json:"id"`
	CreatedAt time.Time `column:"created_at" json:"createdAt"`
	CreatedBy NullInt64 `column:"created_by" json:"createdBy"`
	UpdatedAt time.Time `column:"updated_at" json:"updatedAt"`
	UpdatedBy NullInt64 `column:"updated_by" json:"updatedBy"`
	DeletedAt NullTime  `column:"deleted_at" json:"deletedAt"`
	DeletedBy NullInt64 `column:"deleted_by" json:"deletedBy"`
	Name      string    `column:"name" json:"name"`
	Pass      string    `column:"pass" json:"pass"`
}

func TestQueryType_GetSelectQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		query.SetSelectAll(&testItemTable)
		str, valueList, err := query.GetSelectQuery()
		if err != nil {
			t.Error(err)
			return
		}

		{
			target := str

			check := `SELECT "test_item".* FROM "test_item"`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}

		{
			target := fmt.Sprintf("%v", valueList)

			var checkList []interface{}
			checkList = append(checkList, "")
			check := fmt.Sprintf("%v", checkList)

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error select not exist", func(t *testing.T) {
		query := QueryType{}
		_, _, err := query.GetSelectQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `select not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error table not exist", func(t *testing.T) {
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetSelectAll(&testItemTable.Name)
		_, _, err := query.GetSelectQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `select column meta not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})
}

func TestQueryType_GetSelectCountQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		query.SetSelectAll(&testItemTable)
		str, valueList, err := query.GetSelectQuery()
		if err != nil {
			t.Error(err)
			return
		}

		{
			target := str

			check := `SELECT "test_item".* FROM "test_item"`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}

		{
			target := fmt.Sprintf("%v", valueList)

			var checkList []interface{}
			checkList = append(checkList, "")
			check := fmt.Sprintf("%v", checkList)

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error table not exist", func(t *testing.T) {
		query := QueryType{}
		_, _, err := query.GetSelectCountQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `table not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})
}

func TestQueryType_GetInsertQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		item := TestItem{}
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		query.SetValuesColumn(
			&testItemTable.CreatedAt,
			&testItemTable.CreatedBy,
			&testItemTable.UpdatedAt,
			&testItemTable.UpdatedBy,
			&testItemTable.DeletedAt,
			&testItemTable.DeletedBy,
			&testItemTable.Name,
			&testItemTable.UserId,
		)
		query.SetValues(
			&item.CreatedAt,
			&item.CreatedBy,
			&item.UpdatedAt,
			&item.UpdatedBy,
			&item.DeletedAt,
			&item.DeletedBy,
			&item.Name,
			&item.UserId,
		)

		str, valueList, err := query.GetInsertQuery()
		if err != nil {
			t.Error(err)
			return
		}

		{
			target := str

			check := `INSERT INTO "test_item" ("created_at", "created_by", "updated_at", "updated_by", "deleted_at", "deleted_by", "name", "user_id") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}

		{
			target := fmt.Sprintf("%v", valueList)

			var checkList []interface{}
			checkList = append(checkList, item.CreatedAt)
			checkList = append(checkList, item.CreatedBy)
			checkList = append(checkList, item.UpdatedAt)
			checkList = append(checkList, item.UpdatedBy)
			checkList = append(checkList, item.DeletedAt)
			checkList = append(checkList, item.DeletedBy)
			checkList = append(checkList, item.Name)
			checkList = append(checkList, item.UserId)
			check := fmt.Sprintf("%v", checkList)

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error table not exist", func(t *testing.T) {
		query := QueryType{}
		_, _, err := query.GetInsertQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `table not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})
}

func TestQueryType_GetUpdateQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		item := TestItem{}
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		query.SetSet(&testItemTable.CreatedAt, item.CreatedAt)
		query.SetSet(&testItemTable.CreatedBy, item.CreatedBy)
		query.SetSet(&testItemTable.UpdatedAt, item.UpdatedAt)
		query.SetSet(&testItemTable.UpdatedBy, item.UpdatedBy)
		query.SetSet(&testItemTable.DeletedAt, item.DeletedAt)
		query.SetSet(&testItemTable.DeletedBy, item.DeletedBy)
		query.SetSet(&testItemTable.Name, item.Name)
		query.SetSet(&testItemTable.UserId, item.UserId)
		query.SetWhereString("1 = 1")
		str, valueList, err := query.GetUpdateQuery()
		if err != nil {
			t.Error(err)
			return
		}

		{
			target := str

			check := `UPDATE "test_item" SET "created_at" = $1, "created_by" = $2, "updated_at" = $3, "updated_by" = $4, "deleted_at" = $5, "deleted_by" = $6, "name" = $7, "user_id" = $8 WHERE 1 = 1`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}

		{
			target := fmt.Sprintf("%v", valueList)

			var checkList []interface{}
			checkList = append(checkList, item.CreatedAt)
			checkList = append(checkList, item.CreatedBy)
			checkList = append(checkList, item.UpdatedAt)
			checkList = append(checkList, item.UpdatedBy)
			checkList = append(checkList, item.DeletedAt)
			checkList = append(checkList, item.DeletedBy)
			checkList = append(checkList, item.Name)
			checkList = append(checkList, item.UserId)
			check := fmt.Sprintf("%v", checkList)

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error table not exist", func(t *testing.T) {
		query := QueryType{}
		_, _, err := query.GetUpdateQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `table not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error set not exist", func(t *testing.T) {
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		_, _, err := query.GetUpdateQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `set not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error where not exist", func(t *testing.T) {
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		query.SetSet(&testItemTable.Name, "name")
		_, _, err := query.GetUpdateQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `where not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})
}

func TestQueryType_GetDeleteQuery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//item := TestItem{}
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		query.SetWhereString("1 = 1")
		str, valueList, err := query.GetDeleteQuery()
		if err != nil {
			t.Error(err)
			return
		}

		{
			target := str

			check := `DELETE FROM "test_item" WHERE 1 = 1`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}

		{
			target := fmt.Sprintf("%v", valueList)

			var checkList []interface{}
			//checkList = append(checkList, item.CreatedAt)
			check := fmt.Sprintf("%v", checkList)

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error table not exist", func(t *testing.T) {
		query := QueryType{}
		_, _, err := query.GetDeleteQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `table not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})

	t.Run("error where not exist", func(t *testing.T) {
		testItemTable := TestItem{}

		query := QueryType{}
		query.SetTable(&testItemTable)
		_, _, err := query.GetDeleteQuery()
		{
			target := fmt.Sprintf("%v", err)

			check := `where not exist`

			if target != check {
				t.Error("target:", target)
				t.Error("check :", check)
				return
			}
		}
	})
}
