package strcache

func TestSetMemoryFetcherValues(values map[string]string) func(fetcher *MemoryFetcher) {
	return func(fetcher *MemoryFetcher) {
		fetcher.values = values
	}
}
