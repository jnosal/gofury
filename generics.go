package fury

import "net/http"

type IRetrieve interface {
	Retrieve() (s string, err error)
}

type ICreate interface {
	Create() (s string, err error)
}

type IList interface {
	List() (s string, err error)
}

type IDestroy interface {
	Remove() (s string, err error)
}

type IUpdate interface {
	Update() (s string, err error)
}

type IModify interface {
	Modify() (s string, err error)
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
	if result, err := i.Retrieve(); err == nil {
		data := map[string]string{"result": result}
		meta.Json(http.StatusOK, data)
	} else {
		data := map[string]string{"error": "Not found"}
		meta.Json(http.StatusNotFound, data)
	}
}

func CreateResource(i ICreate, meta *Meta) {
	if result, err := i.Create(); err == nil {
		data := map[string]string{"result": result}
		meta.Json(http.StatusCreated, data)
	} else {
		data := map[string]string{"error": "Validation Error"}
		meta.Json(http.StatusBadRequest, data)
	}
}

func ListResource(i IList, meta *Meta) {
	if result, err := i.List(); err == nil {
		data := map[string]string{"result": result}
		meta.Json(http.StatusOK, data)
	} else {
		data := map[string]string{"error": "Not found"}
		meta.Json(http.StatusNotFound, data)
	}
}

func RemoveResource(i IDestroy, meta *Meta) {
	if _, err := i.Remove(); err == nil {
		meta.Json(http.StatusNoContent, "")
	} else {
		data := map[string]string{"error": "Validation Error"}
		meta.Json(http.StatusBadRequest, data)
	}
}

func UpdateResource(i IUpdate, meta *Meta) {
	if _, err := i.Update(); err == nil {
		meta.Json(http.StatusNoContent, "")
	} else {
		data := map[string]string{"error": "Validation Error"}
		meta.Json(http.StatusBadRequest, data)
	}
}

func ModifyResource(i IModify, meta *Meta) {
	if _, err := i.Modify(); err == nil {
		meta.Json(http.StatusNoContent, "")
	} else {
		data := map[string]string{"error": "Validation Error"}
		meta.Json(http.StatusBadRequest, data)
	}
}
