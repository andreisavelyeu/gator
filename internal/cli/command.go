package cli

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"gator/internal/database"
	"gator/internal/state"
	"gator/internal/utils"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Registered map[string]func(*state.State, Command) error
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	_, ok := c.Registered[name]

	if ok {
		fmt.Println("command has already been registered")
	} else {
		c.Registered[name] = f
	}
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	command, ok := c.Registered[cmd.Name]

	if ok {
		err := command(s, cmd)
		return err
	} else {
		return errors.New("command not found")
	}
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	username := cmd.Args[0]

	user, err := s.Db.GetUser(context.Background(), username)

	if err != nil {
		return err
	}

	s.Config.SetUser(user.Name)

	fmt.Printf("%s has been set", user.Name)
	return nil
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	user := cmd.Args[0]

	newUser := database.CreateUserParams{
		Name:      user,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
	}

	createdUser, err := s.Db.CreateUser(context.Background(), newUser)

	if err != nil {
		return err
	}

	s.Config.SetUser(createdUser.Name)

	fmt.Printf("%s has been created, id: %v", createdUser.Name, createdUser.ID)
	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("users have been deleted")
	return nil
}

func HandlerGetUsers(s *state.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.Config.Current_user_name {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil

}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Add("User-Agent", "gator")

	client := http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return &RSSFeed{}, err
	}

	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return &RSSFeed{}, err
	}

	var feed *RSSFeed

	err = xml.Unmarshal(body, &feed)

	if err != nil {
		return &RSSFeed{}, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for k := range feed.Channel.Item {
		feed.Channel.Item[k].Title = html.UnescapeString(feed.Channel.Item[k].Title)
		feed.Channel.Item[k].Description = html.UnescapeString(feed.Channel.Item[k].Description)
	}
	return feed, nil
}

func HandlerAgg(s *state.State, cmd Command) error {

	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])

	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)

		if err != nil {
			return err
		}

	}

}

func HandlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("not enough arguments")
	}

	url := cmd.Args[1]
	name := cmd.Args[0]

	newFeed := database.CreateFeedParams{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
		UserID:    user.ID,
		Url:       url,
	}

	feed, err := s.Db.CreateFeed(context.Background(), newFeed)

	if err != nil {
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	newRecord, err := s.Db.CreateFeedFollow(context.Background(), feedFollow)

	if err != nil {
		return err
	}

	fmt.Println(newRecord)
	fmt.Println(feed)
	return nil
}

func HandlerGetFeeds(s *state.State, cmd Command) error {
	allFeeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for k, v := range allFeeds {
		fmt.Printf("Feed %v\n", k)
		fmt.Printf("Name: %s\n", v.Name)
		fmt.Printf("Url: %s\n", v.Url)
		fmt.Printf("Created by: %s\n", v.UserName)
	}
	return nil
}

func HandlerFollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	url := cmd.Args[0]

	feed, err := s.Db.GetFeed(context.Background(), url)

	if err != nil {
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	newRecord, err := s.Db.CreateFeedFollow(context.Background(), feedFollow)

	if err != nil {
		return err
	}

	fmt.Println(newRecord)
	return nil
}

func HandlerFollowing(s *state.State, cmd Command, user database.User) error {
	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return err
	}

	for _, v := range feeds {
		fmt.Println(v.Name)
	}
	return nil
}

func HandlerUnfollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.Args) == 0 {
		return errors.New("not enough arguments")
	}

	url := cmd.Args[0]

	deleteFeedParams := database.DeleteFeedFollowByUserAndUrlParams{
		UserID: user.ID,
		Url:    url,
	}
	err := s.Db.DeleteFeedFollowByUserAndUrl(context.Background(), deleteFeedParams)

	if err != nil {
		return err
	}
	return nil
}

func scrapeFeeds(s *state.State) error {
	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)

	if err != nil {
		return err
	}

	params := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		ID: nextFeed.ID,
	}

	err = s.Db.MarkFeedFetched(context.Background(), params)

	if err != nil {
		return err
	}

	for _, item := range feed.Channel.Item {
		publishedAt, _ := utils.ParsePublishedAt(item.PubDate)
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			PublishedAt: publishedAt,
			Url:         item.Link,
			Title:       item.Title,
			Description: sql.NullString{String: item.Description, Valid: true},
			FeedID:      nextFeed.ID,
		}
		s.Db.CreatePost(context.Background(), params)
	}

	return nil
}

func HandleBrowse(s *state.State, cmd Command) error {
	var limitString string

	if len(cmd.Args) == 0 {
		limitString = "2"
	} else {
		limitString = cmd.Args[0]
	}

	limit, err := strconv.Atoi(limitString)
	if err != nil {
		limit = 2
	}

	posts, err := s.Db.GetPosts(context.Background(), int32(limit))

	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Println(post.Title)
	}

	return nil
}
