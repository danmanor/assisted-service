// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// V2UpdateClusterInstallConfigHandlerFunc turns a function with the right signature into a v2 update cluster install config handler
type V2UpdateClusterInstallConfigHandlerFunc func(V2UpdateClusterInstallConfigParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn V2UpdateClusterInstallConfigHandlerFunc) Handle(params V2UpdateClusterInstallConfigParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// V2UpdateClusterInstallConfigHandler interface for that can handle valid v2 update cluster install config params
type V2UpdateClusterInstallConfigHandler interface {
	Handle(V2UpdateClusterInstallConfigParams, interface{}) middleware.Responder
}

// NewV2UpdateClusterInstallConfig creates a new http.Handler for the v2 update cluster install config operation
func NewV2UpdateClusterInstallConfig(ctx *middleware.Context, handler V2UpdateClusterInstallConfigHandler) *V2UpdateClusterInstallConfig {
	return &V2UpdateClusterInstallConfig{Context: ctx, Handler: handler}
}

/*V2UpdateClusterInstallConfig swagger:route PATCH /v2/clusters/{cluster_id}/install-config installer v2UpdateClusterInstallConfig

Override values in the install config.

*/
type V2UpdateClusterInstallConfig struct {
	Context *middleware.Context
	Handler V2UpdateClusterInstallConfigHandler
}

func (o *V2UpdateClusterInstallConfig) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewV2UpdateClusterInstallConfigParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
