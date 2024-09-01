package model

type OrganizersInfo struct {
	ID               uint    `json:"id"`
	Name             string  `json:"name"`
	Phone            string  `json:"phone"`
	Description      string  `json:"description"`
	Address          Address `json:"address"`
	ImageProfilePath string  `json:"image_profile_path"`
	ImageCoverPath   string  `json:"image_cover_path"`
}

type GetOrganizer struct {
	ID               uint              `json:"id"`
	Name             string            `json:"name"`
	Phone            string            `json:"phone"`
	Description      string            `json:"description"`
	Address          Address           `json:"address"`
	ImageProfilePath string            `json:"image_profile_path"`
	ImageCoverPath   string            `json:"image_cover_path"`
	Compatition      []GetCompatitions `json:"compatition"`
}

type UpdateOrganizer struct {
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Description string  `json:"description"`
	Address     Address `json:"address"`
}