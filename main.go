package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyService struct {
	client *spotify.Client
}

func NewSpotifyService() (*SpotifyService, error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET must be set")
	}

	ctx := context.Background()
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get token: %v", err)
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	return &SpotifyService{client: client}, nil
}

func (s *SpotifyService) SearchTracks(query string, limit int) ([]spotify.FullTrack, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.client.Search(ctx, query, spotify.SearchTypeTrack, spotify.Limit(limit))
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	if result.Tracks == nil {
		return nil, fmt.Errorf("no tracks found")
	}

	return result.Tracks.Tracks, nil
}

func (s *SpotifyService) SearchArtists(query string, limit int) ([]spotify.FullArtist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.client.Search(ctx, query, spotify.SearchTypeArtist, spotify.Limit(limit))
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	if result.Artists == nil {
		return nil, fmt.Errorf("no artists found")
	}

	return result.Artists.Artists, nil
}

func PrintTracks(tracks []spotify.FullTrack) {
	fmt.Println("\nНайденные треки:")
	fmt.Println("========================================")
	for i, track := range tracks {
		duration := track.TimeDuration()
		fmt.Printf("%d. %s - %s\n", i+1, track.Name, joinArtists(track.Artists))
		fmt.Printf("   Альбом: %s\n", track.Album.Name)
		fmt.Printf("   Длительность: %d:%02d\n", int(duration.Minutes()), int(duration.Seconds())%60)
		fmt.Printf("   ID: %s\n", track.ID)
		fmt.Printf("   Популярность: %d/100\n", track.Popularity)
		fmt.Println("----------------------------------------")
	}
}

func PrintArtists(artists []spotify.FullArtist) {
	fmt.Println("\nНайденные артисты:")
	fmt.Println("========================================")
	for i, artist := range artists {
		fmt.Printf("%d. %s\n", i+1, artist.Name)
		fmt.Printf("   Жанры: %s\n", strings.Join(artist.Genres, ", "))
		fmt.Printf("   Популярность: %d/100\n", artist.Popularity)
		fmt.Printf("   ID: %s\n", artist.ID)
		fmt.Println("----------------------------------------")
	}
}

func joinArtists(artists []spotify.SimpleArtist) string {
	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name
	}
	return strings.Join(names, ", ")
}

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	service, err := NewSpotifyService()
	if err != nil {
		log.Fatalf("Ошибка создания сервиса: %v", err)
	}

	tracks, err := service.SearchTracks("надо было ставить линукс", 1)
	if err != nil {
		log.Printf("Ошибка поиска треков: %v", err)
	} else {
		PrintTracks(tracks)
	}

	artists, err := service.SearchArtists("cupsize", 3)
	if err != nil {
		log.Printf("Ошибка поиска артистов: %v", err)
	} else {
		PrintArtists(artists)
	}
}
