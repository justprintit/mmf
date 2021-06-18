// Code generated by go-swagger; DO NOT EDIT.

package files

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PostImageReader is a Reader for the PostImage structure.
type PostImageReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PostImageReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewPostImageCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPostImageCreated creates a PostImageCreated with default headers values
func NewPostImageCreated() *PostImageCreated {
	return &PostImageCreated{}
}

/* PostImageCreated describes a response with status code 201, with default header values.

image created
*/
type PostImageCreated struct {
	Payload *PostImageCreatedBody
}

func (o *PostImageCreated) Error() string {
	return fmt.Sprintf("[POST /image][%d] postImageCreated  %+v", 201, o.Payload)
}
func (o *PostImageCreated) GetPayload() *PostImageCreatedBody {
	return o.Payload
}

func (o *PostImageCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(PostImageCreatedBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*PostImageCreatedBody post image created body
swagger:model PostImageCreatedBody
*/
type PostImageCreatedBody struct {

	// filename
	Filename string `json:"filename,omitempty"`
}

// Validate validates this post image created body
func (o *PostImageCreatedBody) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this post image created body based on context it is used
func (o *PostImageCreatedBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *PostImageCreatedBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostImageCreatedBody) UnmarshalBinary(b []byte) error {
	var res PostImageCreatedBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
