// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"sync"

	"github.com/umputun/feed-master/app/youtube"
	ytfeed "github.com/umputun/feed-master/app/youtube/feed"
)

// YoutubeSvcMock is a mock implementation of api.YoutubeSvc.
//
// 	func TestSomethingThatUsesYoutubeSvc(t *testing.T) {
//
// 		// make and configure a mocked api.YoutubeSvc
// 		mockedYoutubeSvc := &YoutubeSvcMock{
// 			RSSFeedFunc: func(cinfo youtube.FeedInfo) (string, error) {
// 				panic("mock out the RSSFeed method")
// 			},
// 			RemoveEntryFunc: func(entry ytfeed.Entry) error {
// 				panic("mock out the RemoveEntry method")
// 			},
// 			StoreRSSFunc: func(chanID string, rss string) error {
// 				panic("mock out the StoreRSS method")
// 			},
// 		}
//
// 		// use mockedYoutubeSvc in code that requires api.YoutubeSvc
// 		// and then make assertions.
//
// 	}
type YoutubeSvcMock struct {
	// RSSFeedFunc mocks the RSSFeed method.
	RSSFeedFunc func(cinfo youtube.FeedInfo) (string, error)

	// RemoveEntryFunc mocks the RemoveEntry method.
	RemoveEntryFunc func(entry ytfeed.Entry) error

	// StoreRSSFunc mocks the StoreRSS method.
	StoreRSSFunc func(chanID string, rss string) error

	// calls tracks calls to the methods.
	calls struct {
		// RSSFeed holds details about calls to the RSSFeed method.
		RSSFeed []struct {
			// Cinfo is the cinfo argument value.
			Cinfo youtube.FeedInfo
		}
		// RemoveEntry holds details about calls to the RemoveEntry method.
		RemoveEntry []struct {
			// Entry is the entry argument value.
			Entry ytfeed.Entry
		}
		// StoreRSS holds details about calls to the StoreRSS method.
		StoreRSS []struct {
			// ChanID is the chanID argument value.
			ChanID string
			// Rss is the rss argument value.
			Rss string
		}
	}
	lockRSSFeed     sync.RWMutex
	lockRemoveEntry sync.RWMutex
	lockStoreRSS    sync.RWMutex
}

// RSSFeed calls RSSFeedFunc.
func (mock *YoutubeSvcMock) RSSFeed(cinfo youtube.FeedInfo) (string, error) {
	if mock.RSSFeedFunc == nil {
		panic("YoutubeSvcMock.RSSFeedFunc: method is nil but YoutubeSvc.RSSFeed was just called")
	}
	callInfo := struct {
		Cinfo youtube.FeedInfo
	}{
		Cinfo: cinfo,
	}
	mock.lockRSSFeed.Lock()
	mock.calls.RSSFeed = append(mock.calls.RSSFeed, callInfo)
	mock.lockRSSFeed.Unlock()
	return mock.RSSFeedFunc(cinfo)
}

// RSSFeedCalls gets all the calls that were made to RSSFeed.
// Check the length with:
//     len(mockedYoutubeSvc.RSSFeedCalls())
func (mock *YoutubeSvcMock) RSSFeedCalls() []struct {
	Cinfo youtube.FeedInfo
} {
	var calls []struct {
		Cinfo youtube.FeedInfo
	}
	mock.lockRSSFeed.RLock()
	calls = mock.calls.RSSFeed
	mock.lockRSSFeed.RUnlock()
	return calls
}

// RemoveEntry calls RemoveEntryFunc.
func (mock *YoutubeSvcMock) RemoveEntry(entry ytfeed.Entry) error {
	if mock.RemoveEntryFunc == nil {
		panic("YoutubeSvcMock.RemoveEntryFunc: method is nil but YoutubeSvc.RemoveEntry was just called")
	}
	callInfo := struct {
		Entry ytfeed.Entry
	}{
		Entry: entry,
	}
	mock.lockRemoveEntry.Lock()
	mock.calls.RemoveEntry = append(mock.calls.RemoveEntry, callInfo)
	mock.lockRemoveEntry.Unlock()
	return mock.RemoveEntryFunc(entry)
}

// RemoveEntryCalls gets all the calls that were made to RemoveEntry.
// Check the length with:
//     len(mockedYoutubeSvc.RemoveEntryCalls())
func (mock *YoutubeSvcMock) RemoveEntryCalls() []struct {
	Entry ytfeed.Entry
} {
	var calls []struct {
		Entry ytfeed.Entry
	}
	mock.lockRemoveEntry.RLock()
	calls = mock.calls.RemoveEntry
	mock.lockRemoveEntry.RUnlock()
	return calls
}

// StoreRSS calls StoreRSSFunc.
func (mock *YoutubeSvcMock) StoreRSS(chanID string, rss string) error {
	if mock.StoreRSSFunc == nil {
		panic("YoutubeSvcMock.StoreRSSFunc: method is nil but YoutubeSvc.StoreRSS was just called")
	}
	callInfo := struct {
		ChanID string
		Rss    string
	}{
		ChanID: chanID,
		Rss:    rss,
	}
	mock.lockStoreRSS.Lock()
	mock.calls.StoreRSS = append(mock.calls.StoreRSS, callInfo)
	mock.lockStoreRSS.Unlock()
	return mock.StoreRSSFunc(chanID, rss)
}

// StoreRSSCalls gets all the calls that were made to StoreRSS.
// Check the length with:
//     len(mockedYoutubeSvc.StoreRSSCalls())
func (mock *YoutubeSvcMock) StoreRSSCalls() []struct {
	ChanID string
	Rss    string
} {
	var calls []struct {
		ChanID string
		Rss    string
	}
	mock.lockStoreRSS.RLock()
	calls = mock.calls.StoreRSS
	mock.lockStoreRSS.RUnlock()
	return calls
}
