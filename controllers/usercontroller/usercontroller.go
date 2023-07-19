package usercontroller

import (
	"net/http"
	"sort"
	"strconv"
	"text/template"

	"github.com/taufiqoo/technical-test-berat/entities"
	"github.com/taufiqoo/technical-test-berat/libraries"
	"github.com/taufiqoo/technical-test-berat/models"
)

var validation = libraries.NewValidation()
var userModel = models.NewUserModel()

func Index(response http.ResponseWriter, request *http.Request) {

	user, _ := userModel.FindAll()

	//buat mengurutkan dari tanggal terbaru
	sort.Slice(user, func(i, j int) bool {
		return user[i].Tanggal > user[j].Tanggal
	})

	// buat menghitung rata-rata
	avgMax, avgMin, avgDiff := userModel.CalculateAverages()

	data := map[string]interface{}{
		"user":    user,
		"avgMax":  avgMax,
		"avgMin":  avgMin,
		"avgDiff": avgDiff,
	}

	temp, err := template.ParseFiles("views/user/index.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(response, data)
}

func Add(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {
		temp, err := template.ParseFiles("views/user/add.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(response, nil)
	} else if request.Method == http.MethodPost {

		request.ParseForm()

		var user entities.User
		user.Tanggal = request.Form.Get("tanggal")
		user.Max = request.Form.Get("max")
		user.Min = request.Form.Get("min")

		var data = make(map[string]interface{})

		vError := validation.Struct(user)

		if vError != nil {
			data["validation"] = vError
		} else {
			data["pesan"] = "Data user berhasil disimpan"
			userModel.Create(user)
		}

		temp, _ := template.ParseFiles("views/user/add.html")
		temp.Execute(response, data)
	}
}

func Edit(response http.ResponseWriter, request *http.Request) {

	if request.Method == http.MethodGet {

		queryString := request.URL.Query()
		id, _ := strconv.ParseInt(queryString.Get("id"), 10, 64)

		var user entities.User
		userModel.Find(id, &user)

		data := map[string]interface{}{
			"user": user,
		}

		temp, err := template.ParseFiles("views/user/edit.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(response, data)
	} else if request.Method == http.MethodPost {

		request.ParseForm()

		var user entities.User
		user.Id, _ = strconv.ParseInt(request.Form.Get("id"), 10, 64)
		user.Tanggal = request.Form.Get("tanggal")
		user.Max = request.Form.Get("max")
		user.Min = request.Form.Get("min")

		var data = make(map[string]interface{})

		vError := validation.Struct(user)

		if vError != nil {
			data["validation"] = vError
		} else {
			data["pesan"] = "Data user berhasil diperbarui"
			userModel.Update(&user)
		}

		temp, _ := template.ParseFiles("views/user/edit.html")
		temp.Execute(response, data)
	}

}

func Delete(response http.ResponseWriter, request *http.Request) {

	queryString := request.URL.Query()
	id, _ := strconv.ParseInt(queryString.Get("id"), 10, 64)

	userModel.Delete(id)

	http.Redirect(response, request, "/user", http.StatusSeeOther)

}

func Show(response http.ResponseWriter, request *http.Request) {
	queryString := request.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)
	if err != nil {
		http.Error(response, "ID Pengguna Tidak Valid", http.StatusBadRequest)
		return
	}

	var user entities.User
	err = userModel.FindById(id, &user)
	if err != nil {
		http.Error(response, "Pengguna Tidak Ditemukan", http.StatusNotFound)
		return
	}

	temp, err := template.ParseFiles("views/user/show.html")
	if err != nil {
		http.Error(response, "Kesalahan Server Internal", http.StatusInternalServerError)
		return
	}

	temp.Execute(response, user)
}
