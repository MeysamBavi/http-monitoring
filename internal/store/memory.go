package store

import (
	"context"

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

type InMemoryUser struct {
	nextId    model.ID
	data      map[model.ID]*model.User
	usernames map[string]model.ID
}

func (u *InMemoryUser) newId() model.ID {
	u.nextId++
	return u.nextId
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
	nextId model.ID
	data   map[model.ID][]*model.URL // user id -> urls
}

func (u *InMemoryUrl) newId() model.ID {
	u.nextId++
	return u.nextId
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

func (u *InMemoryUrl) UpdateStat(_ context.Context, userId model.ID, id model.ID, stat model.DayStat) (*model.URL, error) {

	urls, ok := u.data[userId]
	if !ok {
		return nil, NewNotFoundError("url", "userId", userId)
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
				return url, nil
			}
		}
		// if no day stat was found, add the passed day stat
		url.DayStats = append(url.DayStats, &stat)
		return url, nil
	}

	return nil, NewNotFoundError("url", "id", id)
}

type InMemoryAlert struct {
	nextId model.ID
	data   map[model.ID][]*model.Alert // url id -> alerts
}

func (a *InMemoryAlert) newId() model.ID {
	a.nextId++
	return a.nextId
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
