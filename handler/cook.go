package handler

import (
	"cookvs/model"
	"cookvs/repository/users"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Cook struct {
	Repo *users.SqlRepo
}

func (o *Cook) CheckEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := model.User{
		Email: body.Email,
	}
	user1, err := o.Repo.CheckEmail(user)
	fmt.Println("Request Headers:", r.Header)
	fmt.Println("Request Body:", body)
	if user1 == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		fmt.Println("failed to Check", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(user1)
	if err != nil {
		fmt.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (o *Cook) UpdateByID(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Image    string `json:"image"`
		NickName string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	userid, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := model.User{
		Image:    body.Image,
		NickName: body.NickName,
	}

	err = o.Repo.Update(uint(userid), user)
	if err != nil {
		fmt.Println("failed to update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (o *Cook) UploadImageUser(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Ошибка при получении файла:", err)
		return
	}
	defer file.Close()

	// Указание пути для сохранения файла в папке "uploads"
	uploadPath := "./assets/images/"
	os.MkdirAll(uploadPath, os.ModePerm) // Создаем папку, если её нет

	// Создание нового файла на сервере
	f, err := os.OpenFile(uploadPath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer f.Close()

	// Копирование содержимого файла в созданный файл на сервере
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println("Ошибка при копировании файла:", err)
		return
	}

	fmt.Println("Файл успешно загружен на сервер.")
}

func (o *Cook) UploadImageRecipes(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Указание пути для сохранения файла в папке "uploads"
	uploadPath := "./assets/recipes/"
	os.MkdirAll(uploadPath, os.ModePerm) // Создаем папку, если её нет

	// Создание нового файла на сервере
	f, err := os.OpenFile(uploadPath+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Копирование содержимого файла в созданный файл на сервере
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Ошибка при копировании файла", http.StatusInternalServerError)
		return
	}

	fmt.Println("Файл успешно загружен на сервер.")
	w.WriteHeader(http.StatusOK)
}

func (o *Cook) FindByName(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recipe := model.Recipe{
		Name: body.Name,
	}
	recipe1, err := o.Repo.FindByName(recipe)
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = recipe1

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) FindByCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recipe := model.Recipe{
		Name:     body.Name,
		Category: body.Category,
	}
	recipe1, err := o.Repo.RecipeByCategory(recipe)
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = recipe1

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) FindByTag(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recipe := model.Recipe{
		Name: body.Name,
		Tag:  body.Tag,
	}
	recipe1, err := o.Repo.RecipeByTag(recipe)
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = recipe1

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) ListRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListRecipe")
	res, err := o.Repo.FindAllRecipe()
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = res

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateRecipe")
	var body struct {
		UserID      uint64          `json:"user_id"`
		Name        string          `json:"name"`
		Image       string          `json:"image"`
		Description []model.Step    `json:"description"`
		Products    []model.Product `json:"products"`
		Category    string          `json:"category"`
		Tag         string          `json:"tag"`
		Video       string          `json:"video"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("Decoded JSON Body: %+v\n", body)
		return
	}

	if len(body.Description) == 0 || len(body.Products) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("Description or Products is empty")
		return
	}

	recipes := model.Recipe{
		UserID:      body.UserID,
		Name:        body.Name,
		Image:       body.Image,
		Description: body.Description,
		Products:    body.Products,
		Category:    body.Category,
		Tag:         body.Tag,
		Video:       body.Video,
	}

	err := o.Repo.InsertRecipe(recipes)
	if err != nil {
		fmt.Println("failed to insert recipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(recipes)
	if err != nil {
		fmt.Println("failed to marshal recipe", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Cook) ListRecipeByCategory(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Category string `json:"category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recipe := model.Recipe{
		Category: body.Category,
	}
	recipe1, err := o.Repo.RecipeByCategory(recipe)
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = recipe1

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) ListRecipeByTag(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Tag string `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recipe := model.Recipe{
		Tag: body.Tag,
	}
	recipe1, err := o.Repo.RecipeByTag(recipe)
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = recipe1

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) FindByEmail(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := model.User{
		Email:    body.Email,
		Password: body.Password,
	}
	user1, err := o.Repo.Login(user)
	if user1 == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		fmt.Println("failed to Login", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(user1)
	if err != nil {
		fmt.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)

}

func (o *Cook) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create")
	var body struct {
		Image    string `json:"image"`
		Email    string `json:"email"`
		Password string `json:"password"`
		NickName string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := model.User{
		Image:    body.Image,
		Email:    body.Email,
		Password: body.Password,
		NickName: body.NickName,
	}

	err := o.Repo.Insert(users)
	if err != nil {
		fmt.Println("failed to insert", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(users)
	if err != nil {
		fmt.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (o *Cook) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List")
	res, err := o.Repo.FindAll()
	if err != nil {
		fmt.Println("failder to find all")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		User []model.User `json:"user"`
	}
	response.User = res

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get")
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64
	UserID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := o.Repo.FindById(uint(UserID))
	if err != nil {
		fmt.Println("failed to find by id ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (o *Cook) RecipeByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get")
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64
	RecipeID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe, err := o.Repo.RecipeById(uint(RecipeID))
	if err != nil {
		fmt.Println("failed to find recipe by id ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(recipe); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (o *Cook) ListRecipeByUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ListRecipeByUser")
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64
	UserID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := o.Repo.RecipeByUserId(uint(UserID))
	if err != nil {
		fmt.Println("failder to find all recipe")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response struct {
		Recipe []model.Recipe `json:"recipe"`
	}
	response.Recipe = res

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println("failed Marshal recipe:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (o *Cook) GetByEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get")
	idParam := chi.URLParam(r, "email")

	const base = 10
	const bitSize = 64
	UserID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := o.Repo.FindById(uint(UserID))
	if err != nil {
		fmt.Println("failed to find by id ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		fmt.Println("failed to marshal:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (o *Cook) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete")
	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	UserID, err := strconv.ParseUint(idParam, base, bitSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = o.Repo.DeleteById(uint(UserID))
	if err != nil {
		fmt.Println("failed to delete", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(UserID); err != nil {
		fmt.Println("failed to marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
