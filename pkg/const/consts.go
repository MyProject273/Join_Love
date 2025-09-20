package consts

const (
	// GRPC Server
	GrpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	UserAgentHeader            = "user-agent"
	XForwardedForHeader        = "x-forarded-for"
	AcceptLanguage             = "accept-language"

	// Token
	TokenTypeAccessToken  = 1
	TokenTypeRefreshToken = 2

	// Error code sqlc
	ForeignKeyViolation = "23503"
	UniqueViolation     = "23505"

	// Authorization
	AuthorizationHeaderKey  = "authorization"
	AuthorizationType       = "bearer"
	AuthorizationPayloadKey = "authorization_payload"

	Alphabet = "abcdefghijklmnopqrstuvwxyz"
)
