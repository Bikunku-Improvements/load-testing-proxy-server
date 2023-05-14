package integration

type BikunCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	BlueCredential = []BikunCredential{
		{
			Username: "bikun1",
			Password: "bikun1",
		},
		{
			Username: "bikun4",
			Password: "bikun4",
		},
		{
			Username: "bikun5",
			Password: "bikun5",
		},
		{
			Username: "bikun6",
			Password: "bikun6",
		},
		{
			Username: "bikun3",
			Password: "bikun3",
		},
	}

	RedCredential = []BikunCredential{
		{
			Username: "bikun7",
			Password: "bikun7",
		},
		{
			Username: "bikun8",
			Password: "bikun8",
		},
		{
			Username: "bikun2",
			Password: "bikun2",
		},
		{
			Username: "bikun9",
			Password: "bikun9",
		},
		{
			Username: "bikun10",
			Password: "bikun10",
		},
		{
			Username: "bikun11",
			Password: "bikun11",
		},
	}
)
