/*
 license x
*/

package model

// User storage db data model.
type User struct {
	ID   int64  `sql:"type:serial,primary key" json:"-"`
	Name string `sql:"type:varchar"            json:"-"`
}

//nolint:gochecknoglobals // nil instance of User.
var NilUser User

// UserView model to return the date.
type UserView struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

//nolint:gochecknoglobals // nil instance of UserView.
var NilUserView UserView
