// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// FileUploadStatus file upload status
//
// swagger:model FileUploadStatus
type FileUploadStatus struct {

	// filename
	Filename string `json:"filename,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// upload id
	UploadID string `json:"upload_id,omitempty"`
}

// Validate validates this file upload status
func (m *FileUploadStatus) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this file upload status based on context it is used
func (m *FileUploadStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *FileUploadStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *FileUploadStatus) UnmarshalBinary(b []byte) error {
	var res FileUploadStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
