package filtertitle

import (
	"errors"
	"testing"
)

type titlePairs struct {
	title    string
	location []string
	err      error
}

var titles = []titlePairs{
	{title: "2015-02-15 11:08, Inbrott, Östra Göinge", location: []string{"Östra Göinge"}},
	{title: "2015-03-01 19:00, Trafikolycka, vilt, Östra Göinge", location: []string{"Östra Göinge"}},
	{title: "2015-02-28 22:57, Trafikolycka, smitning från, Hässleholm", location: []string{"Hässleholm"}},
	{title: "2015-01-20 10:28, Trafikolycka, Landskrona", location: []string{"Landskrona"}},
	{title: "2015-02-15 00:49, Sammanfattning kväll, Norrbotten", err: errors.New("Titeln är av typen 'Sammanfattning' och innehåller ingen platsinformation")},
	{title: "2015-02-24 20:21, Trafikolycka, singel, Helsingborg", location: []string{"Helsingborg"}},
}

func Test_FilterTitleWords(t *testing.T) {
	testedAmount := 0
	failedTests := 0

	for _, title := range titles {
		testedAmount++

		location, err := FilterTitleWords(title.title)

		if title.err == nil {
			if len(location) != len(title.location) {
				failedTests++
				t.Error(
					"\n\t\tFor:   \t ", title.title,
					"\n\t\tExpected:", len(title.location), " ", title.location,
					"\n\t\tGot:   \t ", len(location), " ", location,
				)
			} else {
				for i := 0; i < len(location); i++ {
					if location[i] != title.location[i] {
						failedTests++
						t.Error(
							"\n\t\tFor:   \t ", title.title,
							"\n\t\tExpected:", len(title.location), " ", title.location,
							"\n\t\tGot:   \t ", len(location), " ", location,
						)
					}
				}
			}
		} else if title.err.Error() != err.Error() {
			failedTests++
			t.Error(
				"\n\t\tFor:   \t ", title.title,
				"\n\t\tExpected: ", title.err.Error(),
				"\n\t\tGot:   \t ", err.Error(),
			)
		}
	}
	t.Log("--------------------------------------------")
	t.Logf("Failed tests: %v of %v ", failedTests, testedAmount)
	// t.Log("Number of tests failed: ", failedTests)
}
