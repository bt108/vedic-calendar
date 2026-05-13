package renderer

var TithiAbbrev = map[string]string{
	"Pratipat":  "Pra",
	"Dvitiya":   "Dvi",
	"Tritiya":   "Tri",
	"Caturthi":  "Cat",
	"Pancami":   "Pan",
	"Sasti":     "Sas",
	"Saptami":   "Sap",
	"Astami":    "Ast",
	"Navami":    "Nav",
	"Dasami":    "Das",
	"Ekadasi":   "Eka",
	"Dvadasi":   "Dva",
	"Trayodasi": "Tra",
	"Caturdasi": "Cad",
	"Purnima":   "Pur",
	"Amavasya":  "Ama",
}

var NaksatraAbbrev = map[string]string{
	"Asvini":          "Asv",
	"Bharani":         "Bha",
	"Krittika":        "Kri",
	"Rohini":          "Roh",
	"Mrigasira":       "Mri",
	"Ardra":           "Ard",
	"Punarvasu":       "Pun",
	"Pusyami":         "Pus",
	"Aslesa":          "Asl",
	"Magha":           "Mag",
	"Purva-phalguni":  "PPh",
	"Uttara-phalguni": "UPh",
	"Hasta":           "Has",
	"Citra":           "Cit",
	"Swati":           "Swa",
	"Visakha":         "Vis",
	"Anuradha":        "Anu",
	"Jyestha":         "Jye",
	"Mula":            "Mul",
	"Purva-asadha":    "PAs",
	"Uttara-asadha":   "UAs",
	"Sravana":         "Sra",
	"Dhanista":        "Dha",
	"Satabhisa":       "Sat",
	"Purva-bhadra":    "PBh",
	"Uttara-bhadra":   "UBh",
	"Revati":          "Rev",
}

func ShortenTithi(tithi string) string {
	for full, abbr := range TithiAbbrev {
		if len(tithi) >= len(full) && tithi[:len(full)] == full {
			return abbr
		}
	}
	if len(tithi) >= 3 {
		return tithi[:3]
	}
	return tithi
}

func ShortenNaksatra(naksatra string) string {
	if abbr, ok := NaksatraAbbrev[naksatra]; ok {
		return abbr
	}
	if len(naksatra) >= 3 {
		return naksatra[:3]
	}
	return naksatra
}
