// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/justprintit/mmf/swagger/models"
)

// PostObjectReader is a Reader for the PostObject structure.
type PostObjectReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostObjectReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPostObjectOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostObjectOK creates a PostObjectOK with default headers values
func NewPostObjectOK() *PostObjectOK {
	return &PostObjectOK{}
}

/* PostObjectOK describes a response with status code 200, with default header values.

The 3D object object
*/
type PostObjectOK struct {
	Payload *models.ObjectUpload
}

func (o *PostObjectOK) Error() string {
	return fmt.Sprintf("[POST /object][%d] postObjectOK  %+v", 200, o.Payload)
}
func (o *PostObjectOK) GetPayload() *models.ObjectUpload {
	return o.Payload
}

func (o *PostObjectOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ObjectUpload)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}