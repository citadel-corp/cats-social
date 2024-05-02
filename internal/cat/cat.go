package cat

import "time"

type CatRace string
type CatSex string

const (
	Persian          CatRace = "Persian"
	MaineCoon        CatRace = "Maine Coon"
	Siamese          CatRace = "Siamese"
	Ragdoll          CatRace = "Ragdoll"
	Bengal           CatRace = "Bengal"
	Sphynx           CatRace = "Sphynx"
	BritishShorthair CatRace = "British Shorthair"
	Abyssinian       CatRace = "Abyssinian"
	ScottishFold     CatRace = "Scottish Fold"
	Birman           CatRace = "Birman"

	Male   CatSex = "male"
	Female CatSex = "female"
)

var (
	CatRaces []interface{} = []interface{}{Persian, MaineCoon, Siamese, Ragdoll, Bengal, Sphynx, BritishShorthair, Abyssinian, ScottishFold, Birman}
	CatSexes []interface{} = []interface{}{Male, Female}
)

type Cat struct {
	ID          string
	UserID      string
	Name        string
	Race        CatRace
	Sex         CatSex
	Age         int
	Description string
	HasMatched  bool
	ImageURLS   []string
	CreatedAt   time.Time
}