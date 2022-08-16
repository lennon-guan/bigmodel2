package test_bigmodel_generating

//go:generate ../bin/bigmodel $GOFILE
//bigmodel-source: user *models.User
//bigmodel-source: header *views.HeaderInfo
type (
	Activity interface {
		ID() int              // user.Id
		Nickname() (b string) // user
		IP() string           // header.RemoteIP
	}
)
