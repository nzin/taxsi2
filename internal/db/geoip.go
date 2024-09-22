package db

import "gorm.io/gorm"

type GeoipCountry struct {
	CountryCode string
	Allow       bool
}

/*
GetGeoipAllow returns
- an enmpty list (everything is allowed/deny)
- or a specifc list of country code, + allow
*/
func (ds *DbServiceImpl) GetGeoipCountries() (map[string]bool, error) {
	var list []GeoipCountry

	err := ds.db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	ccs := make(map[string]bool)
	for _, cc := range list {
		ccs[cc.CountryCode] = cc.Allow
	}
	return ccs, nil
}

func (ds *DbServiceImpl) SetGeoipCountries(countryCodes []string, allow bool) error {
	ds.db.Unscoped().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&GeoipCountry{})

	for _, cc := range countryCodes {
		c := GeoipCountry{
			CountryCode: cc,
			Allow:       allow,
		}
		err := ds.db.Create(&c).Error
		if err != nil {
			return err
		}
	}
	mode := "allow"
	if !allow {
		mode = "deny"
	}
	return ds.NotifyChange(CHANGELOG_TABLE_GEOIP, mode)
}

/*
Officially assigned ISO 3166 countries

AD	AND	Andorra
AE	ARE	United Arab Emirates (the)
AF	AFG	Afghanistan
AG	ATG	Antigua and Barbuda
AI	AIA	Anguilla
AL	ALB	Albania
AM	ARM	Armenia
AO	AGO	Angola
AQ	ATA	Antarctica
AR	ARG	Argentina
AS	ASM	American Samoa
AT	AUT	Austria
AU	AUS	Australia
AW	ABW	Aruba
AX	ALA	Åland Islands
AZ	AZE	Azerbaijan
BA	BIH	Bosnia and Herzegovina
BB	BRB	Barbados
BD	BGD	Bangladesh
BE	BEL	Belgium
BF	BFA	Burkina Faso
BG	BGR	Bulgaria
BH	BHR	Bahrain
BI	BDI	Burundi
BJ	BEN	Benin
BL	BLM	Saint Barthélemy
BM	BMU	Bermuda
BN	BRN	Brunei Darussalam
BO	BOL	Bolivia (Plurinational State of)
BQ	BES	Bonaire, Sint Eustatius and Saba
BR	BRA	Brazil
BS	BHS	Bahamas (the)
BT	BTN	Bhutan
BV	BVT	Bouvet Island
BW	BWA	Botswana
BY	BLR	Belarus
BZ	BLZ	Belize
CA	CAN	Canada
CC	CCK	Cocos (Keeling) Islands (the)
CD	COD	Congo (the Democratic Republic of the)
CF	CAF	Central African Republic (the)
CG	COG	Congo (the)
CH	CHE	Switzerland
CI	CIV	Côte d'Ivoire
CK	COK	Cook Islands (the)
CL	CHL	Chile
CM	CMR	Cameroon
CN	CHN	China
CO	COL	Colombia
CR	CRI	Costa Rica
CU	CUB	Cuba
CV	CPV	Cabo Verde
CW	CUW	Curaçao
CX	CXR	Christmas Island
CY	CYP	Cyprus
CZ	CZE	Czechia
DE	DEU	Germany
DJ	DJI	Djibouti
DK	DNK	Denmark
DM	DMA	Dominica
DO	DOM	Dominican Republic (the)
DZ	DZA	Algeria
EC	ECU	Ecuador
EE	EST	Estonia
EG	EGY	Egypt
EH	ESH	Western Sahara*
ER	ERI	Eritrea
ES	ESP	Spain
ET	ETH	Ethiopia
FI	FIN	Finland
FJ	FJI	Fiji
FK	FLK	Falkland Islands (the) [Malvinas]
FM	FSM	Micronesia (Federated States of)
FO	FRO	Faroe Islands (the)
FR	FRA	France
GA	GAB	Gabon
GB	GBR	United Kingdom of Great Britain and Northern Ireland (the)
GD	GRD	Grenada
GE	GEO	Georgia
GF	GUF	French Guiana
GG	GGY	Guernsey
GH	GHA	Ghana
GI	GIB	Gibraltar
GL	GRL	Greenland
GM	GMB	Gambia (the)
GN	GIN	Guinea
GP	GLP	Guadeloupe
GQ	GNQ	Equatorial Guinea
GR	GRC	Greece
GS	SGS	South Georgia and the South Sandwich Islands
GT	GTM	Guatemala
GU	GUM	Guam
GW	GNB	Guinea-Bissau
GY	GUY	Guyana
HK	HKG	Hong Kong
HM	HMD	Heard Island and McDonald Islands
HN	HND	Honduras
HR	HRV	Croatia
HT	HTI	Haiti
HU	HUN	Hungary
ID	IDN	Indonesia
IE	IRL	Ireland
IL	ISR	Israel
IM	IMN	Isle of Man
IN	IND	India
IO	IOT	British Indian Ocean Territory (the)
IQ	IRQ	Iraq
IR	IRN	Iran (Islamic Republic of)
IS	ISL	Iceland
IT	ITA	Italy
JE	JEY	Jersey
JM	JAM	Jamaica
JO	JOR	Jordan
JP	JPN	Japan
KE	KEN	Kenya
KG	KGZ	Kyrgyzstan
KH	KHM	Cambodia
KI	KIR	Kiribati
KM	COM	Comoros (the)
KN	KNA	Saint Kitts and Nevis
KP	PRK	Korea (the Democratic People's Republic of)
KR	KOR	Korea (the Republic of)
KW	KWT	Kuwait
KY	CYM	Cayman Islands (the)
KZ	KAZ	Kazakhstan
LA	LAO	Lao People's Democratic Republic (the)
LB	LBN	Lebanon
LC	LCA	Saint Lucia
LI	LIE	Liechtenstein
LK	LKA	Sri Lanka
LR	LBR	Liberia
LS	LSO	Lesotho
LT	LTU	Lithuania
LU	LUX	Luxembourg
LV	LVA	Latvia
LY	LBY	Libya
MA	MAR	Morocco
MC	MCO	Monaco
MD	MDA	Moldova (the Republic of)
ME	MNE	Montenegro
MF	MAF	Saint Martin (French part)
MG	MDG	Madagascar
MH	MHL	Marshall Islands (the)
MK	MKD	North Macedonia
ML	MLI	Mali
MM	MMR	Myanmar
MN	MNG	Mongolia
MO	MAC	Macao
MP	MNP	Northern Mariana Islands (the)
MQ	MTQ	Martinique
MR	MRT	Mauritania
MS	MSR	Montserrat
MT	MLT	Malta
MU	MUS	Mauritius
MV	MDV	Maldives
MW	MWI	Malawi
MX	MEX	Mexico
MY	MYS	Malaysia
MZ	MOZ	Mozambique
NA	NAM	Namibia
NC	NCL	New Caledonia
NE	NER	Niger (the)
NF	NFK	Norfolk Island
NG	NGA	Nigeria
NI	NIC	Nicaragua
NL	NLD	Netherlands (the)
NO	NOR	Norway
NP	NPL	Nepal
NR	NRU	Nauru
NU	NIU	Niue
NZ	NZL	New Zealand
OM	OMN	Oman
PA	PAN	Panama
PE	PER	Peru
PF	PYF	French Polynesia
PG	PNG	Papua New Guinea
PH	PHL	Philippines (the)
PK	PAK	Pakistan
PL	POL	Poland
PM	SPM	Saint Pierre and Miquelon
PN	PCN	Pitcairn
PR	PRI	Puerto Rico
PS	PSE	Palestine, State of
PT	PRT	Portugal
PW	PLW	Palau
PY	PRY	Paraguay
QA	QAT	Qatar
RE	REU	Réunion
RO	ROU	Romania
RS	SRB	Serbia
RU	RUS	Russian Federation (the)
RW	RWA	Rwanda
SA	SAU	Saudi Arabia
SB	SLB	Solomon Islands
SC	SYC	Seychelles
SD	SDN	Sudan (the)
SE	SWE	Sweden
SG	SGP	Singapore
SH	SHN	Saint Helena, Ascension and Tristan da Cunha
SI	SVN	Slovenia
SJ	SJM	Svalbard and Jan Mayen
SK	SVK	Slovakia
SL	SLE	Sierra Leone
SM	SMR	San Marino
SN	SEN	Senegal
SO	SOM	Somalia
SR	SUR	Suriname
SS	SSD	South Sudan
ST	STP	Sao Tome and Principe
SV	SLV	El Salvador
SX	SXM	Sint Maarten (Dutch part)
SY	SYR	Syrian Arab Republic (the)
SZ	SWZ	Eswatini
TC	TCA	Turks and Caicos Islands (the)
TD	TCD	Chad
TF	ATF	French Southern Territories (the)
TG	TGO	Togo
TH	THA	Thailand
TJ	TJK	Tajikistan
TK	TKL	Tokelau
TL	TLS	Timor-Leste
TM	TKM	Turkmenistan
TN	TUN	Tunisia
TO	TON	Tonga
TR	TUR	Türkiye
TT	TTO	Trinidad and Tobago
TV	TUV	Tuvalu
TW	TWN	Taiwan (Province of China)
TZ	TZA	Tanzania, the United Republic of
UA	UKR	Ukraine
UG	UGA	Uganda
UM	UMI	United States Minor Outlying Islands (the)
US	USA	United States of America (the)
UY	URY	Uruguay
UZ	UZB	Uzbekistan
VA	VAT	Holy See (the)
VC	VCT	Saint Vincent and the Grenadines
VE	VEN	Venezuela (Bolivarian Republic of)
VG	VGB	Virgin Islands (British)
VI	VIR	Virgin Islands (U.S.)
VN	VNM	Viet Nam
VU	VUT	Vanuatu
WF	WLF	Wallis and Futuna
WS	WSM	Samoa
YE	YEM	Yemen
YT	MYT	Mayotte
ZA	ZAF	South Africa
ZM	ZMB	Zambia
ZW	ZWE	Zimbabwe


Unofficial but widely used

XK	XKX	Kosovo
*/
