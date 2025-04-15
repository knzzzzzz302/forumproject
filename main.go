package main

import (
	"FORUM-GO/databaseAPI"
	"FORUM-GO/webAPI"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
)

type Post struct {
	Id         int
	Username   string
	Title      string
	Categories []string
	Content    string
	CreatedAt  string
	UpVotes    int
	DownVotes  int
	Comments   []Comment
}

type Comment struct {
	Id        int
	PostId    int
	Username  string
	Content   string
	CreatedAt string
}

// Database
var database *sql.DB

func main() {
	// check if DB exists
	var _, err = os.Stat("database.db")

	// create DB if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create("database.db")
		if err != nil {
			return
		}
		defer file.Close()
	}

	database, _ = sql.Open("sqlite3", "./database.db")

	databaseAPI.CreateUsersTable(database)
	databaseAPI.CreatePostTable(database)
	databaseAPI.CreateCommentTable(database)
	databaseAPI.CreateVoteTable(database)
	databaseAPI.CreateCategoriesTable(database)
	databaseAPI.CreateCategories(database)
	databaseAPI.CreateCategoriesIcons(database)
	databaseAPI.CreateCommentLikesTable(database)
	databaseAPI.CreateCommentDislikesTable(database)
	databaseAPI.CreatePostImagesTable(database)
	
	// Créer le dossier pour stocker les images des posts
	os.MkdirAll("public/uploads/posts", os.ModePerm)

	webAPI.SetDatabase(database)

	fs := http.FileServer(http.Dir("public"))
	router := http.NewServeMux()
	fmt.Println("Starting server on port http://localhost:3030/")

	router.HandleFunc("/", webAPI.Index)
	router.HandleFunc("/register", webAPI.Register)
	router.HandleFunc("/login", webAPI.Login)
	router.HandleFunc("/post", webAPI.DisplayPost)
	router.HandleFunc("/filter", webAPI.GetPostsByApi)
	router.HandleFunc("/newpost", webAPI.NewPost)
	router.HandleFunc("/api/register", webAPI.RegisterApi)
	router.HandleFunc("/api/login", webAPI.LoginApi)
	router.HandleFunc("/api/logout", webAPI.LogoutAPI)
	router.HandleFunc("/api/createpost", webAPI.CreatePostApi)
	router.HandleFunc("/api/comments", webAPI.CommentsApi)
	router.HandleFunc("/api/vote", webAPI.VoteApi)
	router.HandleFunc("/api/deletepost", webAPI.DeletePostHandler)
    router.HandleFunc("/profile", webAPI.DisplayProfile)
    router.HandleFunc("/api/editprofile", webAPI.EditProfileHandler)
    router.HandleFunc("/api/changepassword", webAPI.ChangePasswordHandler)
    router.HandleFunc("/api/uploadprofileimage", webAPI.UploadProfileImageHandler)
	router.HandleFunc("/api/editcomment", webAPI.EditCommentHandler)
	router.Handle("/public/", http.StripPrefix("/public/", fs))
	router.HandleFunc("/editpost", webAPI.EditPostPage)
	router.HandleFunc("/api/editpost", webAPI.EditPostHandler)
	router.HandleFunc("/api/deletecomment", webAPI.DeleteCommentHandler)
	router.HandleFunc("/api/commentlike", webAPI.CommentLikeApi)

	http.ListenAndServe(":3030", router)
}