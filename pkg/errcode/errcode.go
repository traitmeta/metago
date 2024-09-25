package errcode

import (
	"net/http"

	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

const (
	CodeOK     = 200
	MsgOK      = "Successful"
	CodeCustom = 16888
	MsgCustom  = "Custom error"
)

type Err struct {
	code     uint32
	httpCode int
	msg      string
}

func (e *Err) Code() uint32 {
	return e.code
}

func (e *Err) HTTPCode() int {
	return e.httpCode
}

func (e *Err) Error() string {
	return e.msg
}

var (
	NoErr = NewErr(CodeOK, MsgOK)

	ErrCustom                     = NewErr(CodeCustom, MsgCustom)
	ErrUnexpected                 = NewErr(16999, "Network error, please try again later", http.StatusInternalServerError)
	ErrTokenNotValidYet           = NewErr(19999, "Token illegal", http.StatusUnauthorized)
	ErrInvalidUrl                 = NewErr(20000, "URL is illegal")
	ErrInvalidHeader              = NewErr(20001, "Invalid request header")
	ErrInvalidParams              = NewErr(20002, "Parameter is illegal")
	ErrTokenVerify                = NewErr(20003, "Token check error", http.StatusUnauthorized)
	ErrTokenExpire                = NewErr(20004, "Expired token", http.StatusUnauthorized)
	ErrUserLogin                  = NewErr(20005, "User log out", http.StatusUnauthorized)
	ErrUserPrivilegeChange        = NewErr(20006, "Permission changed", http.StatusUnauthorized)
	ErrLockNotAcquire             = NewErr(20007, "Lock not released")
	ErrLockAcquire                = NewErr(20008, "Lock acquisition error")
	ErrLockNotRelease             = NewErr(20009, "Lock is not released")
	ErrLockRelease                = NewErr(20010, "Lock released err")
	ErrTWitterAddress             = NewErr(20011, "Twitter address illeage")
	ErrDiscordAddress             = NewErr(20012, "Discord address illeage")
	ErrAddress                    = NewErr(20013, "Address illeage")
	ErrProjectNotExist            = NewErr(21100, "Project not exist")
	ErrPresaleNotExist            = NewErr(21101, "Presale not exist")
	ErrCorporationNotExist        = NewErr(21102, "Corporation not exist")
	ErrProjectSaleUriNotExist     = NewErr(21103, "Calender uri not exist")
	ErrProjectFeatureNotExist     = NewErr(21105, "Project feature not exist")
	ErrProjectTeamNotExist        = NewErr(21115, "Project team not exist")
	ErrProjectWhitelistNotExist   = NewErr(21116, "Project whitelist not exist")
	ErrProjectAuditRecordNotExist = NewErr(21117, "Project audit record not exist")
	ErrFileNotExist               = NewErr(21201, "File does not exist")
	ErrFileChunkNotExist          = NewErr(21202, "File Chunk does not exist")
	ErrFileUnknown                = NewErr(21203, "File type unknown")
	ErrFileLoadErr                = NewErr(21204, "File load failed")

	ErrUserNotExist = NewErr(11300, "User does not exist")
)

var codeToErr = map[uint32]*Err{
	200: NoErr,

	16888: ErrCustom,
	16999: ErrUnexpected,

	19999: ErrTokenNotValidYet,
	20000: ErrInvalidUrl,
	20001: ErrInvalidHeader,
	20002: ErrInvalidParams,
	20003: ErrTokenVerify,
	20004: ErrTokenExpire,
	20005: ErrUserLogin,
	20006: ErrUserPrivilegeChange,
	20007: ErrLockNotAcquire,
	20008: ErrLockAcquire,
	20009: ErrLockNotRelease,
	20010: ErrLockRelease,
	20011: ErrTWitterAddress,
	20012: ErrDiscordAddress,
	20013: ErrAddress,
	21100: ErrProjectNotExist,
	21101: ErrPresaleNotExist,
	21102: ErrCorporationNotExist,
	21103: ErrProjectSaleUriNotExist,
	21105: ErrProjectFeatureNotExist,
	21115: ErrProjectTeamNotExist,
	21117: ErrProjectAuditRecordNotExist,
	21201: ErrFileNotExist,
	21202: ErrFileChunkNotExist,
	21203: ErrFileUnknown,
	21204: ErrFileLoadErr,
	21300: ErrUserNotExist,
}

func NewErr(code uint32, msg string, httpCode ...int) *Err {
	hc := http.StatusOK
	if len(httpCode) != 0 {
		hc = httpCode[0]
	}

	return &Err{code: code, httpCode: hc, msg: msg}
}

func GetCodeToErr() map[uint32]*Err {
	return codeToErr
}

func SetCodeToErr(code uint32, err *Err) error {
	if _, ok := codeToErr[code]; ok {
		return errors.New("has exist")
	}

	codeToErr[code] = err
	return nil
}

func NewCustomErr(msg string, httpCode ...int) *Err {
	return NewErr(CodeCustom, msg, httpCode...)
}

func IsErr(err error) bool {
	if err == nil {
		return true
	}

	_, ok := err.(*Err)
	return ok
}

func ParseErr(err error) *Err {
	if err == nil {
		return NoErr
	}

	if e, ok := err.(*Err); ok {
		return e
	}

	s, _ := status.FromError(err)
	c := uint32(s.Code())
	if c == CodeCustom {
		return NewCustomErr(s.Message())
	}

	return ParseCode(c)
}

func ParseCode(code uint32) *Err {
	if e, ok := codeToErr[code]; ok {
		return e
	}

	return ErrUnexpected
}
