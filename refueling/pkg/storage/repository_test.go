package storage

import (
	"fmt"
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

var (
	refName = 0
	actName = "0"
)

func TestDataCreatePopulate(t *testing.T) { //* Test Adding
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
				Description:   "Initial",
				RefuelNameRef: 1,
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
				Description:   "Ref 4-3 to 4-4",
				RefuelNameRef: 1,
			},
		},
	}

	got := PopulatePathTo(r1.RefuelName, r1.Acts)

	if got != nil {
		t.Errorf("got %q", got)
	}
}

func TestAddingAct(t *testing.T) {
	var act adding.Act = adding.Act{
		Name: "2",
		CoreConfig: [][]string{
			{"2123131654984132163131641313212\n"},
			{"9797764631154699333144477855235\n"},
		},
		PDC: []string{
			"gsljglsjgsjgjsgjslgjsj",
			"vmx.vm.x,mvlxjljjgoisg",
		}, //todo add few more lines
		Description:   "Ref 4-3 to 4-4",
		RefuelNameRef: 1,
	}

	var acts []adding.Act
	acts = append(acts, act)

	got := PopulatePathTo(act.RefuelNameRef, acts)
	if got != nil {
		t.Errorf("got an error %v", got)
	}
}

func TestStoredFileName(t *testing.T) {
	extension := "PDC"
	got := GetFileName(refName, actName, extension)
	fmt.Println(got)
}

func TestPathToStorage(t *testing.T) {

}

func TestDataQueryConfig(t *testing.T) {
	got := ConfigStorageQuery(refName, actName)
	if len(got) == 0 {
		t.Errorf("got empy array %v", got)
	}

}

// func TestDataQueryPDC(t *testing.T) {
// 	refName := 0
// got := PDCStorageQuery(refName, actName) //* instance of storage
// 	if len(got) == 0 {
// 		t.Errorf("got empy array %v", got)
// 	}
// 	fmt.Println(got)
// }

func TestConfigPDCStoredDelete(t *testing.T) {
	got := ConfigPDCStoredDelete(refName, actName)
	if got != nil {
		t.Errorf("got %v", got)
	}
}

func TestRefuelStoredDelete(t *testing.T) {
	got := RefuelStoredDelete(1)
	if got != nil {
		t.Errorf("got %v", got)
	}
}
