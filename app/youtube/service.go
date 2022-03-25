// Package youtube provides loading audio from video files for given youtube channels
package youtube

import (
	"context"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/umputun/feed-master/app/feed"
	"github.com/umputun/feed-master/app/youtube/channel"
)

//go:generate moq -out mocks/downloader.go -pkg mocks -skip-ensure -fmt goimports . DownloaderService
//go:generate moq -out mocks/channel.go -pkg mocks -skip-ensure -fmt goimports . ChannelService
//go:generate moq -out mocks/store.go -pkg mocks -skip-ensure -fmt goimports . StoreService

// Service loads audio from youtube channels
type Service struct {
	Channels       []ChannelInfo
	Downloader     DownloaderService
	ChannelService ChannelService
	Store          StoreService
	CheckDuration  time.Duration
	RSSFileStore   RSSFileStore
	KeepPerChannel int
	RootURL        string
	processed      map[string]bool
}

// ChannelInfo is a pait of channel ID and name
type ChannelInfo struct {
	Name string
	ID   string
}

// DownloaderService is an interface for downloading audio from youtube
type DownloaderService interface {
	Get(ctx context.Context, id string, fname string) (file string, err error)
}

// ChannelService is an interface for getting channel entries, i.e. the list of videos
type ChannelService interface {
	Get(ctx context.Context, chanID string) ([]channel.Entry, error)
}

// StoreService is an interface for storing and loading metadata about downloaded audio
type StoreService interface {
	Save(entry channel.Entry) (bool, error)
	Load(channelID string, max int) ([]channel.Entry, error)
	Exist(entry channel.Entry) (bool, error)
	RemoveOld(channelID string, keep int) ([]string, error)
}

// Do is a blocking function that downloads audio from youtube channels and updates metadata
func (s *Service) Do(ctx context.Context) error {
	log.Printf("[INFO] Starting youtube service")

	s.processed = make(map[string]bool)
	tick := time.NewTicker(s.CheckDuration)
	defer tick.Stop()

	if err := s.procChannels(ctx); err != nil {
		return errors.Wrap(err, "failed to process channels")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			if err := s.procChannels(ctx); err != nil {
				return errors.Wrap(err, "failed to process channels")
			}
		}
	}
}

// RSSFeed generates RSS feed for given channel
func (s *Service) RSSFeed(cinfo ChannelInfo) (string, error) {
	entries, err := s.Store.Load(cinfo.ID, s.KeepPerChannel)
	if err != nil {
		return "", errors.Wrap(err, "failed to get channel entries")
	}

	if len(entries) == 0 {
		return "", nil
	}

	items := []feed.Item{}
	for _, entry := range entries {

		fileURL := s.RootURL + "/" + path.Base(entry.File)

		var fileSize int
		if fileInfo, fiErr := os.Stat(entry.File); fiErr != nil {
			log.Printf("[WARN] failed to get file size for %s: %v", entry.File, fiErr)
		} else {
			fileSize = int(fileInfo.Size())
		}

		items = append(items, feed.Item{
			Title:       entry.Title,
			Description: entry.Media.Description,
			Link:        entry.Link.Href,
			PubDate:     entry.Published.Format(time.RFC822Z),
			GUID:        entry.ChannelID + "::" + entry.VideoID,
			Author:      entry.Author.Name,
			Enclosure: feed.Enclosure{
				URL:    fileURL,
				Type:   "audio/mpeg",
				Length: fileSize,
			},
			DT: time.Now(),
		})
	}

	rss := feed.Rss2{
		Version:       "2.0",
		ItemList:      items,
		Title:         entries[0].Author.Name,
		Description:   "generated by feed-master",
		Link:          entries[0].Author.URI,
		PubDate:       items[0].PubDate,
		LastBuildDate: time.Now().Format(time.RFC822Z),
	}

	b, err := xml.MarshalIndent(&rss, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal rss")
	}

	return string(b), nil
}

func (s *Service) procChannels(ctx context.Context) error {
	for _, chanInfo := range s.Channels {
		entries, err := s.ChannelService.Get(ctx, chanInfo.ID)
		if err != nil {
			log.Printf("[WARN] failed to get channel entries for %s: %s", chanInfo.ID, err)
			continue
		}
		log.Printf("[INFO] got %d entries for %s, limit to %d", len(entries), chanInfo.Name, s.KeepPerChannel)
		changed := false
		for i, entry := range entries {
			if i >= s.KeepPerChannel {
				break
			}

			// check if we already processed this entry.
			// this is needed to avoid infinite get/remove loop when the original feed is updated in place
			if _, ok := s.processed[entry.UID()]; ok {
				log.Printf("[DEBUG] skipping already processed entry %s", entry.VideoID)
				continue
			}

			// check if entry already exists in store
			exists, exErr := s.Store.Exist(entry)
			if err != nil {
				return errors.Wrapf(exErr, "failed to check if entry %s exists", entry.VideoID)
			}
			if exists {
				continue
			}

			log.Printf("[INFO] new entry [%d] %s, %s, %s", i+1, entry.VideoID, entry.Title, chanInfo.Name)
			file, downErr := s.Downloader.Get(ctx, entry.VideoID, s.makeFileName(entry))
			if downErr != nil {
				log.Printf("[WARN] failed to download %s: %s", entry.VideoID, downErr)
				continue
			}
			log.Printf("[DEBUG] downloaded %s (%s) to %s, channel: %+v", entry.VideoID, entry.Title, file, chanInfo)
			entry.File = file
			if !strings.HasPrefix(entry.Title, chanInfo.Name) {
				entry.Title = chanInfo.Name + ": " + entry.Title
			}
			ok, saveErr := s.Store.Save(entry)
			if saveErr != nil {
				return errors.Wrapf(saveErr, "failed to save entry %+v", entry)
			}
			if !ok {
				log.Printf("[WARN] attempt to save dup entry %+v", entry)
			}
			changed = true
			s.processed[entry.UID()] = true // track processed entries
			log.Printf("[INFO] saved %s (%s) to %s, channel: %+v", entry.VideoID, entry.Title, file, chanInfo)
		}

		if changed { // save rss feed to fs if there are new entries
			rss, rssErr := s.RSSFeed(chanInfo)
			if rssErr != nil {
				log.Printf("[WARN] failed to generate rss for %s: %s", chanInfo.Name, rssErr)
			} else {
				if err := s.RSSFileStore.Save(chanInfo.ID, rss); err != nil {
					log.Printf("[WARN] failed to save rss for %s: %s", chanInfo.Name, err)
				}
			}
		}

		// remove old entries and files
		files, rmErr := s.Store.RemoveOld(chanInfo.ID, s.KeepPerChannel+1)
		if rmErr != nil {
			return errors.Wrapf(rmErr, "failed to remove old meta data for %s", chanInfo.ID)
		}
		for _, f := range files {
			if e := os.Remove(f); e != nil {
				log.Printf("[WARN] failed to remove file %s: %s", f, e)
				continue
			}

			log.Printf("[INFO] removed %s for %s (%s)", f, chanInfo.ID, chanInfo.Name)
		}
	}
	log.Printf("[DEBUG] processed channels completed, total %d", len(s.Channels))
	return nil
}

func (s *Service) makeFileName(entry channel.Entry) string {
	h := sha1.New()
	if _, err := h.Write([]byte(entry.ChannelID + "::" + entry.VideoID)); err != nil {
		return uuid.New().String()
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
