// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// FileUploadRequest file upload request
//
// swagger:model FileUploadRequest
type FileUploadRequest struct {

	// filename
	Filename string `json:"filename,omitempty"`

	// size
	Size int64 `json:"size,omitempty"`
}

// Validate validates this file upload request
func (m *FileUploadRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this file upload request based on context it is used
func (m *FileUploadRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FileUploadRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FileUploadRequest) UnmarshalBinary(b []byte) error {
	var res FileUploadRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
