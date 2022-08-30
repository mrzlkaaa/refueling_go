package storage

import (
	"refueling/refueling/pkg/adding"
	"testing"
	// "github.com/stretchr/testify"
)

//* struct to test adding to db
// r1 := Refuel{
// 	ID:         1,
// 	RefuelName: 1,
// 	Date:       "123456789",
// 	PathTo:     "",
// 	Acts: []Act{
// 		{
// 			ID:          1,
// 			Name:        "1",
// 			Description: "Initial",
// 			RefuelID:    1,
// 		},
// 		{
// 			ID:          2,
// 			Name:        "2",
// 			Description: "Ref 4-3 to 4-4",
// 			RefuelID:    1,
// 		},
// 	},
// }

func TestDataCreatePopulate(t *testing.T) {
	r1 := adding.Refuel{
		RefuelName: 1,
		Date:       "123456789",
		Acts: []adding.Act{
			{
				Name: "0",
				CoreConfig: [][]string{
					{"4442823842835823058092852385820"},
					{"0303030351568888946544646544542"},
				},
				PDC: []string{
					"sfsdfsdfsdfsfsfsdfs",
					"klmmbcmvbmcvbippbcp",
				}, //todo add few more lines
				Description: "Initial",
				RefuelID:    1,
			},
			{
				Name: "1",
				CoreConfig: [][]string{
					{"9898779798797979797897977799882\n"},
					{"9393939441717117121212187873215\n"},
				},
				PDC: []string{
					".,,,cvcnvcnmvncmvnvc",
					"yrhfghdbdkhgiudhgijvkl",
				}, //todo add few more lines
				Description: "Ref 4-3 to 4-4",
				RefuelID:    1,
			},
		},
	}

	got := PopulatePathTo(&r1)

	want := "/mnt/c/Users/Nikita/Desktop/codes/go/src/refueling/refueling/pkg/storage/data/1"
	if got != nil {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestDataQueryRefuel(t *testing.T) {
	//* return array of .txt names
	refName := 0
	gotRes := ConfigStorageQuery(refName)
	if len(gotRes) == 0 {
		t.Errorf("got empy array %v", gotRes)
	}

}
