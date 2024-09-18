package enums

const (
	Autowired         string = "@Autowired"
	Entity                   = "@Entity"
	Controller               = "@Controller"
	Service                  = "@Service"
	Repository               = "@Repository"
	Column                   = "@Column(name = \"%s\")"
	Table                    = "@Table(name = \"%s\")"
	Id                       = "@Id"
	GeneratedValue           = "@GeneratedValue(%s)"
	RequestController        = "@RequestController"
	RequestMapping           = "@RequestMapping(%s)"
	GetMapping               = "@GetMapping%s"
	PostMapping              = "@PostMapping%s"
	DeleteMapping            = "@DeleteMapping%s"
	PutMapping               = "@PutMapping%s"
)
