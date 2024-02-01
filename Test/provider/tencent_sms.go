package provider

import (
	"Technology-Blog/Test/config"
	"Technology-Blog/Test/lib"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TencentSms struct {
	Client *sms.Client
}

func NewTencentSms() *TencentSms {
	credential := common.NewCredential(
		config.Get(config.TencentSMSSecretId),
		config.Get(config.TencentSMSSecretKey),
	)
	newClient, err := sms.NewClient(
		credential,
		config.Get(config.TencentSmsRegion),
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}

	return &TencentSms{
		Client: newClient,
	}
}

func (*TencentSms) Start() {}

func (t *TencentSms) Close() {
	t.Client = nil
}

func (t *TencentSms) sendMessage(phone string, variables []*string) error {
	req := sms.NewSendSmsRequest()
	req.PhoneNumberSet = []*string{&phone}
	req.SmsSdkAppId = lib.StrToPtr(config.Get(config.TencentSmsAppId))
	req.TemplateId = lib.StrToPtr(lib.TencentSMSCommonCodeId)
	req.SignName = lib.StrToPtr(lib.TencentSMSSignName)
	req.TemplateParamSet = variables

	result, err := t.Client.SendSms(req)
	if err != nil {
		return err
	}
	if len(result.Response.SendStatusSet) == 0 {
		return fmt.Errorf("failed to send sms code to phone: %s", phone)
	}
	status := result.Response.SendStatusSet[0]
	if status.Code == nil {
		return fmt.Errorf("failed to send sms code to phone: %s", phone)
	}
	if *status.Code != lib.TencentSMSSuccessCode {
		return fmt.Errorf("failed to send sms code to phone: %s with error code: %s", phone, *status.Code)
	}
	return nil
}

func (*TencentSms) phoneSmsRedisKey(phone string) string {
	return fmt.Sprintf("sms_code_%s", phone)
}

func (t *TencentSms) SendLoginSmsCodeByTencent(phone string, codeLength int) error {
	code := lib.RandNumbers(codeLength)
	fmt.Println("code to sent:", code)
	err := t.sendMessage(phone, []*string{&code})
	if err != nil {
		return err
	}
	key := t.phoneSmsRedisKey(phone)
	_, err = Cache.Pool.Get().Do("SET", key, code, "EX", 600) // 600 s
	if err != nil {
		return err
	}
	return nil
}

func (t *TencentSms) ValidateSmsCode(phone string, code string) (bool, error) {
	key := t.phoneSmsRedisKey(phone)

	reply, err := Cache.Pool.Get().Do("GET", key)
	if err != nil {
		return false, err
	}
	if reply == nil {
		return false, nil
	}

	valid := code == string(reply.([]byte))

	if valid {
		go Cache.Pool.Get().Do("DEL", key)
	}
	return valid, nil
}
