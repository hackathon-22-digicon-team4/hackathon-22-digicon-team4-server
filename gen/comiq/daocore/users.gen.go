// Code generated by "script/dtogen". DO NOT EDIT.
package daocore

import (
  "context"
  "database/sql"
  "strings"
  "time"

  "github.com/Masterminds/squirrel"
)

const UserTableName = "users"

var UserAllColumns = []string{
  "id",
  
  "password",
  
  "created_at",
  
  "updated_at",
  
}

var UserColumnsWOMagics = []string{
  "id",
  
  "password",
  
  
  
}

var UserPrimaryKeyColumns = []string{
  "id",
  
  
  
  
}



type User struct {
  ID string
  Password string
  CreatedAt *time.Time
  UpdatedAt *time.Time
}

func (t *User) Values() []interface{} {
  return []interface{}{t.ID,t.Password,
  }
}

func (t *User) SetMap() map[string]interface{} {
  return map[string]interface{}{"id": t.ID,"password": t.Password,
  }
}


func (t *User) Ptrs() []interface{} {
  return []interface{}{
    &t.ID,
    &t.Password,
    &t.CreatedAt,
    &t.UpdatedAt,
  }
}




func IterateUser(sc interface{ Scan(...interface{}) error}) (User, error) {
  t := User{}
  if err := sc.Scan(t.Ptrs()...); err != nil {
    return User{}, MapError(err)
  }
  return t, nil
}

func SelectOneUserByID(ctx context.Context, txn *sql.Tx, id string) (User, error) {
  query, params, err := squirrel.
    Select(UserAllColumns...).
    From(UserTableName).
    Where(squirrel.Eq{
      "id": id,
    }).
    ToSql()
  if err != nil {
    return User{}, MapError(err)
  }
  stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
    return User{}, MapError(err)
	}
  return IterateUser(stmt.QueryRowContext(ctx, params...))
  
}




func InsertUser(ctx context.Context, txn *sql.Tx, records []*User) error {
  for i := range records {
		if records[i] == nil {
			records = append(records[:i], records[i+1:]...)
		}
	}
	if len(records) == 0 {
    return nil
  }
  sq := squirrel.Insert(UserTableName).Columns(UserColumnsWOMagics...)
	for _, r := range records {
		if r == nil {
			continue
		}
		sq = sq.Values(r.Values()...)
	}
	query, params, err := sq.ToSql()
	if err != nil {
		return err
	}
	stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
		return MapError(err)
	}
	if _, err = stmt.Exec(params...); err != nil {
		return MapError(err)
	}
	return nil
}

func UpdateUser(ctx context.Context, txn *sql.Tx, record User) error {
	sql, params, err := squirrel.Update(UserTableName).SetMap(record.SetMap()).
		Where(squirrel.Eq{
      "id": record.ID,
    }).
		ToSql()
	if err != nil {
		return err
	}
	stmt, err := txn.PrepareContext(ctx, sql)
	if err != nil {
		return MapError(err)
	}
	if _, err = stmt.Exec(params...); err != nil {
		return MapError(err)
	}
	return nil
}

func UpsertUser(ctx context.Context, txn *sql.Tx, record User) error {
	updateSQL, params, err := squirrel.Update(UserTableName).SetMap(record.SetMap()).ToSql()
	if err != nil {
		return err
	}
	updateSQL = strings.TrimPrefix(updateSQL, "UPDATE "+UserTableName+" SET ")
	query, params, err := squirrel.Insert(UserTableName).Columns(UserColumnsWOMagics...).Values(record.Values()...).SuffixExpr(squirrel.Expr("ON DUPLICATE KEY UPDATE "+updateSQL, params...)).ToSql()
	if err != nil {
		return err
	}
	stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
		return MapError(err)
	}
	if _, err = stmt.Exec(params...); err != nil {
		return MapError(err)
	}
	return nil
}

func DeleteOneUserByID(ctx context.Context, txn *sql.Tx, id string) error {
  query, params, err := squirrel.
    Delete(UserTableName).
    Where(squirrel.Eq{
      "id": id,
    }).
    ToSql()
  if err != nil {
    return MapError(err)
  }
  stmt, err := txn.PrepareContext(ctx, query)
	if err != nil {
		return MapError(err)
	}
	if _, err = stmt.Exec(params...); err != nil {
		return MapError(err)
	}
	return nil
}

