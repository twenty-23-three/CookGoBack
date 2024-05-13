package users

import (
	"cookvs/model"
	"database/sql"
	"fmt"
	"math/rand"
)

type SqlRepo struct {
	DB *sql.DB
}

func RandomImagePath() string {
	min := 0
	max := 5
	img := rand.Intn(max-min) + min
	return fmt.Sprintf("http://localhost:3000/assets/images/%v.png", img)
}

func (r *SqlRepo) InsertComment(comments model.Comments) error {
	statement, err := r.DB.Prepare(`INSERT INTO ` + "`comments`" + ` (
		id_recipe,
		image_user,
		name_user,
		comment,
		date ) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert user: %w", err)
	}
	_, err = statement.Exec(comments.IdRecipe, comments.ImageUser, comments.NameUser, comments.Comment, comments.Date)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)

	}

	return nil
}
func (r *SqlRepo) RecipeByCountComments() ([]model.Comments, error) {
	array := []model.Comments{}

	rows, err := r.DB.Query(`
        SELECT number, id_recipe, image_user, name_user, comment, date
        FROM comments
        GROUP BY id_recipe
        ORDER BY COUNT(*) DESC
        LIMIT 1;
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare RecipeByCountComments recipe: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		model_comments := model.Comments{}

		err := rows.Scan(
			&model_comments.Number,
			&model_comments.IdRecipe,
			&model_comments.ImageUser,
			&model_comments.NameUser,
			&model_comments.Comment,
			&model_comments.Date,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		array = append(array, model_comments)
	}

	return array, nil
}

func (r *SqlRepo) CommentsByID(recipe_id uint) ([]model.Comments, error) {

	array := []model.Comments{}

	rows, err := r.DB.Query(`SELECT * FROM `+"`comments`"+`WHERE id_recipe = ? ORDER BY date`, recipe_id)
	if err != nil {
		return []model.Comments{}, fmt.Errorf("failed to prepare CommentsByID recipe: %w", err)
	}

	for rows.Next() {
		model_comments := model.Comments{}

		rows.Scan(&model_comments.Number, &model_comments.IdRecipe, &model_comments.ImageUser, &model_comments.NameUser, &model_comments.Comment, &model_comments.Date)
		array = append(array, model_comments)
	}
	return array, nil
}

func (r *SqlRepo) Insert(user model.User) error {
	statement, err := r.DB.Prepare(`INSERT INTO ` + "`users`" + ` (
		image,
		email,
		password,
		nickname ) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert user: %w", err)
	}
	_, err = statement.Exec(user.Image, user.Email, user.Password, user.NickName)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)

	}

	return nil
}

func (r *SqlRepo) CheckEmail(user model.User) (*model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM`+"`users`"+`
	WHERE email = ?`, user.Email)
	if err != nil {
		return &model_user, fmt.Errorf("failed to prepare find user: %w", err)
	}

	if rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName)
		return &model_user, nil
	}

	return nil, nil

}

func (r *SqlRepo) Login(user model.User) (*model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM`+"`users`"+`
	WHERE email = ? AND password =?`, user.Email, user.Password)
	if err != nil {
		return &model_user, fmt.Errorf("failed to prepare find user: %w", err)
	}

	if rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName)
		return &model_user, nil
	}

	return nil, nil

}

func (r *SqlRepo) FindById(user_id uint) (model.User, error) {
	model_user := model.User{}
	rows, err := r.DB.Query(`SELECT * FROM `+"`users`"+` WHERE user_id = ?`, user_id)
	if err != nil {
		return model_user, fmt.Errorf("failed to prepare find user: %w", err)
	}

	for rows.Next() {
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName)
	}

	return model_user, nil
}

func (r *SqlRepo) DeleteById(user_id uint) error {
	statement, err := r.DB.Prepare(`DELETE FROM ` + "`users`" + ` WHERE user_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare delete user: %w", err)
	}

	_, err = statement.Exec(user_id)
	if err != nil {
		return fmt.Errorf("failed to prepare delete EXEC user: %w", err)
	}
	return nil

}

func (r *SqlRepo) Update(user_id uint, model_user model.User) error {

	statement, err := r.DB.Prepare(`UPDATE ` + "`users`" + ` SET
		image = ?,
        nickname = ?
		WHERE user_id = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare update user: %w", err)
	}
	_, err = statement.Exec(model_user.Image, model_user.NickName, user_id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *SqlRepo) FindAll() ([]model.User, error) {

	array := []model.User{}

	rows, err := r.DB.Query(`SELECT * FROM ` + "`users`" + ``)
	if err != nil {
		return []model.User{}, fmt.Errorf("failed to prepare FindAll user: %w", err)
	}

	for rows.Next() {
		model_user := model.User{}
		rows.Scan(&model_user.UserID, &model_user.Image, &model_user.Email, &model_user.Password, &model_user.NickName)
		array = append(array, model_user)
	}
	return array, nil
}

func (r *SqlRepo) InsertRecipe(recipe model.Recipe) error {
	statement, err := r.DB.Prepare(`INSERT INTO ` + "`recipes`" + ` (
		user_id,
		name,
		image,
		description,
		products,
		category,
		tag,
		video) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert recipe: %w", err)
	}
	_, err = statement.Exec(recipe.UserID, recipe.Name, recipe.Image, recipe.MarshalDescription(), recipe.MarshalProducts(), recipe.Category, recipe.Tag, recipe.Video)
	if err != nil {
		return fmt.Errorf("failed to insert recipe: %w", err)

	}

	return nil
}

func (r *SqlRepo) FindAllRecipe() ([]model.Recipe, error) {

	array := []model.Recipe{}

	rows, err := r.DB.Query(`SELECT * FROM ` + "`recipes`" + `ORDER BY name`)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare FindAll recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		products := ""
		description := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag, &model_recipe.Video)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}
	return array, nil
}

func (r *SqlRepo) FindByName(recipe model.Recipe) ([]model.Recipe, error) {
	array := []model.Recipe{}
	var stroka string
	if recipe.Name != "" {
		stroka = `SELECT * FROM recipes WHERE name LIKE '` + recipe.Name + `%' ORDER BY name;`
	} else {
		stroka = `SELECT * FROM recipes ORDER BY name ASC;`
	}
	fmt.Println(stroka)
	rows, err := r.DB.Query(stroka)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare find recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		description := ""
		products := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag, &model_recipe.Video)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}

	return array, nil

}

func (r *SqlRepo) RecipeByCategory(recipe model.Recipe) ([]model.Recipe, error) {
	array := []model.Recipe{}
	var stroka string
	if recipe.Name != "" {
		stroka = `SELECT * FROM recipes 
		WHERE name LIKE '` + recipe.Name + `%'
		AND category = '` + recipe.Category + `'  
		ORDER BY name;`
	} else {
		stroka = `SELECT * FROM recipes
		WHERE category = '` + recipe.Category + `'  
		ORDER BY name;`
	}
	fmt.Println(stroka)
	rows, err := r.DB.Query(stroka)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare find recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		description := ""
		products := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag, &model_recipe.Video)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}

	return array, nil
}

func (r *SqlRepo) RecipeByTag(recipe model.Recipe) ([]model.Recipe, error) {
	array := []model.Recipe{}
	var stroka string
	if recipe.Name != "" {
		stroka = `SELECT * FROM recipes 
		WHERE name LIKE '` + recipe.Name + `%'
		AND tag = '` + recipe.Tag + `'  
		ORDER BY name;`
	} else {
		stroka = `SELECT * FROM recipes
		WHERE tag = '` + recipe.Tag + `'  
		ORDER BY name;`
	}
	fmt.Println(stroka)
	rows, err := r.DB.Query(stroka)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare find recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		description := ""
		products := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag, &model_recipe.Video)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}

	return array, nil
}

func (r *SqlRepo) RecipeById(recipe_id uint) (model.Recipe, error) {
	model_recipe := model.Recipe{}
	rows, err := r.DB.Query(`SELECT * FROM `+"`recipes`"+` WHERE recipe_id = ?`, recipe_id)
	if err != nil {
		return model_recipe, fmt.Errorf("failed to prepare find user: %w", err)
	}

	for rows.Next() {
		var description string
		var products string
		err := rows.Scan(
			&model_recipe.RecipeID,
			&model_recipe.UserID,
			&model_recipe.Name,
			&model_recipe.Image,
			&description,
			&products,
			&model_recipe.Category,
			&model_recipe.Tag,
			&model_recipe.Video,
		)
		if err != nil {
			return model_recipe, fmt.Errorf("failed to scan row: %w", err)
		}

		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
	}
	return model_recipe, nil
}

func (r *SqlRepo) RecipeByUserId(user_id uint) ([]model.Recipe, error) {
	array := []model.Recipe{}

	rows, err := r.DB.Query(`SELECT * FROM `+"`recipes`"+` WHERE user_id = ? ORDER BY name`, user_id)
	if err != nil {
		return []model.Recipe{}, fmt.Errorf("failed to prepare FindAll recipe: %w", err)
	}

	for rows.Next() {
		model_recipe := model.Recipe{}
		products := ""
		description := ""
		rows.Scan(&model_recipe.RecipeID, &model_recipe.UserID, &model_recipe.Name, &model_recipe.Image, &description, &products, &model_recipe.Category, &model_recipe.Tag, &model_recipe.Video)
		model_recipe.UnmarshalSteps(description)
		model_recipe.UnmarshalProducts(products)
		array = append(array, model_recipe)
	}
	return array, nil
}
