// Code generated by go-swagger; DO NOT EDIT.

package objects

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetObjectsObjectIDPrintsParams creates a new GetObjectsObjectIDPrintsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetObjectsObjectIDPrintsParams() *GetObjectsObjectIDPrintsParams {
	return &GetObjectsObjectIDPrintsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetObjectsObjectIDPrintsParamsWithTimeout creates a new GetObjectsObjectIDPrintsParams object
// with the ability to set a timeout on a request.
func NewGetObjectsObjectIDPrintsParamsWithTimeout(timeout time.Duration) *GetObjectsObjectIDPrintsParams {
	return &GetObjectsObjectIDPrintsParams{
		timeout: timeout,
	}
}

// NewGetObjectsObjectIDPrintsParamsWithContext creates a new GetObjectsObjectIDPrintsParams object
// with the ability to set a context for a request.
func NewGetObjectsObjectIDPrintsParamsWithContext(ctx context.Context) *GetObjectsObjectIDPrintsParams {
	return &GetObjectsObjectIDPrintsParams{
		Context: ctx,
	}
}

// NewGetObjectsObjectIDPrintsParamsWithHTTPClient creates a new GetObjectsObjectIDPrintsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetObjectsObjectIDPrintsParamsWithHTTPClient(client *http.Client) *GetObjectsObjectIDPrintsParams {
	return &GetObjectsObjectIDPrintsParams{
		HTTPClient: client,
	}
}

/* GetObjectsObjectIDPrintsParams contains all the parameters to send to the API endpoint
   for the get objects object ID prints operation.

   Typically these are written to a http.Request.
*/
type GetObjectsObjectIDPrintsParams struct {

	/* ObjectID.

	   The object identifier number
	*/
	ObjectID float64

	/* Page.

	   Page number. Default is 1
	*/
	Page *string

	/* PerPage.

	   Number of results per page. Default is 20
	*/
	PerPage *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get objects object ID prints params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetObjectsObjectIDPrintsParams) WithDefaults() *GetObjectsObjectIDPrintsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get objects object ID prints params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetObjectsObjectIDPrintsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) WithTimeout(timeout time.Duration) *GetObjectsObjectIDPrintsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) WithContext(ctx context.Context) *GetObjectsObjectIDPrintsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) WithHTTPClient(client *http.Client) *GetObjectsObjectIDPrintsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithObjectID adds the objectID to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) WithObjectID(objectID float64) *GetObjectsObjectIDPrintsParams {
	o.SetObjectID(objectID)
	return o
}

// SetObjectID adds the objectId to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) SetObjectID(objectID float64) {
	o.ObjectID = objectID
}

// WithPage adds the page to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) WithPage(page *string) *GetObjectsObjectIDPrintsParams {
	o.SetPage(page)
	return o
}

// SetPage adds the page to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) SetPage(page *string) {
	o.Page = page
}

// WithPerPage adds the perPage to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) WithPerPage(perPage *string) *GetObjectsObjectIDPrintsParams {
	o.SetPerPage(perPage)
	return o
}

// SetPerPage adds the perPage to the get objects object ID prints params
func (o *GetObjectsObjectIDPrintsParams) SetPerPage(perPage *string) {
	o.PerPage = perPage
}

// WriteToRequest writes these params to a swagger request
func (o *GetObjectsObjectIDPrintsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param object_id
	if err := r.SetPathParam("object_id", swag.FormatFloat64(o.ObjectID)); err != nil {
		return err
	}

	if o.Page != nil {

		// query param page
		var qrPage string

		if o.Page != nil {
			qrPage = *o.Page
		}
		qPage := qrPage
		if qPage != "" {

			if err := r.SetQueryParam("page", qPage); err != nil {
				return err
			}
		}
	}

	if o.PerPage != nil {

		// query param per_page
		var qrPerPage string

		if o.PerPage != nil {
			qrPerPage = *o.PerPage
		}
		qPerPage := qrPerPage
		if qPerPage != "" {

			if err := r.SetQueryParam("per_page", qPerPage); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
