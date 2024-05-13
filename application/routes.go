package application

import (
	"cookvs/handler"
	"cookvs/repository/users"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Route("/users", a.loadUserRoutes)
	router.Route("/recipes", a.loadRecipeRoutes)
	router.Handle("/assets/*", http.FileServer(http.Dir(".")))
	router.Handle("/assets/image", http.FileServer(http.Dir(".")))

	a.router = router
}

func (a *App) loadUserRoutes(router chi.Router) {
	userHandler := &handler.Cook{
		Repo: &users.SqlRepo{
			DB: a.db,
		},
	}

	router.Post("/login", userHandler.FindByEmail)
	router.Post("/check", userHandler.CheckEmail)
	router.Post("/", userHandler.Create)
	router.Get("/", userHandler.List)
	router.Put("/{id}", userHandler.UpdateByID)
	router.Delete("/{id}", userHandler.DeleteByID)
	router.Post("/uploaduser", userHandler.UploadImageUser)
}

func (a *App) loadRecipeRoutes(router chi.Router) {
	recipeHandler := &handler.Cook{
		Repo: &users.SqlRepo{
			DB: a.db,
		},
	}

	router.Get("/", recipeHandler.ListRecipe)
	router.Post("/", recipeHandler.CreateRecipe)
	router.Post("/uploadrecipe", recipeHandler.UploadImageRecipes)
	router.Post("/NameFind", recipeHandler.FindByName)
	router.Post("/Category", recipeHandler.FindByCategory)
	router.Post("/Tag", recipeHandler.FindByTag)
	router.Get("/{id}", recipeHandler.RecipeByID)
	router.Post("/user{id}", recipeHandler.ListRecipeByUser)
	router.Post("/category", recipeHandler.ListRecipeByCategory)
	router.Post("/tag", recipeHandler.ListRecipeByTag)
	router.Post("/addcomments", recipeHandler.AddComments)
	router.Post("/comments{id}", recipeHandler.CommentsList)
	router.Post("/popularrecipe", recipeHandler.RecipeByCountComments)
}
