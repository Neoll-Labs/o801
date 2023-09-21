/*
 license x
*/

package models

type User struct {
	ID   int64  `sql:"type:serial,primary key"`
	Name string `sql:"type:varchar"`
}

var NilUser User
