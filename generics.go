package fury

import "net/http"


type Retrievable interface {
	Retrieve() string
}


type Creatable interface {
	Create()
}


type Listable interface {
	List()
}


type Removable interface {
	Remove()
}


func RetrieveResource(retriever Retrievable, meta *Meta) {
	result := retriever.Retrieve()
	data := map[string]string{"result": result}
	meta.Json(http.StatusOK, data)
}


func CreateResource(creator Creatable, meta *Meta) {

}


func ListResource(lister Listable, meta *Meta) {

}


func DeleteResource(remover Removable, meta *Meta) {

}