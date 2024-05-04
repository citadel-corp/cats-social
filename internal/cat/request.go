package cat

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var imgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_\-\:x\=\(\)]+)*(\.(jpg|jpeg|png))?$`, s)
	return match
}, "image url is not valid")

type CreateUpdateCatPayload struct {
	Name        string   `json:"name"`
	Race        CatRace  `json:"race"`
	Sex         CatSex   `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	Description string   `json:"description"`
	ImageURLS   []string `json:"imageUrls"`
}

func (p CreateUpdateCatPayload) Validate() error {
	for i := range p.ImageURLS {
		if len(p.ImageURLS[i]) == 0 {
			return errors.New("tags must not be empty")
		}
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&p.Race, validation.Required, validation.Length(5, 15), validation.In(CatRaces...)),
		validation.Field(&p.Sex, validation.Required, validation.In(CatSexes...)),
		validation.Field(&p.AgeInMonth, validation.Required, validation.Min(1), validation.Max(120082)),
		validation.Field(&p.Description, validation.Required, validation.Length(1, 200)),
		validation.Field(&p.ImageURLS, validation.Required, validation.Length(1, 0), validation.Each(validation.Required, validation.NotNil, imgUrlValidationRule)),
	)
}
