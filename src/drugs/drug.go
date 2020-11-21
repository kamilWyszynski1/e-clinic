package drugs

import "encoding/xml"

type ProduktyLecznicze struct {
	XMLName          xml.Name `xml:"produktyLecznicze"`
	Text             string   `xml:",chardata"`
	StanNaDzien      string   `xml:"stanNaDzien,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	ProduktLeczniczy []struct {
		Text                      string `xml:",chardata"`
		NazwaProduktu             string `xml:"nazwaProduktu,attr"`
		RodzajPreparatu           string `xml:"rodzajPreparatu,attr"`
		NazwaPowszechnieStosowana string `xml:"nazwaPowszechnieStosowana,attr"`
		Moc                       string `xml:"moc,attr"`
		Postac                    string `xml:"postac,attr"`
		PodmiotOdpowiedzialny     string `xml:"podmiotOdpowiedzialny,attr"`
		TypProcedury              string `xml:"typProcedury,attr"`
		NumerPozwolenia           string `xml:"numerPozwolenia,attr"`
		WaznoscPozwolenia         string `xml:"waznoscPozwolenia,attr"`
		KodATC                    string `xml:"kodATC,attr"`
		ID                        string `xml:"id,attr"`
		SubstancjeCzynne          struct {
			Text             string   `xml:",chardata"`
			SubstancjaCzynna []string `xml:"substancjaCzynna"`
		} `xml:"substancjeCzynne"`
		Opakowania struct {
			Text       string `xml:",chardata"`
			Opakowanie []struct {
				Text                  string `xml:",chardata"`
				Wielkosc              string `xml:"wielkosc,attr"`
				JednostkaWielkosci    string `xml:"jednostkaWielkosci,attr"`
				KodEAN                string `xml:"kodEAN,attr"`
				KategoriaDostepnosci  string `xml:"kategoriaDostepnosci,attr"`
				Skasowane             string `xml:"skasowane,attr"`
				NumerEu               string `xml:"numerEu,attr"`
				DystrybutorRownolegly string `xml:"dystrybutorRownolegly,attr"`
				ID                    string `xml:"id,attr"`
			} `xml:"opakowanie"`
		} `xml:"opakowania"`
	} `xml:"produktLeczniczy"`
}
