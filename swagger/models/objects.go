// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Objects objects
//
// swagger:model Objects
type Objects struct {

	// items
	Items *ObjectsItemsTuple0 `json:"items,omitempty"`

	// total count
	TotalCount int64 `json:"total_count,omitempty"`
}

// Validate validates this objects
func (m *Objects) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Objects) validateItems(formats strfmt.Registry) error {
	if swag.IsZero(m.Items) { // not required
		return nil
	}

	if m.Items != nil {
		if err := m.Items.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("items")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this objects based on the context it is used
func (m *Objects) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Objects) contextValidateItems(ctx context.Context, formats strfmt.Registry) error {

	if m.Items != nil {
		if err := m.Items.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("items")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Objects) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Objects) UnmarshalBinary(b []byte) error {
	var res Objects
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// ObjectsItemsTuple0 ObjectsItemsTuple0 a representation of an anonymous Tuple type
//
// swagger:model ObjectsItemsTuple0
type ObjectsItemsTuple0 struct {

	// p0
	// Required: true
	P0 *Object `json:"-"` // custom serializer

}

// UnmarshalJSON unmarshals this tuple type from a JSON array
func (m *ObjectsItemsTuple0) UnmarshalJSON(raw []byte) error {
	// stage 1, get the array but just the array
	var stage1 []json.RawMessage
	buf := bytes.NewBuffer(raw)
	dec := json.NewDecoder(buf)
	dec.UseNumber()

	if err := dec.Decode(&stage1); err != nil {
		return err
	}

	// stage 2: hydrates struct members with tuple elements
	if len(stage1) > 0 {
		var dataP0 Object
		buf = bytes.NewBuffer(stage1[0])
		dec := json.NewDecoder(buf)
		dec.UseNumber()
		if err := dec.Decode(&dataP0); err != nil {
			return err
		}
		m.P0 = &dataP0

	}

	return nil
}

// MarshalJSON marshals this tuple type into a JSON array
func (m ObjectsItemsTuple0) MarshalJSON() ([]byte, error) {
	data := []interface{}{
		m.P0,
	}

	return json.Marshal(data)
}

// Validate validates this objects items tuple0
func (m *ObjectsItemsTuple0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateP0(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ObjectsItemsTuple0) validateP0(formats strfmt.Registry) error {

	if err := validate.Required("P0", "body", m.P0); err != nil {
		return err
	}

	if m.P0 != nil {
		if err := m.P0.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("P0")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this objects items tuple0 based on the context it is used
func (m *ObjectsItemsTuple0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateP0(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ObjectsItemsTuple0) contextValidateP0(ctx context.Context, formats strfmt.Registry) error {

	if m.P0 != nil {
		if err := m.P0.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("P0")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ObjectsItemsTuple0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ObjectsItemsTuple0) UnmarshalBinary(b []byte) error {
	var res ObjectsItemsTuple0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}