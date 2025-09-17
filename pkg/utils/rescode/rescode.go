package rescode

import (
	"errors"
	"net/http"
)

/*
001	NOT_FOUND	Dữ liệu không tồn tại
002	INVALID	Dữ liệu không hợp lệ (định dạng, logic...)
003	REQUIRED	Thiếu dữ liệu bắt buộc
004	ALREADY_EXISTS	Đã tồn tại dữ liệu trùng
005	PERMISSION_DENIED	Không có quyền thực hiện hành động
006	EXPIRED	Dữ liệu hoặc token đã hết hạn
007	MISMATCH	Dữ liệu không khớp (VD: password, code...)
008	CONFLICT	Xung đột trạng thái, logic nghiệp vụ
009	NOT_ALLOWED	Không được phép thực hiện
010	INTERNAL_ERROR	Lỗi nội bộ hệ thống
011	CREATION_FAILED	Tạo mới thất bại
012	UPDATE_FAILED	Cập nhật thất bại
013	DELETE_FAILED	Xoá thất bại
014	GET_FAILED	Truy vấn thất bại
015	BINDING_ERROR	Lỗi bind dữ liệu (JSON, form, param...)
016	VALIDATION_FAILED	Lỗi xác thực dữ liệu
017	AUTH_REQUIRED	Chưa đăng nhập hoặc chưa xác thực
018	TOO_MANY_REQUESTS	Spam, quá nhiều request
019	RATE_LIMITED	Bị giới hạn tần suất truy cập
020	UNSUPPORTED_OPERATION	Hành động không được hỗ trợ
021	DEPENDENCY_FAILED	Thất bại do lỗi từ hệ thống khác
022	DATABASE_ERROR	Lỗi truy vấn cơ sở dữ liệu
023	CACHE_ERROR	Lỗi Redis/memory cache
024	QUEUE_ERROR	Lỗi hàng đợi xử lý (RabbitMQ, Kafka...)
025	THIRD_PARTY_ERROR	Lỗi từ dịch vụ bên thứ ba
026	CONFIGURATION_ERROR	Lỗi cấu hình hệ thống
027	UNAUTHORIZED_ACTION	Đăng nhập rồi nhưng vẫn không đủ quyền
028	RESOURCE_LOCKED	Tài nguyên đang bị khóa (chờ xử lý khác)
029	ALREADY_PROCESSED	Đã xử lý trước đó (idempotent action)
030	UNEXPECTED_STATE	Trạng thái không hợp lệ
031	SESSION_EXPIRED	Phiên làm việc đã hết hạn
032	NOT_IMPLEMENTED	Chưa được triển khai
033	INVALID_CREDENTIALS	Sai tài khoản hoặc mật khẩu
034	VERIFICATION_FAILED	Mã xác thực sai hoặc hết hạn
035	UNSUPPORTED_MEDIA_TYPE	Kiểu dữ liệu không hỗ trợ (ví dụ: file type)
036	INVALID_FILE_FORMAT	File không đúng định dạng
037	MAX_UPLOAD_SIZE_EXCEEDED	File quá lớn
038	DUPLICATE_REQUEST	Gửi request giống nhau nhiều lần
039	SIGNATURE_INVALID	Sai chữ ký số, hoặc không xác thực được
040	PAYMENT_REQUIRED	Cần thanh toán để tiếp tục
041	BALANCE_NOT_ENOUGH	Không đủ tiền
042	ACCOUNT_LOCKED	Tài khoản bị khóa
043	ACCOUNT_DISABLED	Tài khoản bị vô hiệu hoá
044	ACCOUNT_NOT_VERIFIED	Tài khoản chưa xác minh
*/
type Code struct {
	Value     int
	Key       string
	Http      int
	DetailKey string
	Err       error
}

func (c Code) Code() int {
	return c.Value
}

func (c Code) MessageKey() string {
	return c.Key
}

func (c Code) HTTPStatus() int {
	return c.Http
}

func (c Code) WithDetailKey(detail string) Code {
	c.DetailKey = detail
	return c
}

func (c Code) Error() string {
	if c.Err != nil {
		return c.Err.Error()
	}
	return c.Key
}

var (
	// Success
	Success = Code{200, "success", http.StatusOK, "", nil}

	// User
	UserNotFound      = Code{10001, "user.not_found", http.StatusNotFound, "", errors.New("user not found")}
	UserInvalid       = Code{10002, "user.invalid", http.StatusBadRequest, "", errors.New("invalid user data")}
	UserRequire       = Code{10003, "user.require", http.StatusBadRequest, "", errors.New("user required")}
	UserAlreadyExists = Code{10004, "user.already_exists", http.StatusConflict, "", errors.New("user already exists")}
	UserNotVerified   = Code{10044, "user.not_verified", http.StatusUnauthorized, "", errors.New("user not verified")}
	UserDisable       = Code{10043, "user.locked", http.StatusForbidden, "", errors.New("user disabled")}
	UserCreateFailed  = Code{10011, "user.create_failed", http.StatusInternalServerError, "", errors.New("failed to create user")}
	UserUpdateFailed  = Code{10012, "user.update_failed", http.StatusInternalServerError, "", errors.New("failed to update user")}
	UserDeleteFailed  = Code{10013, "user.delete_failed", http.StatusInternalServerError, "", errors.New("failed to delete user")}
	UserGetFailed     = Code{10014, "user.get_failed", http.StatusInternalServerError, "", errors.New("failed to get user")}

	// Token

	// Auth
	InvalidCredentials = Code{20001, "auth.invalid_credentials", http.StatusUnauthorized, "", errors.New("invalid credentials")}
	TokenExpired       = Code{20002, "auth.token_expired", http.StatusUnauthorized, "", errors.New("token expired")}
	TokenInvalid       = Code{20003, "auth.token_invalid", http.StatusUnauthorized, "", errors.New("invalid token")}
	PermissionDenied   = Code{20004, "auth.permission_denied", http.StatusForbidden, "", errors.New("permission denied")}
	AuthRequired       = Code{20005, "auth.required", http.StatusUnauthorized, "", errors.New("authentication required")}

	// File
	FileTooLarge        = Code{33001, "file.too_large", http.StatusRequestEntityTooLarge, "", errors.New("file too large")}
	UnsupportedFileType = Code{33002, "file.unsupported_type", http.StatusUnsupportedMediaType, "", errors.New("unsupported file type")}
	InvalidFileFormat   = Code{33003, "file.invalid_format", http.StatusBadRequest, "", errors.New("invalid file format")}

	// Common
	NotFound              = Code{99001, "common.not_found", http.StatusNotFound, "", errors.New("not found")}
	Invalid               = Code{99002, "common.invalid", http.StatusBadRequest, "", errors.New("invalid data")}
	Required              = Code{99003, "common.required", http.StatusBadRequest, "", errors.New("required data missing")}
	AlreadyExists         = Code{99004, "common.already_exists", http.StatusConflict, "", errors.New("already exists")}
	Expired               = Code{99006, "common.expired", http.StatusUnauthorized, "", errors.New("expired")}
	Mismatch              = Code{99007, "common.mismatch", http.StatusBadRequest, "", errors.New("data mismatch")}
	Conflict              = Code{99008, "common.conflict", http.StatusConflict, "", errors.New("conflict")}
	NotAllowed            = Code{99009, "common.not_allowed", http.StatusForbidden, "", errors.New("not allowed")}
	Internal              = Code{99010, "common.internal", http.StatusInternalServerError, "", errors.New("internal error")}
	CreationFailed        = Code{99011, "common.creation_failed", http.StatusInternalServerError, "", errors.New("creation failed")}
	UpdateFailed          = Code{99012, "common.update_failed", http.StatusInternalServerError, "", errors.New("update failed")}
	DeleteFailed          = Code{99013, "common.delete_failed", http.StatusInternalServerError, "", errors.New("delete failed")}
	ValidationFailed      = Code{99016, "common.validation_failed", http.StatusBadRequest, "", errors.New("validation failed")}
	TooManyRequests       = Code{99018, "common.too_many_requests", http.StatusTooManyRequests, "", errors.New("too many requests")}
	Database              = Code{99022, "common.database", http.StatusInternalServerError, "", errors.New("database error")}
	Cache                 = Code{99023, "common.cache", http.StatusInternalServerError, "", errors.New("cache error")}
	Queue                 = Code{99024, "common.queue", http.StatusInternalServerError, "", errors.New("queue error")}
	ThirdParty            = Code{99025, "common.third_party", http.StatusBadGateway, "", errors.New("third party error")}
	Config                = Code{99026, "common.configuration", http.StatusInternalServerError, "", errors.New("configuration error")}
	UnAuthorized          = Code{99027, "common.un_authorized", http.StatusUnauthorized, "", errors.New("unauthorized")}
	ResourceLocked        = Code{99028, "common.resource_locked", http.StatusConflict, "", errors.New("resource locked")}
	AlreadyProcessed      = Code{99029, "common.already_processed", http.StatusConflict, "", errors.New("already processed")}
	UnexpectedState       = Code{99030, "common.unexpected_state", http.StatusInternalServerError, "", errors.New("unexpected state")}
	SessionExpired        = Code{99031, "common.session_expired", http.StatusUnauthorized, "", errors.New("session expired")}
	NotImplemented        = Code{99032, "common.not_implemented", http.StatusNotImplemented, "", errors.New("not implemented")}
	VerificationFailed    = Code{99034, "common.verification_failed", http.StatusBadRequest, "", errors.New("verification failed")}
	UnsupportedMediaType  = Code{99035, "common.unsupported_media_type", http.StatusUnsupportedMediaType, "", errors.New("unsupported media type")}
	MaxUploadSizeExceeded = Code{99037, "common.max_upload_size_exceeded", http.StatusRequestEntityTooLarge, "", errors.New("max upload size exceeded")}
	DuplicateRequest      = Code{99038, "common.duplicate_request", http.StatusConflict, "", errors.New("duplicate request")}
	SignatureInvalid      = Code{99039, "common.signature_invalid", http.StatusBadRequest, "", errors.New("signature invalid")}
	PaymentRequired       = Code{99040, "common.payment_required", http.StatusPaymentRequired, "", errors.New("payment required")}
	BalanceNotEnough      = Code{99041, "common.balance_not_enough", http.StatusForbidden, "", errors.New("balance not enough")}
	AccountLocked         = Code{99042, "common.account_locked", http.StatusForbidden, "", errors.New("account locked")}
	AccountDisabled       = Code{99043, "common.account_disabled", http.StatusForbidden, "", errors.New("account disabled")}
	AccountNotVerified    = Code{99044, "common.account_not_verified", http.StatusUnauthorized, "", errors.New("account not verified")}
)
