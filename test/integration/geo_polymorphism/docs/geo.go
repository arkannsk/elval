package docs

//go:generate elval-gen gen -input . -openapi

// PointDocs документирует geojson.Point
// @oa:description "A single geographic coordinate [lon, lat]"
type PointDocs struct {
	// @oa:enum "Point"
	Type string `json:"type"`

	// @oa:title "Coordinates"
	// @oa:description "[longitude, latitude]"
	// @oa:example [37.6173, 55.7558]
	Coordinates [2]float64 `json:"coordinates"`
}

// PolygonDocs документирует geojson.Polygon
// @oa:description "A closed geometric shape"
type PolygonDocs struct {
	// @oa:enum "Polygon"
	Type string `json:"type"`

	// @oa:title "Coordinates"
	// @oa:description "Array of linear rings"
	Coordinates [][3]float64 `json:"coordinates"`
}

// FeatureDocs документирует geojson.Feature
// @oa:description "GeoJSON Feature with geometry and properties"
// @oa:discriminator.propertyName "geometry.type"
// @oa:discriminator.mapping "Point:PointFeature"
// @oa:discriminator.mapping "Polygon:PolygonFeature"
type FeatureDocs struct {
	// @oa:title "Geometry"
	// @oa:description "The geometric shape"
	// @oa:oneOf "PointDocs,PolygonDocs"
	Geometry any `json:"geometry"`

	// @oa:title "Properties"
	// @oa:description "Arbitrary properties"
	Properties map[string]any `json:"properties,omitempty"`
}
