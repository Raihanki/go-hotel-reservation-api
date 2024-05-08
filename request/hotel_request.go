package request

type HotelRequest struct {
	Name        string `json:"name" form:"name" validate:"required,min=3,max=199"`
	Type        string `json:"type" form:"type" validate:"required,min=3,max=199"`
	City        string `json:"city" form:"city" validate:"required,min=3,max=199"`
	Address     string `json:"address" form:"address" validate:"required,min=3,max=500"`
	Photo       string `json:"photo" form:"photo"`
	Description string `json:"description" form:"description" validate:"required,min=3 max=1000"`
	Rating      string `json:"rating" form:"rating" validate:"min=1,max=5,number"`
}
