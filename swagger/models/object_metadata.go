// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ObjectMetadata object metadata
//
// swagger:model ObjectMetadata
type ObjectMetadata struct {

	// client url
	ClientURL string `json:"client_url,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// dimensions
	Dimensions string `json:"dimensions,omitempty"`

	// filament quantity
	FilamentQuantity string `json:"filament_quantity,omitempty"`

	// files
	Files []*FileUploadRequest `json:"files"`

	// how to
	HowTo string `json:"how_to,omitempty"`

	// images
	Images []*ImageUploadRequest `json:"images"`

	// licenses
	Licenses []*License `json:"licenses"`

	// name
	Name string `json:"name,omitempty"`

	// support free
	SupportFree bool `json:"support_free,omitempty"`

	// tags
	Tags string `json:"tags,omitempty"`

	// time to do from
	TimeToDoFrom int64 `json:"time_to_do_from,omitempty"`

	// time to do to
	TimeToDoTo int64 `json:"time_to_do_to,omitempty"`

	// 2: Public, 0: Private
	Visibility *int64 `json:"visibility,omitempty"`
}

// Validate validates this object metadata
func (m *ObjectMetadata) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFiles(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImages(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLicenses(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ObjectMetadata) validateFiles(formats strfmt.Registry) error {
	if swag.IsZero(m.Files) { // not required
		return nil
	}

	for i := 0; i < len(m.Files); i++ {
		if swag.IsZero(m.Files[i]) { // not required
			continue
		}

		if m.Files[i] != nil {
			if err := m.Files[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("files" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ObjectMetadata) validateImages(formats strfmt.Registry) error {
	if swag.IsZero(m.Images) { // not required
		return nil
	}

	for i := 0; i < len(m.Images); i++ {
		if swag.IsZero(m.Images[i]) { // not required
			continue
		}

		if m.Images[i] != nil {
			if err := m.Images[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("images" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ObjectMetadata) validateLicenses(formats strfmt.Registry) error {
	if swag.IsZero(m.Licenses) { // not required
		return nil
	}

	for i := 0; i < len(m.Licenses); i++ {
		if swag.IsZero(m.Licenses[i]) { // not required
			continue
		}

		if m.Licenses[i] != nil {
			if err := m.Licenses[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("licenses" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this object metadata based on the context it is used
func (m *ObjectMetadata) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFiles(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateImages(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateLicenses(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ObjectMetadata) contextValidateFiles(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Files); i++ {

		if m.Files[i] != nil {
			if err := m.Files[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("files" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ObjectMetadata) contextValidateImages(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Images); i++ {

		if m.Images[i] != nil {
			if err := m.Images[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("images" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ObjectMetadata) contextValidateLicenses(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Licenses); i++ {

		if m.Licenses[i] != nil {
			if err := m.Licenses[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("licenses" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ObjectMetadata) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ObjectMetadata) UnmarshalBinary(b []byte) error {
	var res ObjectMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
