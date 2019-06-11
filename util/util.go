package util

import "github.com/syndtr/goleveldb/leveldb"

func ClearDB(DB *leveldb.DB) bool {
	DBIter := DB.NewIterator(nil, nil)
	for DBIter.Next() {
		key := DBIter.Key()
		e := DB.Delete(key, nil)
		ERR("DB Delete Error", e)
	}

	return false
}
