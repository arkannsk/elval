package geo_polymorphism

import (
	"github.com/paulmach/orb/geojson"
)

//go:generate elval-gen gen -input . -openapi

// CreateLocationRequest запрос на создание локации
type CreateLocationRequest struct {
	// @oa:rewrite.ref "github.com/arkannsk/elval/test/integration/geo_polymorphism/pkg/docs.FeatureDocs"
	Feature geojson.Feature `json:"feature"`

	// Обычное поле
	UserID string `json:"user_id"`
}

// ListLocationsResponse ответ со списком локаций
type ListLocationsResponse struct {
	// Слайс структур с переопределением
	// @oa:rewrite.ref "github.com/arkannsk/elval/test/integration/geo_polymorphism/pkg/docs.FeatureDocs"
	Features []geojson.Feature `json:"features"`
}
