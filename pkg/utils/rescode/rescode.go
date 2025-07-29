package rescode

import "net/http"

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

var (
	// Success
	Success = Code{200, "success", http.StatusOK, ""}

	Invalid = Code{99002, "common.invalid", http.StatusBadRequest, ""}
)
