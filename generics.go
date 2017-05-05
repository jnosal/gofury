package fury

import "net/http"

type IRetrieve interface {
	Retrieve() (s string, err error)
}

type ICreate interface {
	Create() (err error)
}

type IList interface {
	List() (err error)
}

type IDestroy interface {
	Remove() (err error)
}

type IUpdate interface {
	Update() (err error)
}

type IModify interface {
	Modify() (err error)
}

type IListCreate interface {
	IList
	ICreate
}

type IModifyUpdate interface {
	IModify
	IUpdate
}

type IRetrieveUpdate interface {
	IRetrieve
	IModify
	IUpdate
}

type IDetailDestroy interface {
	IRetrieve
	IDestroy
}

func RetrieveResource(i IRetrieve, meta *Meta) {
	result, _ := i.Retrieve()
	data := map[string]string{"result": result}
	meta.Json(http.StatusOK, data)
}

func CreateResource(i ICreate, meta *Meta) {

}

func ListResource(i IList, meta *Meta) {

}

func RemoveResource(i IDestroy, meta *Meta) {

}

func UpdateResource(i IUpdate, meta *Meta) {

}

func ModifyResource(i IUpdate, meta *Meta) {

}
