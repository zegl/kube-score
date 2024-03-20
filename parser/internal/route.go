package internal

import (
	routev1 "github.com/openshift/api/route/v1"
	ks "github.com/zegl/kube-score/domain"
)

var _ ks.Route = (*RouteV1)(nil)

type RouteV1 struct {
	Obj      routev1.Route
	Location ks.FileLocation
}

func (r RouteV1) FileLocation() ks.FileLocation {
	return r.Location
}

func (r RouteV1) Route() routev1.Route {
	return r.Obj
}
