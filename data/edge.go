package data

import "github.com/jmoiron/sqlx"

func (db *DB) AddEdgesToDB(tx *sqlx.Tx, edges []Edge) error {

	_, err := tx.NamedExec("INSERT INTO edgecosts (id) VALUES (:id)", edges)

	return err

}
