package mocks

import(
	"code.google.com/p/gomock/gomock"
  "github.com/onsi/ginkgo/thirdparty/gomocktestreporter"
)


// Support Interfaces
//----------------------------------------------------------------------

func NewCtrl() *gomock.Controller {
	return gomock.NewController(gomocktestreporter.New())
}
