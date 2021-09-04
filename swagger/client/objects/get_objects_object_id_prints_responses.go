// Code generated by go-swagger; DO NOT EDIT.

package objects

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

// GetObjectsObjectIDPrintsReader is a Reader for the GetObjectsObjectIDPrints structure.
type GetObjectsObjectIDPrintsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetObjectsObjectIDPrintsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetObjectsObjectIDPrintsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetObjectsObjectIDPrintsOK creates a GetObjectsObjectIDPrintsOK with default headers values
func NewGetObjectsObjectIDPrintsOK() *GetObjectsObjectIDPrintsOK {
	return &GetObjectsObjectIDPrintsOK{}
}

/* GetObjectsObjectIDPrintsOK describes a response with status code 200, with default header values.

The 3D file information
*/
type GetObjectsObjectIDPrintsOK struct {
	Payload *GetObjectsObjectIDPrintsOKBody
}

func (o *GetObjectsObjectIDPrintsOK) Error() string {
	return fmt.Sprintf("[GET /objects/{object_id}/prints][%d] getObjectsObjectIdPrintsOK  %+v", 200, o.Payload)
}
func (o *GetObjectsObjectIDPrintsOK) GetPayload() *GetObjectsObjectIDPrintsOKBody {
	return o.Payload
}

func (o *GetObjectsObjectIDPrintsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetObjectsObjectIDPrintsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*GetObjectsObjectIDPrintsOKBody get objects object ID prints o k body
swagger:model GetObjectsObjectIDPrintsOKBody
*/
type GetObjectsObjectIDPrintsOKBody struct {

	// items
	Items []*models.Print `json:"items"`

	// total count
	TotalCount int64 `json:"total_count,omitempty"`
}

// Validate validates this get objects object ID prints o k body
func (o *GetObjectsObjectIDPrintsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetObjectsObjectIDPrintsOKBody) validateItems(formats strfmt.Registry) error {
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
					return ve.ValidateName("getObjectsObjectIdPrintsOK" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get objects object ID prints o k body based on the context it is used
func (o *GetObjectsObjectIDPrintsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetObjectsObjectIDPrintsOKBody) contextValidateItems(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Items); i++ {

		if o.Items[i] != nil {
			if err := o.Items[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getObjectsObjectIdPrintsOK" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetObjectsObjectIDPrintsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetObjectsObjectIDPrintsOKBody) UnmarshalBinary(b []byte) error {
	var res GetObjectsObjectIDPrintsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}