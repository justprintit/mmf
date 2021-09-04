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

// OneCategory one category
//
// swagger:model OneCategory
type OneCategory struct {

	// children
	Children *OneCategoryChildren `json:"children,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// parent
	Parent *OneCategoryParent `json:"parent,omitempty"`

	// slug
	Slug string `json:"slug,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this one category
func (m *OneCategory) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateChildren(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCategory) validateChildren(formats strfmt.Registry) error {
	if swag.IsZero(m.Children) { // not required
		return nil
	}

	if m.Children != nil {
		if err := m.Children.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("children")
			}
			return err
		}
	}

	return nil
}

func (m *OneCategory) validateParent(formats strfmt.Registry) error {
	if swag.IsZero(m.Parent) { // not required
		return nil
	}

	if m.Parent != nil {
		if err := m.Parent.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("parent")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this one category based on the context it is used
func (m *OneCategory) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateChildren(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateParent(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCategory) contextValidateChildren(ctx context.Context, formats strfmt.Registry) error {

	if m.Children != nil {
		if err := m.Children.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("children")
			}
			return err
		}
	}

	return nil
}

func (m *OneCategory) contextValidateParent(ctx context.Context, formats strfmt.Registry) error {

	if m.Parent != nil {
		if err := m.Parent.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("parent")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OneCategory) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCategory) UnmarshalBinary(b []byte) error {
	var res OneCategory
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// OneCategoryChildren one category children
//
// swagger:model OneCategoryChildren
type OneCategoryChildren struct {

	// items
	Items []*OneCategoryChildrenItemsItems0 `json:"items"`

	// total count
	TotalCount int64 `json:"total_count,omitempty"`
}

// Validate validates this one category children
func (m *OneCategoryChildren) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateItems(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCategoryChildren) validateItems(formats strfmt.Registry) error {
	if swag.IsZero(m.Items) { // not required
		return nil
	}

	for i := 0; i < len(m.Items); i++ {
		if swag.IsZero(m.Items[i]) { // not required
			continue
		}

		if m.Items[i] != nil {
			if err := m.Items[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("children" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this one category children based on the context it is used
func (m *OneCategoryChildren) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateItems(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OneCategoryChildren) contextValidateItems(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Items); i++ {

		if m.Items[i] != nil {
			if err := m.Items[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("children" + "." + "items" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *OneCategoryChildren) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCategoryChildren) UnmarshalBinary(b []byte) error {
	var res OneCategoryChildren
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// OneCategoryChildrenItemsItems0 one category children items items0
//
// swagger:model OneCategoryChildrenItemsItems0
type OneCategoryChildrenItemsItems0 struct {

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// slug
	Slug string `json:"slug,omitempty"`
}

// Validate validates this one category children items items0
func (m *OneCategoryChildrenItemsItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this one category children items items0 based on context it is used
func (m *OneCategoryChildrenItemsItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *OneCategoryChildrenItemsItems0) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCategoryChildrenItemsItems0) UnmarshalBinary(b []byte) error {
	var res OneCategoryChildrenItemsItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

// OneCategoryParent one category parent
//
// swagger:model OneCategoryParent
type OneCategoryParent struct {

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// slug
	Slug string `json:"slug,omitempty"`
}

// Validate validates this one category parent
func (m *OneCategoryParent) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this one category parent based on context it is used
func (m *OneCategoryParent) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *OneCategoryParent) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OneCategoryParent) UnmarshalBinary(b []byte) error {
	var res OneCategoryParent
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}