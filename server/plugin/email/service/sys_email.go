package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/email/utils"
)

type EmailService struct {
}

//@author: [maplepie](https://github.com/maplepie)
//@function: EmailTest
//@description: 发送邮件测试
//@return: err error

func (e *EmailService) EmailTest() (err error) {
	subject := "test"
	body := "test"
	err = utils.EmailTest(subject, body)
	return err
}
