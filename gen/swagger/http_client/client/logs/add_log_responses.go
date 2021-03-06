// Code generated by go-swagger; DO NOT EDIT.

package logs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/vk-cs/iot-go-agent-sdk/gen/swagger/http_client/models"
)

// AddLogReader is a Reader for the AddLog structure.
type AddLogReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *AddLogReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewAddLogOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewAddLogBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewAddLogUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := NewAddLogTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewAddLogInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewAddLogOK creates a AddLogOK with default headers values
func NewAddLogOK() *AddLogOK {
	return &AddLogOK{}
}

/* AddLogOK describes a response with status code 200, with default header values.

OK
*/
type AddLogOK struct {
}

func (o *AddLogOK) Error() string {
	return fmt.Sprintf("[POST /logs][%d] addLogOK ", 200)
}

func (o *AddLogOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewAddLogBadRequest creates a AddLogBadRequest with default headers values
func NewAddLogBadRequest() *AddLogBadRequest {
	return &AddLogBadRequest{}
}

/* AddLogBadRequest describes a response with status code 400, with default header values.

Bad params suplied
*/
type AddLogBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *AddLogBadRequest) Error() string {
	return fmt.Sprintf("[POST /logs][%d] addLogBadRequest  %+v", 400, o.Payload)
}
func (o *AddLogBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AddLogBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddLogUnauthorized creates a AddLogUnauthorized with default headers values
func NewAddLogUnauthorized() *AddLogUnauthorized {
	return &AddLogUnauthorized{}
}

/* AddLogUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type AddLogUnauthorized struct {
	Payload *models.ErrorResponse
}

func (o *AddLogUnauthorized) Error() string {
	return fmt.Sprintf("[POST /logs][%d] addLogUnauthorized  %+v", 401, o.Payload)
}
func (o *AddLogUnauthorized) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AddLogUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddLogTooManyRequests creates a AddLogTooManyRequests with default headers values
func NewAddLogTooManyRequests() *AddLogTooManyRequests {
	return &AddLogTooManyRequests{}
}

/* AddLogTooManyRequests describes a response with status code 429, with default header values.

Too Many Requests
*/
type AddLogTooManyRequests struct {
	Payload *models.ErrorResponse
}

func (o *AddLogTooManyRequests) Error() string {
	return fmt.Sprintf("[POST /logs][%d] addLogTooManyRequests  %+v", 429, o.Payload)
}
func (o *AddLogTooManyRequests) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AddLogTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAddLogInternalServerError creates a AddLogInternalServerError with default headers values
func NewAddLogInternalServerError() *AddLogInternalServerError {
	return &AddLogInternalServerError{}
}

/* AddLogInternalServerError describes a response with status code 500, with default header values.

Internal server error
*/
type AddLogInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *AddLogInternalServerError) Error() string {
	return fmt.Sprintf("[POST /logs][%d] addLogInternalServerError  %+v", 500, o.Payload)
}
func (o *AddLogInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *AddLogInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
