package integration

var Blue = []byte(`[
  { "lat": -6.34842, "long": 106.82976 },
  { "lat": -6.34857, "long": 106.83071 },
  { "lat": -6.3487, "long": 106.83136 },
  { "lat": -6.34873, "long": 106.83145 },
  { "lat": -6.34888, "long": 106.83169 },
  { "lat": -6.34908, "long": 106.83186 },
  { "lat": -6.34928, "long": 106.83193 },
  { "lat": -6.34939, "long": 106.83194 },
  { "lat": -6.35025, "long": 106.83176 },
  { "lat": -6.35112, "long": 106.83157 },
  { "lat": -6.35129, "long": 106.83156 },
  { "lat": -6.35142, "long": 106.83157 },
  { "lat": -6.35167, "long": 106.83165 },
  { "lat": -6.35203, "long": 106.83182 },
  { "lat": -6.35225, "long": 106.83189 },
  { "lat": -6.35248, "long": 106.83191 },
  { "lat": -6.35299, "long": 106.83189 },
  { "lat": -6.35402, "long": 106.8317 },
  { "lat": -6.35428, "long": 106.83161 },
  { "lat": -6.35458, "long": 106.83141 },
  { "lat": -6.35497, "long": 106.83107 },
  { "lat": -6.35587, "long": 106.83036 },
  { "lat": -6.35602, "long": 106.83027 },
  { "lat": -6.35614, "long": 106.83022 },
  { "lat": -6.35634, "long": 106.83018 },
  { "lat": -6.35649, "long": 106.83018 },
  { "lat": -6.35671, "long": 106.83021 },
  { "lat": -6.35698, "long": 106.83032 },
  { "lat": -6.35746, "long": 106.83071 },
  { "lat": -6.35795, "long": 106.83108 },
  { "lat": -6.35813, "long": 106.83118 },
  { "lat": -6.35828, "long": 106.83122 },
  { "lat": -6.35846, "long": 106.83124 },
  { "lat": -6.35887, "long": 106.83129 },
  { "lat": -6.36036, "long": 106.83152 },
  { "lat": -6.3622, "long": 106.83182 },
  { "lat": -6.36277, "long": 106.8319 },
  { "lat": -6.36318, "long": 106.83195 },
  { "lat": -6.36325, "long": 106.83192 },
  { "lat": -6.36328, "long": 106.83188 },
  { "lat": -6.3633, "long": 106.83178 },
  { "lat": -6.36327, "long": 106.83163 },
  { "lat": -6.3631, "long": 106.83145 },
  { "lat": -6.36149, "long": 106.82989 },
  { "lat": -6.36067, "long": 106.82907 },
  { "lat": -6.36033, "long": 106.82871 },
  { "lat": -6.35977, "long": 106.82815 },
  { "lat": -6.3597, "long": 106.82806 },
  { "lat": -6.35963, "long": 106.82791 },
  { "lat": -6.35961, "long": 106.82784 },
  { "lat": -6.35957, "long": 106.82616 },
  { "lat": -6.35957, "long": 106.82559 },
  { "lat": -6.35957, "long": 106.8253 },
  { "lat": -6.35961, "long": 106.82508 },
  { "lat": -6.35969, "long": 106.82491 },
  { "lat": -6.36067, "long": 106.82392 },
  { "lat": -6.36139, "long": 106.82317 },
  { "lat": -6.36177, "long": 106.82278 },
  { "lat": -6.36189, "long": 106.82269 },
  { "lat": -6.36208, "long": 106.82261 },
  { "lat": -6.3624, "long": 106.82251 },
  { "lat": -6.36338, "long": 106.8222 },
  { "lat": -6.36505, "long": 106.82167 },
  { "lat": -6.36548, "long": 106.82157 },
  { "lat": -6.36572, "long": 106.82157 },
  { "lat": -6.36582, "long": 106.82162 },
  { "lat": -6.36625, "long": 106.82196 },
  { "lat": -6.36657, "long": 106.82232 },
  { "lat": -6.36675, "long": 106.82259 },
  { "lat": -6.36682, "long": 106.82281 },
  { "lat": -6.36683, "long": 106.82316 },
  { "lat": -6.36684, "long": 106.82385 },
  { "lat": -6.36684, "long": 106.82385 },
  { "lat": -6.36685, "long": 106.8249 },
  { "lat": -6.3669, "long": 106.82497 },
  { "lat": -6.36702, "long": 106.82508 },
  { "lat": -6.36709, "long": 106.82511 },
  { "lat": -6.36856, "long": 106.82509 },
  { "lat": -6.36873, "long": 106.82511 },
  { "lat": -6.36895, "long": 106.82517 },
  { "lat": -6.36905, "long": 106.82523 },
  { "lat": -6.36953, "long": 106.82565 },
  { "lat": -6.37004, "long": 106.82613 },
  { "lat": -6.37091, "long": 106.82697 },
  { "lat": -6.37127, "long": 106.82729 },
  { "lat": -6.37135, "long": 106.82737 },
  { "lat": -6.37147, "long": 106.82757 },
  { "lat": -6.37152, "long": 106.82767 },
  { "lat": -6.37154, "long": 106.8278 },
  { "lat": -6.37155, "long": 106.82867 },
  { "lat": -6.37158, "long": 106.82928 },
  { "lat": -6.37159, "long": 106.82976 },
  { "lat": -6.37156, "long": 106.82991 },
  { "lat": -6.37136, "long": 106.83025 },
  { "lat": -6.37123, "long": 106.83046 },
  { "lat": -6.37113, "long": 106.83055 },
  { "lat": -6.37091, "long": 106.83073 },
  { "lat": -6.3708, "long": 106.83079 },
  { "lat": -6.37066, "long": 106.83086 },
  { "lat": -6.37051, "long": 106.83088 },
  { "lat": -6.37038, "long": 106.83089 },
  { "lat": -6.37008, "long": 106.83089 },
  { "lat": -6.37008, "long": 106.83089 },
  { "lat": -6.36913, "long": 106.83093 },
  { "lat": -6.36896, "long": 106.83094 },
  { "lat": -6.36882, "long": 106.831 },
  { "lat": -6.36859, "long": 106.83118 },
  { "lat": -6.36811, "long": 106.83171 },
  { "lat": -6.36782, "long": 106.83199 },
  { "lat": -6.36769, "long": 106.83204 },
  { "lat": -6.36709, "long": 106.83207 },
  { "lat": -6.36489, "long": 106.8321 },
  { "lat": -6.36447, "long": 106.8321 },
  { "lat": -6.36423, "long": 106.83207 },
  { "lat": -6.3637, "long": 106.83197 },
  { "lat": -6.36353, "long": 106.83187 },
  { "lat": -6.36335, "long": 106.83177 },
  { "lat": -6.36323, "long": 106.83173 },
  { "lat": -6.36308, "long": 106.83168 },
  { "lat": -6.36296, "long": 106.83166 },
  { "lat": -6.3627, "long": 106.83166 },
  { "lat": -6.36251, "long": 106.83174 },
  { "lat": -6.36234, "long": 106.83175 },
  { "lat": -6.36197, "long": 106.8317 },
  { "lat": -6.36031, "long": 106.83143 },
  { "lat": -6.35881, "long": 106.83117 },
  { "lat": -6.35812, "long": 106.83107 },
  { "lat": -6.35795, "long": 106.83098 },
  { "lat": -6.35743, "long": 106.83058 },
  { "lat": -6.3571, "long": 106.8303 },
  { "lat": -6.35695, "long": 106.83021 },
  { "lat": -6.35671, "long": 106.83011 },
  { "lat": -6.35657, "long": 106.83009 },
  { "lat": -6.35628, "long": 106.8301 },
  { "lat": -6.35608, "long": 106.83015 },
  { "lat": -6.35595, "long": 106.83021 },
  { "lat": -6.35566, "long": 106.83043 },
  { "lat": -6.35436, "long": 106.83144 },
  { "lat": -6.35414, "long": 106.83155 },
  { "lat": -6.35382, "long": 106.83163 },
  { "lat": -6.35322, "long": 106.83175 },
  { "lat": -6.35284, "long": 106.8318 },
  { "lat": -6.35261, "long": 106.83183 },
  { "lat": -6.35239, "long": 106.83184 },
  { "lat": -6.35215, "long": 106.83178 },
  { "lat": -6.35198, "long": 106.83171 },
  { "lat": -6.35184, "long": 106.83164 },
  { "lat": -6.35143, "long": 106.83151 },
  { "lat": -6.35114, "long": 106.83149 },
  { "lat": -6.35102, "long": 106.83151 },
  { "lat": -6.3502, "long": 106.83169 },
  { "lat": -6.34952, "long": 106.83184 },
  { "lat": -6.34938, "long": 106.83185 },
  { "lat": -6.34922, "long": 106.83182 },
  { "lat": -6.34907, "long": 106.83175 },
  { "lat": -6.34896, "long": 106.83164 },
  { "lat": -6.34885, "long": 106.83148 },
  { "lat": -6.34879, "long": 106.83134 },
  { "lat": -6.34849, "long": 106.82971 },
  { "lat": -6.34834, "long": 106.82899 },
  { "lat": -6.34826, "long": 106.82901 },
  { "lat": -6.34838, "long": 106.82956 },
  { "lat": -6.34842, "long": 106.82976 }
]`)

var Red = []byte(`[
  { "lat": -6.34842, "long": 106.82976 },
  { "lat": -6.34857, "long": 106.83071 },
  { "lat": -6.3487, "long": 106.83136 },
  { "lat": -6.34873, "long": 106.83145 },
  { "lat": -6.34888, "long": 106.83169 },
  { "lat": -6.34908, "long": 106.83186 },
  { "lat": -6.34928, "long": 106.83193 },
  { "lat": -6.34939, "long": 106.83194 },
  { "lat": -6.35025, "long": 106.83176 },
  { "lat": -6.35112, "long": 106.83157 },
  { "lat": -6.35129, "long": 106.83156 },
  { "lat": -6.35142, "long": 106.83157 },
  { "lat": -6.35167, "long": 106.83165 },
  { "lat": -6.35203, "long": 106.83182 },
  { "lat": -6.35225, "long": 106.83189 },
  { "lat": -6.35248, "long": 106.83191 },
  { "lat": -6.35299, "long": 106.83189 },
  { "lat": -6.35402, "long": 106.8317 },
  { "lat": -6.35428, "long": 106.83161 },
  { "lat": -6.35458, "long": 106.83141 },
  { "lat": -6.35497, "long": 106.83107 },
  { "lat": -6.35587, "long": 106.83036 },
  { "lat": -6.35602, "long": 106.83027 },
  { "lat": -6.35614, "long": 106.83022 },
  { "lat": -6.35634, "long": 106.83018 },
  { "lat": -6.35649, "long": 106.83018 },
  { "lat": -6.35671, "long": 106.83021 },
  { "lat": -6.35698, "long": 106.83032 },
  { "lat": -6.35746, "long": 106.83071 },
  { "lat": -6.35795, "long": 106.83108 },
  { "lat": -6.35813, "long": 106.83118 },
  { "lat": -6.35828, "long": 106.83122 },
  { "lat": -6.35846, "long": 106.83124 },
  { "lat": -6.35887, "long": 106.83129 },
  { "lat": -6.36036, "long": 106.83152 },
  { "lat": -6.3622, "long": 106.83182 },
  { "lat": -6.36277, "long": 106.8319 },
  { "lat": -6.36337, "long": 106.83199 },
  { "lat": -6.36427, "long": 106.83215 },
  { "lat": -6.36463, "long": 106.83218 },
  { "lat": -6.36709, "long": 106.83215 },
  { "lat": -6.3675, "long": 106.83214 },
  { "lat": -6.36774, "long": 106.83212 },
  { "lat": -6.36793, "long": 106.83199 },
  { "lat": -6.36823, "long": 106.83171 },
  { "lat": -6.36877, "long": 106.83115 },
  { "lat": -6.36884, "long": 106.83109 },
  { "lat": -6.36904, "long": 106.83105 },
  { "lat": -6.36989, "long": 106.83099 },
  { "lat": -6.3703, "long": 106.83098 },
  { "lat": -6.37049, "long": 106.83097 },
  { "lat": -6.37077, "long": 106.83091 },
  { "lat": -6.37091, "long": 106.83083 },
  { "lat": -6.37112, "long": 106.83066 },
  { "lat": -6.3713, "long": 106.83046 },
  { "lat": -6.37152, "long": 106.8301 },
  { "lat": -6.37164, "long": 106.82989 },
  { "lat": -6.37166, "long": 106.82974 },
  { "lat": -6.37166, "long": 106.82928 },
  { "lat": -6.37163, "long": 106.82854 },
  { "lat": -6.37162, "long": 106.82775 },
  { "lat": -6.37159, "long": 106.82764 },
  { "lat": -6.37154, "long": 106.82752 },
  { "lat": -6.37141, "long": 106.82731 },
  { "lat": -6.37124, "long": 106.82713 },
  { "lat": -6.36955, "long": 106.82557 },
  { "lat": -6.3692, "long": 106.82525 },
  { "lat": -6.36897, "long": 106.82509 },
  { "lat": -6.36875, "long": 106.82502 },
  { "lat": -6.36857, "long": 106.825 },
  { "lat": -6.36722, "long": 106.825 },
  { "lat": -6.36712, "long": 106.82495 },
  { "lat": -6.36699, "long": 106.82482 },
  { "lat": -6.36694, "long": 106.82471 },
  { "lat": -6.36692, "long": 106.82462 },
  { "lat": -6.36693, "long": 106.82404 },
  { "lat": -6.36692, "long": 106.82296 },
  { "lat": -6.36689, "long": 106.82275 },
  { "lat": -6.36681, "long": 106.82252 },
  { "lat": -6.36662, "long": 106.82227 },
  { "lat": -6.36642, "long": 106.82203 },
  { "lat": -6.36622, "long": 106.82181 },
  { "lat": -6.36609, "long": 106.82171 },
  { "lat": -6.36609, "long": 106.82171 },
  { "lat": -6.36578, "long": 106.82152 },
  { "lat": -6.36567, "long": 106.82149 },
  { "lat": -6.36553, "long": 106.82148 },
  { "lat": -6.36523, "long": 106.82154 },
  { "lat": -6.36323, "long": 106.82216 },
  { "lat": -6.36288, "long": 106.82227 },
  { "lat": -6.36214, "long": 106.8225 },
  { "lat": -6.36188, "long": 106.82257 },
  { "lat": -6.36172, "long": 106.82265 },
  { "lat": -6.36146, "long": 106.82271 },
  { "lat": -6.36135, "long": 106.82282 },
  { "lat": -6.36129, "long": 106.82289 },
  { "lat": -6.36124, "long": 106.82297 },
  { "lat": -6.36123, "long": 106.823 },
  { "lat": -6.36126, "long": 106.82305 },
  { "lat": -6.36111, "long": 106.82331 },
  { "lat": -6.36089, "long": 106.82357 },
  { "lat": -6.35971, "long": 106.82476 },
  { "lat": -6.35958, "long": 106.82493 },
  { "lat": -6.35952, "long": 106.82508 },
  { "lat": -6.35949, "long": 106.8253 },
  { "lat": -6.35949, "long": 106.82589 },
  { "lat": -6.35951, "long": 106.82741 },
  { "lat": -6.35952, "long": 106.82775 },
  { "lat": -6.35955, "long": 106.8279 },
  { "lat": -6.35961, "long": 106.82806 },
  { "lat": -6.35984, "long": 106.82835 },
  { "lat": -6.36062, "long": 106.82916 },
  { "lat": -6.36221, "long": 106.83072 },
  { "lat": -6.36265, "long": 106.83116 },
  { "lat": -6.36271, "long": 106.83125 },
  { "lat": -6.36275, "long": 106.83137 },
  { "lat": -6.3627, "long": 106.83166 },
  { "lat": -6.36251, "long": 106.83174 },
  { "lat": -6.36234, "long": 106.83175 },
  { "lat": -6.36197, "long": 106.8317 },
  { "lat": -6.36031, "long": 106.83143 },
  { "lat": -6.35881, "long": 106.83117 },
  { "lat": -6.35812, "long": 106.83107 },
  { "lat": -6.35795, "long": 106.83098 },
  { "lat": -6.35743, "long": 106.83058 },
  { "lat": -6.3571, "long": 106.8303 },
  { "lat": -6.35695, "long": 106.83021 },
  { "lat": -6.35671, "long": 106.83011 },
  { "lat": -6.35657, "long": 106.83009 },
  { "lat": -6.35628, "long": 106.8301 },
  { "lat": -6.35608, "long": 106.83015 },
  { "lat": -6.35595, "long": 106.83021 },
  { "lat": -6.35566, "long": 106.83043 },
  { "lat": -6.35436, "long": 106.83144 },
  { "lat": -6.35414, "long": 106.83155 },
  { "lat": -6.35382, "long": 106.83163 },
  { "lat": -6.35322, "long": 106.83175 },
  { "lat": -6.35284, "long": 106.8318 },
  { "lat": -6.35261, "long": 106.83183 },
  { "lat": -6.35239, "long": 106.83184 },
  { "lat": -6.35215, "long": 106.83178 },
  { "lat": -6.35198, "long": 106.83171 },
  { "lat": -6.35184, "long": 106.83164 },
  { "lat": -6.35143, "long": 106.83151 },
  { "lat": -6.35114, "long": 106.83149 },
  { "lat": -6.35102, "long": 106.83151 },
  { "lat": -6.3502, "long": 106.83169 },
  { "lat": -6.34952, "long": 106.83184 },
  { "lat": -6.34938, "long": 106.83185 },
  { "lat": -6.34922, "long": 106.83182 },
  { "lat": -6.34907, "long": 106.83175 },
  { "lat": -6.34896, "long": 106.83164 },
  { "lat": -6.34885, "long": 106.83148 },
  { "lat": -6.34879, "long": 106.83134 },
  { "lat": -6.34849, "long": 106.82971 },
  { "lat": -6.34834, "long": 106.82899 },
  { "lat": -6.34826, "long": 106.82901 },
  { "lat": -6.34838, "long": 106.82956 },
  { "lat": -6.34842, "long": 106.82976 }
]`)

var RedTest = []byte(`[
	{ "lat": -6.34842, "long": 106.82976 },
	{ "lat": -6.34857, "long": 106.83071 },
	{ "lat": -6.3487, "long": 106.83136 },
	{ "lat": -6.34873, "long": 106.83145 }
]`)

type RouteData struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}
