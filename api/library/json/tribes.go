package json

type Tribes struct {
	Count int     `json:"total_count,omitempty"`
	Items []Tribe `json:"items,omitempty"`
}

type Tribe struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	UserAvatar string      `json:"user_avatar"`
	URL        string      `json:"url"`
	Groups     TribeGroups `json:"groups,omitempty"`
}

type TribeGroups struct {
	Count int          `json:"total_count,omitempty"`
	Items []TribeGroup `json:"items,omitempty"`
}

type TribeGroup struct {
	Id           GroupId `json:"id"`
	Name         string  `json:"name"`
	TotalObjects int     `json:"total_count_objects,omitempty"`
	Date         string  `json:"date,omitempty"`
}