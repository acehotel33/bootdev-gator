package commands

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/acehotel33/bootdev-gator/internal/database"
	"github.com/acehotel33/bootdev-gator/internal/rss"
	"github.com/acehotel33/bootdev-gator/internal/state"
	"github.com/google/uuid"
)

type Command struct {
	name      string
	arguments []string
}

type Commands struct {
	commandsMap map[string]func(*state.State, Command) error
}

func InitializeCommands() (*Commands, error) {
	cmds := &Commands{
		commandsMap: map[string]func(*state.State, Command) error{},
	}
	cmds.Register("login", HandlerLogin)
	cmds.Register("register", HandlerRegister)
	cmds.Register("reset", HandlerReset)
	cmds.Register("users", HandlerGetAllUsers)
	cmds.Register("agg", middlewareLoggedIn(HandlerAggregator))
	cmds.Register("addfeed", middlewareLoggedIn(HandlerAddFeed))
	cmds.Register("feeds", HandlerFeeds)
	cmds.Register("follow", middlewareLoggedIn(HandlerFollow))
	cmds.Register("following", middlewareLoggedIn(HandlerFollowing))
	cmds.Register("unfollow", middlewareLoggedIn(HandlerUnfollow))
	cmds.Register("browse", middlewareLoggedIn(HandlerBrowse))
	return cmds, nil
}

func RunCommand(state *state.State, cmds *Commands) error {
	if len(os.Args) < 2 {
		return errors.New("not enough arguments")
	}

	commandName := os.Args[1]
	commandArgs := os.Args[2:]

	cmd := Command{
		name:      commandName,
		arguments: commandArgs,
	}

	if err := cmds.Run(state, cmd); err != nil {
		return err
	}

	return nil
}

func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.commandsMap[name] = f
}

func (c *Commands) Run(s *state.State, cmd Command) error {
	if f, ok := c.commandsMap[cmd.name]; ok {
		if err := f(s, cmd); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("invalid command")
}

func HandlerRegister(s *state.State, cmd Command) error {
	if len(cmd.arguments) != 1 {
		return errors.New("invalid arguments")
	}

	username := cmd.arguments[0]
	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}

	user, err := s.DB.CreateUser(context.Background(), createUserParams)
	if err != nil {
		return errors.New("could not register user")
	}

	s.Cfg.SetUser(username)
	fmt.Printf("user %v created\n", username)
	fmt.Println(user)

	return nil
}

func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.arguments) == 0 || len(cmd.arguments) > 1 {
		return errors.New("invalid arguments")
	}

	username := cmd.arguments[0]
	_, err := s.DB.GetUser(context.Background(), username)
	if err != nil {
		return errors.New("could not log in")
	}
	if err := s.Cfg.SetUser(username); err != nil {
		return err
	}
	fmt.Printf("user has been set to: %v\n", username)
	return nil
}

func HandlerReset(s *state.State, cmd Command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("invalid arguments")
	}

	if err := s.DB.ResetUsers(context.Background()); err != nil {
		return err
	}

	fmt.Println("users reset")
	return nil
}

func HandlerGetAllUsers(s *state.State, cmd Command) error {
	users, err := s.DB.GetAllUsers(context.Background())
	if err != nil {
		return errors.New("could not get users")
	}

	for _, user := range users {
		line := "* " + user.Name
		if user.Name == s.Cfg.CurrentUsername {
			line = line + " (current)"
		}
		fmt.Println(line)
	}

	return nil
}

func HandlerAggregator(s *state.State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("invalid arguments")
	}

	timeBetweenReqs := cmd.arguments[0]
	timeDuration, err := time.ParseDuration(timeBetweenReqs)
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeDuration)

	// rssFeed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	return err
	// }
	//
	// fmt.Println(rssFeed)
	// return nil

	ticker := time.NewTicker(timeDuration)

	// fmt.Println("Ticker created")

	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}

}

func HandlerAddFeed(s *state.State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 2 {
		return errors.New("invalid arguments")
	}
	feedName := cmd.arguments[0]
	feedURL := cmd.arguments[1]

	currentUserID := user.ID

	_, err := rss.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	createFeedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID: uuid.NullUUID{
			UUID:  currentUserID,
			Valid: true,
		},
	}

	dbFeed, err := s.DB.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return err
	}

	createFeedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: uuid.NullUUID{
			UUID:  currentUserID,
			Valid: true,
		},
		FeedID: uuid.NullUUID{
			UUID:  dbFeed.ID,
			Valid: true,
		},
	}
	_, err = s.DB.CreateFeedFollow(context.Background(), createFeedFollowParams)
	if err != nil {
		return err
	}

	fmt.Println(dbFeed)
	return nil
}

func HandlerFeeds(s *state.State, cmd Command) error {
	if len(cmd.arguments) != 0 {
		return errors.New("invalid arguments")
	}

	feed, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	type returnParams struct {
		Name     string
		Url      string
		Username string
	}

	for _, item := range feed {
		dbUser, err := s.DB.GetUserByID(context.Background(), item.UserID.UUID)
		if err != nil {
			return err
		}
		userName := dbUser.Name

		params := returnParams{
			Name:     item.Name,
			Url:      item.Url,
			Username: userName,
		}

		fmt.Println(params)
	}

	return nil
}

func HandlerFollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("invalid arguments")
	}

	userDB := user
	feedURL := cmd.arguments[0]
	feedDB, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: uuid.NullUUID{
			UUID:  userDB.ID,
			Valid: true,
		},
		FeedID: uuid.NullUUID{
			UUID:  feedDB.ID,
			Valid: true,
		},
	}

	feedFollowDB, err := s.DB.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	type returnStruct struct {
		feedName string
		username string
	}

	returnParams := returnStruct{
		feedName: feedFollowDB.FeedName,
		username: feedFollowDB.UserName,
	}

	fmt.Println(returnParams)
	return nil
}

func HandlerFollowing(s *state.State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 0 {
		return errors.New("invalid arguments")
	}

	userDB := user

	following, err := s.DB.GetFeedFollowsForUser(context.Background(), uuid.NullUUID{
		UUID:  userDB.ID,
		Valid: true,
	})
	if err != nil {
		return err
	}

	for i := range following {
		fmt.Println(following[i].FeedName)
	}

	return nil
}

func HandlerUnfollow(s *state.State, cmd Command, user database.User) error {
	if len(cmd.arguments) != 1 {
		return errors.New("invalid arguments")
	}

	feedURL := cmd.arguments[0]
	feedDB, err := s.DB.GetFeedByUrl(context.Background(), feedURL)
	if err != nil {
		return err
	}

	deleteFeedFollowParams := database.DeleteFeedFollowParams{
		UserID: uuid.NullUUID{
			UUID:  user.ID,
			Valid: true,
		},
		FeedID: uuid.NullUUID{
			UUID:  feedDB.ID,
			Valid: true,
		},
	}
	_, err = s.DB.DeleteFeedFollow(context.Background(), deleteFeedFollowParams)
	if err != nil {
		return err
	}

	return nil
}

func scrapeFeeds(s *state.State, user database.User) error {
	nextFeed, err := s.DB.GetNextFeedToFetch(context.Background(), uuid.NullUUID{UUID: user.ID, Valid: true})
	if err != nil {
		return err
	}

	markFeedFetchedParams := database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	markedFeed, err := s.DB.MarkFeedFetched(context.Background(), markFeedFetchedParams)
	if err != nil {
		return err
	}

	fetchedFeed, err := rss.FetchFeed(context.Background(), markedFeed.Url)
	if err != nil {
		return err
	}

	fmt.Println("----------")
	fmt.Printf("Fetching feed: %v - %v\n", markedFeed.Name, markedFeed.Url)
	fmt.Println("----------")

	fetchedItems := fetchedFeed.Channel.Item
	for _, item := range fetchedItems {

		postDB, err := s.DB.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: timeStringParser(item.PubDate), Valid: true},
			FeedID:      uuid.NullUUID{UUID: markedFeed.ID, Valid: true},
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint \"posts_url_key\"") {
				// fmt.Printf("ignoring existing post: %v\n", item.Title)
				continue
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("%v added to posts DB\n", postDB.Title)
		}

	}

	return nil
}

func timeStringParser(s string) time.Time {

	formats := []string{
		time.RFC1123,
		time.RFC822,
		"2006-01-02T15:04:05-07:00"}

	var parsedTime time.Time
	var err error
	for _, format := range formats {
		parsedTime, err = time.Parse(format, s)
		if err == nil {
			break
		}
	}
	if err != nil {
		parsedTime = time.Time{}
	}

	return parsedTime
}

func HandlerBrowse(s *state.State, cmd Command, user database.User) error {
	if len(cmd.arguments) > 1 {
		return errors.New("invalid arguments")
	}

	var limit int32
	limit = 2
	if len(cmd.arguments) == 1 {
		limitInt, err := strconv.Atoi(cmd.arguments[0])
		limit = int32(limitInt)
		if err != nil {
			limit = 2
		}

	}
	postsDB, err := s.DB.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: uuid.NullUUID{UUID: user.ID, Valid: true}, Limit: limit})
	if err != nil {
		return err
	}

	for _, post := range postsDB {
		fmt.Println("---------")
		fmt.Println(post.Title.String)
		fmt.Println(post.Description.String)
		fmt.Println("---------")
		fmt.Println()
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state.State, cmd Command, user database.User) error) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Cfg.CurrentUsername)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
