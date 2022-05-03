package web

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"playlisttogether/backend/config"
	"playlisttogether/backend/memories"
	"playlisttogether/backend/playlist"
	"playlisttogether/backend/spotify"
	"time"

	"github.com/gorilla/mux"
)

func Serve(conf *config.Config, db *sql.DB) {

	headerActive := conf.Headers
	//creating Spotify Instance
	spotifyInstance := spotify.NewSpotifyInstance(&headerActive)

	//creating Playlist Instance
	playlistInstance := playlist.NewPlaylistInstance([]byte(conf.Secret), &headerActive, conf.ExpireTime)

	r := mux.NewRouter()
	r.HandleFunc("/search", spotifyInstance.Searcher).Methods(http.MethodGet)
	r.HandleFunc("/login", playlistInstance.Login(db))

	//handlers with middleware
	tvr := r.PathPrefix("/playlist").Subrouter()
	tvr.Use(playlistInstance.JwtMiddleware(db))
	tvr.HandleFunc("/logout", playlistInstance.Logout(db))
	//playlist
	tvr.HandleFunc("/addsongtobucket", playlist.AddSongToBucket(db))
	tvr.HandleFunc("/getbucket", playlistInstance.GetBucket(db))
	tvr.HandleFunc("/getplaylist", playlistInstance.GetPlaylist(db))
	tvr.HandleFunc("/delete", playlistInstance.Delete(db))
	tvr.HandleFunc("/dislike", playlistInstance.Dislike(db))
	//playlistuser
	tvr.HandleFunc("/getusers", playlistInstance.GetUsers(db))
	//playlists
	tvr.HandleFunc("/getplaylists", playlistInstance.GetPlaylists(db))
	tvr.HandleFunc("/createplaylist", playlist.CreatePlaylist(db))
	tvr.HandleFunc("/getsonguris", playlistInstance.GetSongUris(db))
	tvr.HandleFunc("/saveplaylist", playlistInstance.UpdatePlaylist(db))
	//memories
	tvr.HandleFunc("/getmemories", memories.GetMemories(db, conf.ImagePath))

	//TODO
	// rename Playlist endpoint
	// delete Playlist endpoint
	// leave Playlist endpoint
	// maybe kick User endpoint

	//not used
	//tvr.HandleFunc("/getamount", playlistInstance.GetAmount(db))
	//tvr.HandleFunc("/merge", playlistInstance.Merge(db))

	//the duration to shut down server
	wait := time.Second * 15

	srv := &http.Server{
		Addr: conf.Host + ":" + conf.HostPort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	//log.Fatal(http.ListenAndServe(conf.Host+":"+conf.HostPort, r))
}
