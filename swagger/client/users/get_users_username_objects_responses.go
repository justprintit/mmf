// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/justprintit/mmf/swagger/models"
)

// GetUsersUsernameObjectsReader is a Reader for the GetUsersUsernameObjects structure.
type GetUsersUsernameObjectsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetUsersUsernameObjectsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetUsersUsernameObjectsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetUsersUsernameObjectsOK creates a GetUsersUsernameObjectsOK with default headers values
func NewGetUsersUsernameObjectsOK() *GetUsersUsernameObjectsOK {
	return &GetUsersUsernameObjectsOK{}
}

/* GetUsersUsernameObjectsOK describes a response with status code 200, with default header values.

List of objects
*/
type GetUsersUsernameObjectsOK struct {
	Payload *models.Objects
}

func (o *GetUsersUsernameObjectsOK) Error() string {
	return fmt.Sprintf("[GET /users/{username}/objects][%d] getUsersUsernameObjectsOK  %+v", 200, o.Payload)
}
func (o *GetUsersUsernameObjectsOK) GetPayload() *models.Objects {
	return o.Payload
}

func (o *GetUsersUsernameObjectsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Objects)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}