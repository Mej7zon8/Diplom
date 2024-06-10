package categorize

import (
	"github.com/k773/utils/synctools"
	"log/slog"
	"messenger/data/entities"
	chatstore "messenger/data/store/chat-store"
	"sync"
)

type job struct {
	s sync.Mutex

	concurrencyLimiter *synctools.ConcurrencyLimiter
}

// newJob creates a new categorization job.
func newJob() *job {
	return &job{
		concurrencyLimiter: synctools.NewConcurrencyLimiter(2),
	}
}

func (j *job) Run() {
	// List all chats:
	chats, e := chatstore.Instance.GetAllChats()
	if e != nil {
		slog.With("module", "jobs/jobs/categorize", "method", "Run").Warn("failed to list all chats", "error", e)
		return
	}
	var wg sync.WaitGroup
	// Run the job for each chat:
	for _, chat := range chats {
		wg.Add(1)
		j.concurrencyLimiter.Lock()
		go func(chat entities.ChatInfo) {
			defer j.concurrencyLimiter.Unlock()
			defer wg.Done()
			j.runForChat(chat.ID)
		}(chat)
	}
	wg.Wait()
}

func (j *job) runForChat(ref entities.ChatRef) {
	// Execute the job:
	e := j.WithWorker(ref, func(w *chatWorker) error {
		return w.Run()
	})
	// Log the error:
	if e != nil {
		slog.With("module", "jobs/jobs/categorize", "method", "runForChat").Warn("job failed for the chat", "error", e, "chat", ref)
	}
}

// WithWorker spawns a new worker for the job and executes the provided function with it.
// The callback is executed only after the concurrency limiter allows it.
func (j *job) WithWorker(chatRef entities.ChatRef, do func(w *chatWorker) error) error {
	return do(newChatWorker(chatRef))
}
