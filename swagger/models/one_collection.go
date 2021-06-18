// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// OneCollection one collection
//
// swagger:model OneCollection
type OneCollection struct {

	// cover object
	CoverObject *Object `json:"cover_object,omitempty"`

	// The value is specified in ISO 8601 (YYYY-MM-DDThh:mm:ss.sZ) format.
	CreatedAt string `json:"created_at,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// featured
	Featured bool `json:"featured,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// owner
	Owner *User `json:"owner,omitempty"`

	// public
	Public bool `json:"public,omitempty"`

	// slug
	Slug string `json:"slug,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this one collection
func (m *OneCollection) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCoverObject(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOwner(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCollection) validateCoverObject(formats strfmt.Registry) error {
	if swag.IsZero(m.CoverObject) { // not required
		return nil
	}

	if m.CoverObject != nil {
		if err := m.CoverObject.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cover_object")
			}
			return err
		}
	}

	return nil
}

func (m *OneCollection) validateOwner(formats strfmt.Registry) error {
	if swag.IsZero(m.Owner) { // not required
		return nil
	}

	if m.Owner != nil {
		if err := m.Owner.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("owner")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this one collection based on the context it is used
func (m *OneCollection) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCoverObject(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateOwner(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCollection) contextValidateCoverObject(ctx context.Context, formats strfmt.Registry) error {

	if m.CoverObject != nil {
		if err := m.CoverObject.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cover_object")
			}
			return err
		}
	}

	return nil
}

func (m *OneCollection) contextValidateOwner(ctx context.Context, formats strfmt.Registry) error {

	if m.Owner != nil {
		if err := m.Owner.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("owner")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OneCollection) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCollection) UnmarshalBinary(b []byte) error {
	var res OneCollection
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
