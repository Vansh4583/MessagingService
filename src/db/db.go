package db

func Put(key string, value string) {
	if value == "" {
		delete(db, key)
	} else {
		db[key] = value
	}
}

func Get(key string) string {
	return db[key]
}

var db = make(map[string]string)
