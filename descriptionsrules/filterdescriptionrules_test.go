package main

import "testing"

type descriptionPairs struct {
	description string
	location    []string
}

var descriptions = []descriptionPairs{
	{"Rökutveckling, Blåregnsgatan.", []string{"Blåregnsgatan"}},
	{"Inbrott i bostad, Sösdala.", []string{"Sösdala"}},
	{"Personbilar i kollision på E6.", []string{}},
	{"Bilbrand på Tians väg.", []string{"Tians", "väg"}},
	{"Butiksrån, Norra Grängesbergsgatan.", []string{"Norra", "Grängesbergsgatan"}},
	{"Flera fordon i kollision på E22, Stavröd.", []string{"Stavröd"}},
	{"Två män gripna misstänkta för stöld, Hamnen Trelleborg", []string{"Hamnen", "Trelleborg"}},
	{"Inbrott i villa, Vellinge-Månstorp.", []string{"Vellinge-Månstorp"}},
	{"Rån mot apotek, Gösta Lundhs gata.", []string{"Gösta", "Lundhs", "gata"}},
	{"Lastbil välter, Inre ringvägen, i höjd med Sege industriområde.", []string{"Inre", "ringvägen", "Sege"}},
	{"Lastbil välter, Inre gränd.", []string{"Inre", "gränd"}},
	{"Två bilar krockar, korsningen Klörupsvägen och Havrejordsvägen.", []string{"Klörupsvägen", "Havrejordsvägen"}}, //Hitta ord som innehåller "vägen"
	{"Två personbilar i sidokollision, E22, Bäckaskog.", []string{"Bäckaskog"}},
	{"Personrån.", []string{}},
	{"Stopp och kontroll av stulen bil, väg 111, Laröd.", []string{"111", "Laröd"}},
	{"Trafikolycka/rattfylleri, Vankivavägen, Hässleholm.", []string{"Vankivavägen", "Hässleholm"}},
	{"Trafifkolycka, Länsväg 769, mellan Skurup och Rydsgård.", []string{"769", "Skurup", "Rydsgård"}},
	{"Trafikolycka på Länsväg 1329, Norra Rörum.", []string{"1329", "Norra", "Rörum"}},
	{"Två bilar krockar, korsningen Skolgatan och Andra Avenyn.", []string{"Skolgatan", "Andra", "Avenyn"}},
	{"Person påkörd, Storgatan. Senare ändrad till sjukdom/olycksfall.", []string{"Storgatan"}},
	{"Inbrott i företagslokal, Fosie industriområde.", []string{"Fosie"}},
	{"Två personbilar i kollision, Hjalmar Brantings väg.", []string{"Hjalmar", "Brantings", "väg"}},
	{"Brand i sopstation på Ramels väg i Malmö", []string{"Ramels", "väg", "Malmö"}},
	{"Misshandel, Anna Lindhs plats.", []string{"Anna", "Lindhs", "plats"}},
	{"Misshandel på Onsala danska väg i Gottskär. Detta var inte bra.", []string{"Onsala", "Gottskär"}},
	{"Något hände på Lv 985, Hakarp, Huskvarna.", []string{"985", "Hakarp", "Huskvarna"}},
	{"Georg Lückligs väg och larm om ungdomar som bråkar. På platsen.", []string{"Lückligs", "väg"}},
	{"personbil-minibuss på vägen mellan Hinneryd-Traryd. I minibussen", []string{"Hinneryd-Traryd"}},
	{"Per Högströmsgatan", []string{"Högströmsgatan"}},
	{"Singelolycka, väg 1560, Tomelillabygden.", []string{"1560", "Tomelillabygden"}},
	{"Skit hände på Trumvägen område Hamptjärnsmoran Boden är en åra", []string{"Trumvägen", "Hamptjärnsmoran", "Boden"}},
}

func Test_Rule1(t *testing.T) {
	testedAmount := 0
	failedTests := 0

	for _, desc := range descriptions {
		testedAmount++

		location := Rule1(desc.description)

		if len(location) != len(desc.location) {
			failedTests++
			t.Error(
				"\n\t\tFor:   \t ", desc.description,
				"\n\t\tExpected:", len(desc.location), " ", desc.location,
				"\n\t\tGot:   \t ", len(location), " ", location,
			)
		} else {
			for i := 0; i < len(location); i++ {
				if location[i] != desc.location[i] {
					failedTests++
					t.Error(
						"\n\t\tFor:   \t ", desc.description,
						"\n\t\tExpected:", len(desc.location), " ", desc.location,
						"\n\t\tGot:   \t ", len(location), " ", location,
					)
				}
			}
		}
	}
	t.Log("--------------------------------------------")
	t.Logf("Failed tests: %v of %v ", failedTests, testedAmount)
	// t.Log("Number of tests failed: ", failedTests)
}
