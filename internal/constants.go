package internal

const (
	//config
	APP_PORT       string = "APP_PORT"
	MONGO_HOST     string = "MONGO_HOST"
	MONGO_DATABASE string = "MONGO_DATABASE"
	MONGO_PORT     string = "MONGO_PORT"
	//messages
	MsgResponseStartProcess     string = "Start process -->   "
	MsgResponseStartApplication string = "Starting the application...."
	MsgResponseConnectedMongoDB string = "Connected to MongoDB!"
	MsgResponseStartingNow      string = "STARTING NOW..."
	MsgResponseListingAll       string = "LISTANDO TODOS..."
	MsgResponseGettingOne       string = "OBTENIENDO UN OBJETO..."
	MsgResponseDeletingOne      string = "OBJETO BORRADO CORRECTAMENTE"
	MsgResponseCreatingOne      string = "OBJETO INSERTADO CORRECTAMENTE"
	MsgResponseUpdatingOne      string = "OBJETO EDITADO CORRECTAMENTE"
	MsgResponseObjectExists     string = "OBJETO EXISTENTE"
	MsgResponseServerError      string = "ERROR DE SERVIDOR"
	MsgResponseNoData           string = "SIN DATOS EXISTENTES"
	MsgResponseInvalidRequest   string = "PAYLOAD/REQUEST INVALIDO"
	//URLs
	URLStartingNow string = "/"
	URLApi         string = "/api"
	URLIndex       string = "/index"
	URLListingAll  string = "/list_persons"
	URLGettingOne  string = "/get_person"
	URLDeletingOne string = "/delete_person"
	URLCreatingOne string = "/create_person"
	URLUpdatingOne string = "/update_person"
	//types responses
	ERROR   string = "ERROR"
	SUCCESS string = "SUCCESS"
	TEST    string = "TEST"
	OPTIONS string = "OPTIONS"
	//Database
	CollectionPerson string = "person"
	Field__id        string = "_id"
	FieldUpdated     string = "updated"
	FieldLastname    string = "lastname"
	FieldCi          string = "ci"
	OrderDesc        string = "-"
	OrderAsc         string = ""
	MongoDB__set     string = "$set"
	//functions
	FuncCreatePerson string = "CreatePerson"
	FuncUpdatePerson string = "UpdatePerson"
	FuncListPersons  string = "ListPersons"
	FuncGetPerson    string = "GetPerson"
	FuncDeletePerson string = "DeletePerson"
	//keyvals
	KeyType               string = "type"
	KeyURL                string = "URL"
	KeyMessage            string = "message"
	KeyMethod             string = "method"
	KeyTook               string = "took"
	KeyValues             string = "values"
	KeyResponseMessage    string = "Message"
	KeyResponseStatusCode string = "StatusCode"
	KeyResponseData       string = "Data"
	//http methods
	HTTP_POST string = "POST"
	HTTP_GET  string = "GET"
	//misc
	ValueEmpty string = ""
)
