// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/justprintit/mmf/swagger/models"
)

// GetUsersUsernameObjectsLikedReader is a Reader for the GetUsersUsernameObjectsLiked structure.
type GetUsersUsernameObjectsLikedReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetUsersUsernameObjectsLikedReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetUsersUsernameObjectsLikedOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetUsersUsernameObjectsLikedOK creates a GetUsersUsernameObjectsLikedOK with default headers values
func NewGetUsersUsernameObjectsLikedOK() *GetUsersUsernameObjectsLikedOK {
	return &GetUsersUsernameObjectsLikedOK{}
}

/* GetUsersUsernameObjectsLikedOK describes a response with status code 200, with default header values.

List of objects
*/
type GetUsersUsernameObjectsLikedOK struct {
	Payload *GetUsersUsernameObjectsLikedOKBody
}

func (o *GetUsersUsernameObjectsLikedOK) Error() string {
	return fmt.Sprintf("[GET /users/{username}/objects_liked][%d] getUsersUsernameObjectsLikedOK  %+v", 200, o.Payload)
}
func (o *GetUsersUsernameObjectsLikedOK) GetPayload() *GetUsersUsernameObjectsLikedOKBody {
	return o.Payload
}

func (o *GetUsersUsernameObjectsLikedOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetUsersUsernameObjectsLikedOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetUsersUsernameObjectsLikedOKBody get users username objects liked o k body
swagger:model GetUsersUsernameObjectsLikedOKBody
*/
type GetUsersUsernameObjectsLikedOKBody struct {

	// items
	Items []*models.Object `json:"items"`

	// total count
	TotalCount int64 `json:"total_count,omitempty"`
}

// Validate validates this get users username objects liked o k body
func (o *GetUsersUsernameObjectsLikedOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetUsersUsernameObjectsLikedOKBody) validateItems(formats strfmt.Registry) error {
	if swag.IsZero(o.Items) { // not required
		return nil
	}

	for i := 0; i < len(o.Items); i++ {
		if swag.IsZero(o.Items[i]) { // not required
			continue
		}

		if o.Items[i] != nil {
			if err := o.Items[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getUsersUsernameObjectsLikedOK" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get users username objects liked o k body based on the context it is used
func (o *GetUsersUsernameObjectsLikedOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetUsersUsernameObjectsLikedOKBody) contextValidateItems(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Items); i++ {

		if o.Items[i] != nil {
			if err := o.Items[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getUsersUsernameObjectsLikedOK" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetUsersUsernameObjectsLikedOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetUsersUsernameObjectsLikedOKBody) UnmarshalBinary(b []byte) error {
	var res GetUsersUsernameObjectsLikedOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}