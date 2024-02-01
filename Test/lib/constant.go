package lib

const (
	ImageMaxSize int = 2000
)

const (
	Uid         = "Uid"
	CtxAuthUser = "User"
)

const (
	HeaderXOperator = "x-operator"
	HeaderXToken    = "x-token"

	HeaderXAccessKey = "X-Access-Key"
	HeaderXSignature = "X-Signature"
)

const (
	GinCtxKeyClaims   = "claims"
	GinCtxKeyOperator = "operator"

	GinCtxKeyBody = "body"
)

const (
	TencentSMSSignMethod   = "TC3-HMAC-SHA256"
	TencentSMSSignName     = "海湃领客"
	TencentSMSCommonCodeId = "1847408"
	TencentSMSSuccessCode  = "Ok"
)

const (
	PhoneRegex    = `^1[3-9]\d{9}$`
	ImageUrlRegex = "^(http|https)://[a-zA-Z0-9]+(.[a-zA-Z0-9]+)+([a-zA-Z0-9-._?,'+/\\~:#[]@!$&*])*(.png|.jpg|.jpeg)$"
)

const (
	SdProgressPath      = "/sdapi/v1/progress"
	OssUserObjectPath   = "user"
	OssSystemObjectPath = "system"
)

type GqlCtxKey string

const (
	GqlCtxKeyGin      GqlCtxKey = "gin"
	GqlCtxKeyOperator GqlCtxKey = "operator"
	GqlCtxKeyUserId   GqlCtxKey = "user"
)

const (
	RouterGroupGql = "/gql"
)
