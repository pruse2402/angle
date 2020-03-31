package errs

import (
	"fmt"
)

// UiError holds custom error for handling DB Layer
// Specific errors
type UIErr struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	ValError error  `json:"errors,omitempty"`
}

type CloneErr struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	ValError error  `json:"errors,omitempty"`
}

type ValidationError struct {
	Errors map[string]interface{}
}

func (ve *ValidationError) Error() string {
	s, n := "", 0
	for _, e := range ve.Errors {
		if e != "" {
			if n == 0 {
				s = e.(string)
			}
		}
		n++
	}
	switch n {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return fmt.Sprintf("%s (and %d other errors)", s, n-1)
}

// Error function for satisfying error interface
func (err *UIErr) Error() string {
	return fmt.Sprintf("Message: %s, Code: %d\n", err.Message, err.Code)
}

// Code returns a code
func (err *UIErr) ErrCode() int {
	return err.Code
}

func (err *CloneErr) Error() string {
	return fmt.Sprintf("Message: %s, Code: %d\n", err.Message, err.Code)
}

func New(errMsg string, code int) *UIErr {
	return &UIErr{Message: errMsg, Code: code}
}

//func WrapUIErr(err error, msg string) *UIErr {
//	log.Println("Error: ", msg)
//
//	if msg == "" {
//		msg = "Oops, Something went wrong. Please try again"
//	}
//
//	if mgo.IsDup(err) {
//		return &UIErr{409, msg}
//	}
//
//	if err.Error() == "EOF" {
//		return &UIErr{503, msg}
//	}
//
//	var uiErr UIErr
//	switch err.(type) {
//
//	case *mgo.QueryError:
//		uiErr.Code = 410
//	case *mgo.LastError:
//		uiErr.Code = 400
//	default:
//		uiErr.Code = 500
//	}
//
//	uiErr.Message = msg
//	return &uiErr
//}

// ApiErr holds error for API integration
type ApiErr struct {
	StatusCode int               `json:"statusCode"`
	ExternalId string            `json:"externalId"`
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors"`
}

func (err ApiErr) Error() string {
	return fmt.Sprintf("Message: %s, Code: %d\n", err.Message, err.StatusCode)
}

type AppError struct {
	Code    int                    `json:"-"`
	Message string                 `json:"message"`
	Err     string                 `json:"error,omitempty"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func (ae *AppError) Error() string {
	return ae.Message
}
