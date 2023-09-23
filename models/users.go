/*
 license x
*/

package models

type User struct {
	ID   int64  `sql:"type:serial,primary key" json:"id"`
	Name string `sql:"type:varchar"            json:"name"`
}

//nolint:gochecknoglobals // nil instance of User.
var NilUser User
