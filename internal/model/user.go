/*
 license x
*/

package model

type User struct {
	ID   int64  `sql:"type:serial,primary key" json:"-"`
	Name string `sql:"type:varchar"            json:"-"`
}

//nolint:gochecknoglobals // nil instance of User.
var NilUser User
