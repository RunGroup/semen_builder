package builder

const DESC = "DESC"
const ASC = "ASC"

type Order struct {
	Direction string
	Field     string
}
