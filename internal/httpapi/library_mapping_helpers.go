package httpapi

func optionalLibrarySeriesType(value *SeriesType) *string {
	if value == nil {
		return nil
	}
	seriesType := string(*value)
	return &seriesType
}

func optionalInt64(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func libraryMatchSourceResponse(value *string) *LibraryScanItemMatchSource {
	if value == nil {
		return nil
	}
	source := LibraryScanItemMatchSource(*value)
	return &source
}
