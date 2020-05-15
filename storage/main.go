package storage

import "github.com/gin-gonic/gin"

// Item : a struct to store a key and value combo
type Item struct {
	key   string
	value string
}

// DBS : The database struct to store stuff from the server
type DBS struct {
	Users  []string
	Items  []Item
	Number int
	DB     map[string]string
}

// PrintAll : Method to output the number field
func (d *DBS) PrintAll() map[string]interface{} {
	return gin.H{
		"number": d.Number,
	}
}

// UserList : Method to ouput item list
func (d *DBS) UserList() map[string]interface{} {
	return gin.H{
		"users": d.Users,
	}
}
