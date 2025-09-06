package db

import "fmt"

func Put(key string, value string) {
	fmt.Printf("DB PUT: key=%q value=%q\n", key, value)
	if value == "" {
		delete(db, key)
	} else {
		db[key] = value
	}
}

func Get(key string) string {
	fmt.Printf("DB GET: key=%q\n", key)
	return db[key]
}

var db = make(map[string]string)
