package model

const (
	RegionFieldID           = "id"
	RegionFieldName         = "name"
	RegionFieldSlug         = "slug"
	RegionFieldCountry      = "country"
	RegionFieldDisplayOrder = "display_order"
	RegionFieldParentID     = "parent_id"
	RegionFieldLevel        = "level"
	RegionFieldLat          = "lat"
	RegionFieldLng          = "lng"
	RegionFieldPlaceID      = "place_id"
)

type Region struct {
	ID           int64   `db:"id"`
	Name         string  `db:"name"`
	Slug         string  `db:"slug"`
	Country      string  `db:"country"`
	DisplayOrder int8    `db:"display_order"`
	ParentID     int64   `db:"parent_id"`
	Level        int8    `db:"level"`
	Lat          float64 `db:"lat"`
	Lng          float64 `db:"lng"`
	PlaceID      string  `db:"place_id"`
}

func (b *Region) FieldMap() (fields []string, values []interface{}) {
	fields = []string{
		RegionFieldID,
		RegionFieldName,
		RegionFieldSlug,
		RegionFieldCountry,
		RegionFieldDisplayOrder,
		RegionFieldParentID,
		RegionFieldLevel,
		RegionFieldLat,
		RegionFieldLng,
		RegionFieldPlaceID,
	}

	values = []interface{}{
		&b.ID,
		&b.Name,
		&b.Slug,
		&b.Country,
		&b.DisplayOrder,
		&b.ParentID,
		&b.Level,
		&b.Lat,
		&b.Lng,
		&b.PlaceID,
	}
	return
}

func (*Region) TableName() string {
	return "regions"
}
