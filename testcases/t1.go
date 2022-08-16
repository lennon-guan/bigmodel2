package test_bigmodel_generating

//go:generate ../bin/bigmodel $GOFILE
//bigmodel-source: user *models.User
//bigmodel-source: header *views.HeaderInfo
type RiskControlModel interface {
	ID() int              // user.Id
	Nickname() (b string) // user
	IP() string           // header.RemoteIP
}

// 下面的方法是生成实现类后，自行实现的数据源获取方法
func (m *_RiskControlModelInnerImpl) getUser() *model.User       { return nil }
func (m *_RiskControlModelInnerImpl) getHeader() *request.Header { return nil }
