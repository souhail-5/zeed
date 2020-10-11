package changelog

type File struct {
	Name     string
	Channel  string
	Priority int
	Hash     string
	Content  string
}

type Entry struct {
	Text     string
	Priority int
	Channel  Channel
}

type Channel struct {
	Id      string
	Entries []Entry
}

type ByPriority []Entry

func (entries ByPriority) Len() int           { return len(entries) }
func (entries ByPriority) Swap(i, j int)      { entries[i], entries[j] = entries[j], entries[i] }
func (entries ByPriority) Less(i, j int) bool { return entries[i].Priority > entries[j].Priority }
