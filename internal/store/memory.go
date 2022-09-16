package store

import (
	"context"
	"fmt"
	"time"

	"github.com/MeysamBavi/http-monitoring/internal/model"
	"go.uber.org/zap"
)

type InMemoryStore struct {
	user   *InMemoryUser
	url    *InMemoryUrl
	alert  *InMemoryAlert
	logger *zap.Logger
}

func NewInMemoryStore(logger *zap.Logger) Store {
	return &InMemoryStore{
		user:   &InMemoryUser{data: make(map[model.ID]*model.User), usernames: make(map[string]model.ID)},
		url:    &InMemoryUrl{data: make(map[model.ID][]*model.URL)},
		alert:  &InMemoryAlert{data: make(map[model.ID][]*model.Alert)},
		logger: logger,
	}
}

func (s *InMemoryStore) User() User {
	return s.user
}

func (s *InMemoryStore) Url() Url {
	return s.url
}

func (s *InMemoryStore) Alert() Alert {
	return s.alert
}

type idGen int

func (ign *idGen) newId() model.ID {
	(*ign)++
	r, _ := model.ParseId(fmt.Sprint((*ign)))
	return r
}

type InMemoryUser struct {
	idGen
	data      map[model.ID]*model.User
	usernames map[string]model.ID
}

func (u *InMemoryUser) Get(_ context.Context, id model.ID) (*model.User, error) {
	user, ok := u.data[id]
	if !ok {
		return nil, NewNotFoundError("user", "id", id)
	}

	return user, nil
}

func (u *InMemoryUser) GetByUsername(_ context.Context, username string) (*model.User, error) {
	userId, ok := u.usernames[username]
	if !ok {
		return nil, NewNotFoundError("user", "username", username)
	}

	return u.data[userId], nil
}

func (u *InMemoryUser) Add(_ context.Context, user *model.User) error {
	if _, ok := u.usernames[user.Username]; ok {
		return NewDuplicateError("user", "username", user.Username)
	}

	user.Id = u.newId()

	u.data[user.Id] = user
	u.usernames[user.Username] = user.Id

	return nil
}

type InMemoryUrl struct {
	idGen
	data map[model.ID][]*model.URL // user id -> urls
}

func (u *InMemoryUrl) GetByUserId(_ context.Context, id model.ID) ([]*model.URL, error) {
	urls := u.data[id]
	return urls, nil
}

func (u *InMemoryUrl) Add(_ context.Context, url *model.URL) error {
	url.Id = u.newId()

	urls := u.data[url.UserId]
	u.data[url.UserId] = append(urls, url)

	return nil
}

func (u *InMemoryUrl) GetDayStat(_ context.Context, userId model.ID, id model.ID, date model.Date) (model.DayStat, error) {

	urls, ok := u.data[userId]
	if !ok {
		return model.DayStat{}, NewNotFoundError("url", "userId", userId)
	}

	// find url among user urls
	for _, url := range urls {
		if url.Id != id {
			continue
		}

		// find day stat among url day stats
		for _, ds := range url.DayStats {
			if ds.Date == date {
				return *ds, nil
			}
		}

		return model.DayStat{}, NewNotFoundError("stat", "date", date)
	}

	return model.DayStat{}, NewNotFoundError("url", "id", id)
}

func (u *InMemoryUrl) UpdateStat(_ context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, model.DayStat, error) {

	urls, ok := u.data[userId]
	if !ok {
		return nil, model.DayStat{}, NewNotFoundError("url", "userId", userId)
	}

	// find url among user urls
	for _, url := range urls {
		if url.Id != id {
			continue
		}

		// find day stat among url day stats
		for _, ds := range url.DayStats {
			// apply change
			if ds.Date == stat.Date {
				ds.FailureCount += stat.FailureCount
				ds.SuccessCount += stat.SuccessCount
				return url, *ds, nil
			}
		}
		// if no day stat was found, add the passed day stat
		url.DayStats = append(url.DayStats, &stat)
		return url, stat, nil
	}

	return nil, model.DayStat{}, NewNotFoundError("url", "id", id)
}

func (u *InMemoryUrl) ForAll(_ context.Context, callBack func(model.URL)) error {
	for _, urls := range u.data {
		for _, url := range urls {
			callBack(*url)
		}
	}

	return nil
}

// for test only
func (u *InMemoryUrl) ListenForChanges(_ context.Context, out chan<- UrlChangeEvent) error {
	i, _ := model.ParseId("1")
	for {
		time.Sleep(10 * time.Second)
		out <- UrlChangeEvent{
			Url: model.URL{
				Id:        i,
				UserId:    i,
				Url:       "https://httpbin.org/status/206",
				Threshold: 20,
				Interval:  model.Interval{Duration: 30 * time.Second},
			},
			Operation: UrlChangeOperationInsert,
		}
	}
}

type InMemoryAlert struct {
	idGen
	data map[model.ID][]*model.Alert // url id -> alerts
}

func (a *InMemoryAlert) GetByUrlId(_ context.Context, urlId model.ID) ([]*model.Alert, error) {
	alerts := a.data[urlId]
	return alerts, nil
}

func (a *InMemoryAlert) Add(_ context.Context, alert *model.Alert) error {
	alert.Id = a.newId()

	alerts := a.data[alert.UrlId]
	a.data[alert.UrlId] = append(alerts, alert)

	return nil
}
